package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/internal/models"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/internal/pkg/modules"
)

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (implementation *Implementation) protect(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), remoteAddressKey, r.RemoteAddr)
		r = r.WithContext(ctx)
		path := r.URL.Path
		if strings.HasPrefix(path, "/internal/") {
			if !implementation.trustedService(r) {
				writeError(w, r, models.ErrForbidden)
				return
			}
			next.ServeHTTP(w, r)
			return
		}
		if !strings.HasPrefix(path, "/api/") {
			next.ServeHTTP(w, r)
			return
		}
		w.Header().Set("Cache-Control", "no-store")
		if path == "/api/v1/auth/login" || path == "/api/v1/auth/setup" || path == "/api/v1/auth/refresh" {
			w.Header().Set("Cache-Control", "no-store")
			next.ServeHTTP(w, r)
			return
		}
		claims, err := implementation.Module.VerifyAccess(r.Context(), strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer "))
		if err != nil {
			writeError(w, r, implementation.auditFailure(r.Context(), err, models.AuditEventKindAccessDenied, "Неавторизованный запрос "+r.Method+" "+path))
			return
		}
		if !modules.Allowed(claims, r.Method, path) {
			if err := implementation.Module.Audit(r.Context(), "account:"+claims.Subject, "Отказ в доступе "+r.Method+" "+path, models.AuditEventKindAccessDenied); err != nil {
				writeError(w, r, err)
				return
			}
			writeError(w, r, models.ErrForbidden)
			return
		}
		r = r.WithContext(context.WithValue(r.Context(), claimsKey, claims))
		w.Header().Set("Cache-Control", "no-store")
		if r.Method == http.MethodGet {
			next.ServeHTTP(w, r)
			return
		}
		actor := "account:" + claims.Subject
		event, err := implementation.auditEntry(r, actor)
		if err != nil {
			writeError(w, r, err)
			return
		}
		if err := implementation.Module.AuditHTTP(r.Context(), event); err != nil {
			writeError(w, r, err)
			return
		}
		recorder := &statusWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(recorder, r)
		auditContext, cancel := context.WithTimeout(context.WithoutCancel(r.Context()), 2*time.Second)
		defer cancel()
		event.Status = &recorder.status
		if err := implementation.Module.AuditHTTP(auditContext, event); err != nil {
			slog.Error("audit outcome delivery failed", "path", path, "error", err)
		}
	})
}

func (i *Implementation) auditEntry(r *http.Request, actor string) (models.HTTPAudit, error) {
	route, _, err := i.auditRoutes.FindRoute(r)
	if err != nil {
		return models.HTTPAudit{}, err
	}
	data, err := json.Marshal(route.Operation.Extensions["x-audit"])
	if err != nil {
		return models.HTTPAudit{}, err
	}
	var entry models.HTTPAudit
	if err = json.Unmarshal(data, &entry); err != nil {
		return entry, err
	}
	entry.Actor = actor
	entry.Method = models.HTTPAuditMethod(r.Method)
	entry.Path = r.URL.Path
	if !entry.EventType.Valid() || entry.Title == "" {
		return entry, fmt.Errorf("invalid OpenAPI audit metadata for %s", route.Operation.OperationID)
	}
	return entry, nil
}
