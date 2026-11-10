package api

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"errors"
	"fmt"
	auditcontract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/api/audit"
	internalapi "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/api/internalapi"
	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/api/security"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/internal/models"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/internal/pkg/modules"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	openapirouter "github.com/getkin/kin-openapi/routers/gorillamux"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	nethttpmiddleware "github.com/oapi-codegen/nethttp-middleware"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Deps struct {
	Module      *modules.Security
	InternalKey string
}
type Implementation struct {
	auditRoutes routers.Router
	Deps
	mu       sync.Mutex
	attempts map[string]loginWindow
}
type loginWindow struct {
	start time.Time
	count int
}
type contextKey int

const (
	claimsKey contextKey = iota
	remoteAddressKey
)

var _ auditcontract.StrictServerInterface = (*Implementation)(nil)

var _ contract.StrictServerInterface = (*Implementation)(nil)
var _ internalapi.StrictServerInterface = (*Implementation)(nil)

func init() {
	openapi3filter.RegisterBodyDecoder("application/json", strictJSONBodyDecoder)
}

func strictJSONBodyDecoder(body io.Reader, _ http.Header, _ *openapi3.SchemaRef, _ openapi3filter.EncodingFn) (any, error) {
	decoder := json.NewDecoder(body)
	decoder.UseNumber()
	var value any
	if err := decoder.Decode(&value); err != nil {
		return nil, &openapi3filter.ParseError{Kind: openapi3filter.KindInvalidFormat, Cause: err}
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		if err == nil {
			err = fmt.Errorf("multiple JSON values in request body")
		}
		return nil, &openapi3filter.ParseError{Kind: openapi3filter.KindInvalidFormat, Cause: err}
	}
	return value, nil
}

