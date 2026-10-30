package api

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/gorilla/mux"
	nethttpmiddleware "github.com/oapi-codegen/nethttp-middleware"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/api/security"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/internal/models"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/internal/pkg/modules"
)

type SecurityAPI struct {
	Module      *modules.Security
	InternalKey string
	mu          sync.Mutex
	attempts    map[string]loginWindow
}

type loginWindow struct {
	start time.Time
	count int
}

func (api *SecurityAPI) Register(router *mux.Router) {
	api.attempts = map[string]loginWindow{}
	spec, err := contract.GetSwagger()
	if err != nil {
		panic(err)
	}
	validator := nethttpmiddleware.OapiRequestValidatorWithOptions(spec, &nethttpmiddleware.Options{DoNotValidateServers: true, Options: openapi3filter.Options{AuthenticationFunc: func(context.Context, *openapi3filter.AuthenticationInput) error { return nil }}, ErrorHandlerWithOpts: func(_ context.Context, err error, w http.ResponseWriter, r *http.Request, _ nethttpmiddleware.ErrorHandlerOpts) {
		securityError(w, models.ErrInvalidInput)
	}})
	contract.HandlerWithOptions(api, contract.GorillaServerOptions{BaseRouter: router, Middlewares: []contract.MiddlewareFunc{validator}, ErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) { securityError(w, models.ErrInvalidInput) }})
	router.HandleFunc("/security/openapi.json", func(w http.ResponseWriter, r *http.Request) { writeJSON(w, 200, spec) }).Methods("GET")
	router.HandleFunc("/internal/authorize", api.authorize).Methods("POST")
	router.HandleFunc("/internal/audit", api.audit).Methods("POST")
	router.Use(api.protect)
}

func securityError(w http.ResponseWriter, err error) {
	status, code := 500, "INTERNAL"
	message := "Внутренняя ошибка сервиса"
	switch {
	case errors.Is(err, models.ErrUnauthorized):
		status, code = 401, "UNAUTHORIZED"
		message = err.Error()
	case errors.Is(err, models.ErrForbidden):
		status, code = 403, "FORBIDDEN"
		message = err.Error()
	case errors.Is(err, models.ErrInvalidInput):
		status, code = 400, "INVALID_REQUEST"
		message = err.Error()
	case errors.Is(err, models.ErrConflict):
		status, code = 409, "CONFLICT"
		message = err.Error()
	case errors.Is(err, models.ErrNotFound):
		status, code = 404, "NOT_FOUND"
		message = err.Error()
	}
	if status == 500 {
		slog.Error("security request failed", "error", err)
	}
	writeJSON(w, status, map[string]any{"error": map[string]string{"code": code, "message": message}})
}

func decodeSecurity(w http.ResponseWriter, r *http.Request, value any) bool {
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		securityError(w, models.ErrInvalidInput)
		return false
	}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if decoder.Decode(value) != nil || decoder.Decode(new(any)) != io.EOF {
		securityError(w, models.ErrInvalidInput)
		return false
	}
	return true
}

func (api *SecurityAPI) claims(r *http.Request) (models.Claims, error) {
	return api.Module.Verify(strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer "))
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
