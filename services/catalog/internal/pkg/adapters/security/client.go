package security

import (
	"context"
	"encoding/json"
	"fmt"
	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/catalog/api"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/catalog/internal/pkg/adapters/audit"
	"github.com/getkin/kin-openapi/routers"
	openapirouter "github.com/getkin/kin-openapi/routers/gorillamux"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type contextKey int

const (
	tokenKey contextKey = iota
	staffKey
)

func Authorization(ctx context.Context) string {
	value, _ := ctx.Value(tokenKey).(string)
	return value
}
func StaffID(ctx context.Context) string { value, _ := ctx.Value(staffKey).(string); return value }

const serviceRequestTimeout = 2 * time.Second

type Authorizer struct {
	client      *ClientWithResponses
	auditClient *audit.ClientWithResponses
	routes      routers.Router
}

func New() (*Authorizer, error) {
	base := os.Getenv("AUTH_URL")
	u, err := url.Parse(base)
	if err != nil || u.Host == "" || (u.Scheme != "http" && u.Scheme != "https") {
		return nil, fmt.Errorf("AUTH_URL is required")
	}
	data, err := os.ReadFile(os.Getenv("SERVICE_KEY_FILE"))
	if err != nil {
		return nil, err
	}
	key := strings.TrimSpace(string(data))
	if len(key) < 32 {
		return nil, fmt.Errorf("invalid service key")
	}
	client, err := NewClientWithResponses(base, WithHTTPClient(&http.Client{Timeout: serviceRequestTimeout}), WithRequestEditorFn(func(_ context.Context, r *http.Request) error { r.Header.Set("X-Service-Key", key); return nil }))
	if err != nil {
		return nil, err
	}
	auditClient, err := audit.NewClientWithResponses(base, audit.WithHTTPClient(&http.Client{Timeout: serviceRequestTimeout}), audit.WithRequestEditorFn(func(_ context.Context, r *http.Request) error { r.Header.Set("X-Service-Key", key); return nil }))
	if err != nil {
		return nil, err
	}
	spec, err := contract.GetSwagger()
	if err != nil {
		return nil, err
	}
	spec.Servers = nil
	routes, err := openapirouter.NewRouter(spec)
	if err != nil {
		return nil, err
	}
	return &Authorizer{client: client, auditClient: auditClient, routes: routes}, nil
}
func writeProblem(w http.ResponseWriter, status int, code ErrorCode, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(Error{Error: ErrorDetail{Code: code, Message: message}}); err != nil {
		slog.Warn("encode access error", "error", err)
	}
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}
func (a *Authorizer) audit(ctx context.Context, input audit.HTTPAudit) error {
	if a.auditClient == nil {
		return fmt.Errorf("audit client is not configured")
	}
	response, err := a.auditClient.AppendHTTPAuditWithResponse(ctx, input)
	if err != nil {
		return err
	}
	if response.JSON200 == nil || !response.JSON200.Ok {
		return fmt.Errorf("identity audit HTTP %d", response.StatusCode())
	}
	return nil
}

func (a *Authorizer) auditEntry(r *http.Request, actor string) (audit.HTTPAudit, error) {
	if a.routes == nil {
		return audit.HTTPAudit{}, fmt.Errorf("OpenAPI router is not configured")
	}
	route, _, err := a.routes.FindRoute(r)
	if err != nil {
		return audit.HTTPAudit{}, err
	}
	data, err := json.Marshal(route.Operation.Extensions["x-audit"])
	if err != nil {
		return audit.HTTPAudit{}, err
	}
	var entry audit.HTTPAudit
	if err = json.Unmarshal(data, &entry); err != nil {
		return entry, err
	}
	entry.Actor = actor
	entry.Method = audit.HTTPAuditMethod(r.Method)
	entry.Path = r.URL.Path
	if !entry.EventType.Valid() || entry.Title == "" {
		return entry, fmt.Errorf("invalid OpenAPI audit metadata for %s", route.Operation.OperationID)
	}
	return entry, nil
}
func (a *Authorizer) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/api/") {
			next.ServeHTTP(w, r)
			return
		}
		w.Header().Set("Cache-Control", "no-store")
		token := r.Header.Get("Authorization")
		response, err := a.client.AuthorizeWithResponse(r.Context(), &AuthorizeParams{Authorization: &token}, AuthorizeJSONRequestBody{Method: HTTPMethod(r.Method), Path: r.URL.Path})
		if err != nil {
			writeProblem(w, http.StatusServiceUnavailable, DEPENDENCYUNAVAILABLE, "Сервис авторизации недоступен")
			return
		}
		if response.JSON401 != nil && response.JSON401.Error.Code == UNAUTHORIZED {
			writeProblem(w, http.StatusUnauthorized, UNAUTHORIZED, response.JSON401.Error.Message)
			return
		}
		if response.JSON403 != nil && response.JSON403.Error.Code == FORBIDDEN {
			writeProblem(w, http.StatusForbidden, FORBIDDEN, response.JSON403.Error.Message)
			return
		}
		claims := response.JSON200
		if claims == nil || claims.Sub == uuid.Nil || !claims.Role.Valid() {
			writeProblem(w, http.StatusServiceUnavailable, DEPENDENCYUNAVAILABLE, "Некорректный ответ авторизации")
			return
		}
		staffID := ""
		if claims.StaffID != nil {
			staffID = claims.StaffID.String()
		}
		ctx := context.WithValue(r.Context(), tokenKey, token)
		ctx = context.WithValue(ctx, staffKey, staffID)
		r = r.WithContext(ctx)
		if r.Method == http.MethodGet {
			next.ServeHTTP(w, r)
			return
		}
		actor := "account:" + claims.Sub.String()
		entry, err := a.auditEntry(r, actor)
		if err != nil {
			writeProblem(w, http.StatusServiceUnavailable, DEPENDENCYUNAVAILABLE, "Не удалось определить операцию аудита")
			return
		}
		if err := a.audit(ctx, entry); err != nil {
			writeProblem(w, http.StatusServiceUnavailable, DEPENDENCYUNAVAILABLE, "Сервис аудита недоступен")
			return
		}
		recorder := &statusWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(recorder, r)
		auditContext, cancel := context.WithTimeout(context.WithoutCancel(ctx), 2*time.Second)
		defer cancel()
		entry.Status = &recorder.status
		if err := a.audit(auditContext, entry); err != nil {
			slog.Error("audit outcome delivery failed", "path", r.URL.Path, "error", err)
		}
	})
}