func New(deps Deps) *Implementation {
	return &Implementation{Deps: deps, attempts: map[string]loginWindow{}}
}
func (implementation *Implementation) Register(router *mux.Router) {
	spec, err := contract.GetSwagger()
	if err != nil {
		panic(err)
	}
	spec.InternalizeRefs(context.Background(), nil)
	implementation.auditRoutes, err = openapirouter.NewRouter(spec)
	if err != nil {
		panic(err)
	}
	strict := contract.NewStrictHandlerWithOptions(implementation, nil, contract.StrictHTTPServerOptions{RequestErrorHandlerFunc: writeRequestError, ResponseErrorHandlerFunc: writeError})
	contract.HandlerWithOptions(strict, contract.GorillaServerOptions{BaseRouter: router, Middlewares: []contract.MiddlewareFunc{requestValidator(spec)}, ErrorHandlerFunc: writeRequestError})
	privateSpec, err := internalapi.GetSwagger()
	if err != nil {
		panic(err)
	}
	private := internalapi.NewStrictHandlerWithOptions(implementation, nil, internalapi.StrictHTTPServerOptions{RequestErrorHandlerFunc: writeRequestError, ResponseErrorHandlerFunc: writeError})
	internalapi.HandlerWithOptions(private, internalapi.GorillaServerOptions{BaseRouter: router, Middlewares: []internalapi.MiddlewareFunc{requestValidator(privateSpec)}, ErrorHandlerFunc: writeRequestError})
	auditSpec, err := auditcontract.GetSwagger()
	if err != nil {
		panic(err)
	}
	auditHandler := auditcontract.NewStrictHandlerWithOptions(implementation, nil, auditcontract.StrictHTTPServerOptions{RequestErrorHandlerFunc: writeRequestError, ResponseErrorHandlerFunc: writeError})
	auditcontract.HandlerWithOptions(auditHandler, auditcontract.GorillaServerOptions{BaseRouter: router, Middlewares: []auditcontract.MiddlewareFunc{requestValidator(auditSpec)}, ErrorHandlerFunc: writeRequestError})
	router.HandleFunc("/security/openapi.json", func(w http.ResponseWriter, _ *http.Request) { writeJSON(w, 200, spec) }).Methods("GET")
	router.Use(implementation.protect)
}
func requestValidator(spec *openapi3.T) func(http.Handler) http.Handler {
	return nethttpmiddleware.OapiRequestValidatorWithOptions(spec, &nethttpmiddleware.Options{DoNotValidateServers: true, Options: openapi3filter.Options{AuthenticationFunc: func(context.Context, *openapi3filter.AuthenticationInput) error { return nil }}, ErrorHandlerWithOpts: func(_ context.Context, _ error, w http.ResponseWriter, r *http.Request, _ nethttpmiddleware.ErrorHandlerOpts) {
		writeRequestError(w, r, nil)
	}})
}
func writeRequestError(w http.ResponseWriter, r *http.Request, _ error) {
	writeError(w, r, models.ErrInvalidInput)
}
func writeError(w http.ResponseWriter, r *http.Request, err error) {
	status, code := 500, contract.INTERNAL
	message := "Внутренняя ошибка сервиса"
	switch {
	case errors.Is(err, models.ErrUnavailable):
		status, code = 503, contract.DEPENDENCYUNAVAILABLE
	case errors.Is(err, models.ErrUnauthorized):
		status, code = 401, contract.UNAUTHORIZED
	case errors.Is(err, models.ErrForbidden):
		status, code = 403, contract.FORBIDDEN
	case errors.Is(err, models.ErrInvalidInput):
		status, code = 400, contract.INVALIDREQUEST
	case errors.Is(err, models.ErrConflict):
		status, code = 409, contract.CONFLICT
	case errors.Is(err, models.ErrNotFound):
		status, code = 404, contract.NOTFOUND
	}
	if status == 500 {
		slog.Error("identity request failed", "error", err)
	} else {
		message = err.Error()
	}
	if strings.HasPrefix(r.URL.Path, "/internal/") {
		writeJSON(w, status, internalapi.Error{Error: internalapi.ErrorDetail{Code: internalapi.ErrorCode(code), Message: message}})
		return
	}
	response := contract.SecurityError{}
	response.Error.Code = code
	response.Error.Message = message
	writeJSON(w, status, response)
}
func requestClaims(ctx context.Context) (models.Claims, error) {
	claims, ok := ctx.Value(claimsKey).(models.Claims)
	if !ok {
		return claims, models.ErrUnauthorized
	}
	return claims, nil
}
func pagination(limit, offset *int) (int, int) {
	l, o := 50, 0
	if limit != nil {
		l = *limit
	}
	if offset != nil {
		o = *offset
	}
	if l < 1 {
		l = 1
	}
	if l > 200 {
		l = 200
	}
	if o < 0 {
		o = 0
	}
	return l, o
}
func toStaffID(id *string) (*uuid.UUID, error) {
	if id == nil {
		return nil, nil
	}
	value, err := uuid.Parse(*id)
	if err != nil {
		return nil, fmt.Errorf("map staff id: %w", err)
	}
	return &value, nil
}
func toAccount(a models.Account) (contract.Account, error) {
	id, err := uuid.Parse(a.ID)
	if err != nil {
		return contract.Account{}, err
	}
	staffID, err := toStaffID(a.StaffID)
	if err != nil {
		return contract.Account{}, err
	}
	return contract.Account{ID: id, Username: a.Username, CreatedAt: a.CreatedAt, UpdatedAt: a.UpdatedAt, IsActive: a.IsActive, StaffID: staffID, Role: contract.AccountRole(a.Role), StaffName: a.StaffName}, nil
}
func toSession(s models.Session) (contract.Session, error) {
	id, err := uuid.Parse(s.AccountID)
	if err != nil {
		return contract.Session{}, err
	}
	staffID, err := toStaffID(s.StaffID)
	if err != nil {
		return contract.Session{}, err
	}
	return contract.Session{AccountID: id, StaffID: staffID, AccessToken: s.AccessToken, RefreshToken: s.RefreshToken, ExpiresAt: s.ExpiresAt, Role: s.Role}, nil
}
func (implementation *Implementation) Authorize(ctx context.Context, request internalapi.AuthorizeRequestObject) (internalapi.AuthorizeResponseObject, error) {
	token := ""
	if request.Params.Authorization != nil {
		token = strings.TrimPrefix(*request.Params.Authorization, "Bearer ")
	}
	claims, err := implementation.Module.VerifyAccess(ctx, token)
	input := request.Body
	if err != nil {
		return nil, implementation.auditFailure(ctx, err, models.AuditEventKindAccessDenied, "Неавторизованный запрос "+string(input.Method)+" "+input.Path)
	}
	if !modules.Allowed(claims, string(input.Method), input.Path) {
		if err := implementation.Module.Audit(ctx, "account:"+claims.Subject, "Отказ в доступе "+string(input.Method)+" "+input.Path, models.AuditEventKindAccessDenied); err != nil {
			return nil, err
		}
		return nil, models.ErrForbidden
	}
	subject, err := uuid.Parse(claims.Subject)
	if err != nil {
		return nil, err
	}
	jti, err := uuid.Parse(claims.ID)
	if err != nil {
		return nil, err
	}
	var staffID *uuid.UUID
	if claims.StaffID != "" {
		staffID, err = toStaffID(&claims.StaffID)
		if err != nil {
			return nil, err
		}
	}
	return internalapi.Authorize200JSONResponse{Sub: subject, StaffID: staffID, Role: internalapi.ClaimsRole(claims.Role), Iss: claims.Issuer, Aud: claims.Audience, Exp: claims.Expires, Iat: claims.Issued, Jti: jti}, nil
}
func (implementation *Implementation) AppendHTTPAudit(ctx context.Context, request auditcontract.AppendHTTPAuditRequestObject) (auditcontract.AppendHTTPAuditResponseObject, error) {
	input := request.Body
	event := models.HTTPAudit{Actor: input.Actor, Method: models.HTTPAuditMethod(input.Method), Path: input.Path, Status: input.Status, EventType: models.AuditEventKind(input.EventType), Title: input.Title}
	if err := implementation.Module.AuditHTTP(ctx, event); err != nil {
		return nil, err
	}
	return auditcontract.AppendHTTPAudit200JSONResponse{Ok: true}, nil
}
func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		slog.Warn("encode response", "error", err)
	}
}

func (implementation *Implementation) trustedService(r *http.Request) bool {
	return implementation.InternalKey != "" && subtle.ConstantTimeCompare([]byte(r.Header.Get("X-Service-Key")), []byte(implementation.InternalKey)) == 1
}
