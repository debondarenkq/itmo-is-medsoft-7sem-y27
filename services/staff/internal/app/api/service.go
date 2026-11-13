package api

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/api"
	internalapi "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/api/internalapi"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/internal/models"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/internal/pkg/adapters/security"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/internal/pkg/modules"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	nethttpmiddleware "github.com/oapi-codegen/nethttp-middleware"
)

type Deps struct{ Module *modules.Module }

type Implementation struct{ Deps }

var _ contract.StrictServerInterface = (*Implementation)(nil)

func New(deps Deps) *Implementation { return &Implementation{Deps: deps} }

func (implementation *Implementation) Register(router *mux.Router) {
	router.HandleFunc("/openapi.json", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, json.RawMessage(contract.Specification()))
	}).Methods(http.MethodGet)
	auth, err := security.New()
	if err != nil {
		panic(err)
	}
	router.Use(auth.Middleware)
	spec, err := contract.GetSwagger()
	if err != nil {
		panic(err)
	}
	validator := nethttpmiddleware.OapiRequestValidatorWithOptions(spec, &nethttpmiddleware.Options{
		DoNotValidateServers: true,
		Options:              openapi3filter.Options{AuthenticationFunc: func(context.Context, *openapi3filter.AuthenticationInput) error { return nil }},
		ErrorHandlerWithOpts: func(_ context.Context, err error, w http.ResponseWriter, request *http.Request, options nethttpmiddleware.ErrorHandlerOpts) {
			writeValidationError(w, options.StatusCode, err)
		},
	})
	strict := contract.NewStrictHandlerWithOptions(implementation, nil, contract.StrictHTTPServerOptions{
		RequestErrorHandlerFunc:  writeRequestError,
		ResponseErrorHandlerFunc: writeError,
	})
	contract.HandlerWithOptions(strict, contract.GorillaServerOptions{
		BaseRouter:       router,
		Middlewares:      []contract.MiddlewareFunc{validator},
		ErrorHandlerFunc: writeRequestError,
	})
}

func writeError(w http.ResponseWriter, request *http.Request, err error) {
	status, code := http.StatusInternalServerError, contract.INTERNAL
	switch {
	case errors.Is(err, models.ErrForbidden):
		status, code = http.StatusForbidden, contract.FORBIDDEN
	case errors.Is(err, models.ErrInvalidInput):
		status, code = http.StatusBadRequest, contract.INVALIDREQUEST
	case errors.Is(err, models.ErrNotFound):
		status, code = http.StatusNotFound, contract.NOTFOUND
	case errors.Is(err, models.ErrConflict):
		status, code = http.StatusConflict, contract.CONFLICT
	}
	message := err.Error()
	if status == http.StatusInternalServerError {
		slog.Error("request failed", "method", request.Method, "path", request.URL.Path, "error", err)
		message = "Внутренняя ошибка сервиса"
	}
	if strings.HasPrefix(request.URL.Path, "/internal/") {
		writeJSON(w, status, internalapi.Error{Error: internalapi.ErrorDetail{Code: internalapi.ErrorCode(code), Message: message}})
		return
	}
	writeJSON(w, status, contract.Error{Error: contract.ErrorDetail{Code: code, Message: message}})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		slog.Warn("encode response", "error", err)
	}
}

func writeRequestError(w http.ResponseWriter, r *http.Request, _ error) {
	if strings.HasPrefix(r.URL.Path, "/internal/") {
		writeError(w, r, models.ErrInvalidInput)
		return
	}
	writeJSON(w, http.StatusBadRequest, contract.Error{Error: contract.ErrorDetail{
		Code: contract.INVALIDREQUEST, Message: "Некорректный формат запроса",
	}})
}

func writeValidationError(w http.ResponseWriter, status int, err error) {
	detail := contract.ErrorDetail{Code: contract.INVALIDREQUEST, Message: "Проверьте заполнение полей"}
	var schemaError *openapi3.SchemaError
	if errors.As(err, &schemaError) && schemaError.Schema != nil {
		message := "Некорректное значение"
		switch schemaError.SchemaField {
		case "minLength", "maxLength":
			if maximum := schemaError.Schema.MaxLength; maximum != nil {
				message = fmt.Sprintf("От %d до %d символов", schemaError.Schema.MinLength, *maximum)
			}
		}
		if path := schemaError.JSONPointer(); len(path) > 0 {
			fields := map[string]string{strings.Join(path, "."): message}
			detail.Fields = &fields
		}
	}
	writeJSON(w, status, contract.Error{Error: detail})
}

func fromStaffInput(value contract.StaffInput) models.Staff {
	return models.Staff{
		FirstName: value.FirstName,
		LastName:  value.LastName,
		Position:  value.Position,
		Role:      roleValue(value.Role),
	}
}

func toStaff(value models.Staff) (contract.Staff, error) {
	id, err := uuid.Parse(value.ID)
	if err != nil {
		return contract.Staff{}, fmt.Errorf("map staff id: %w", err)
	}
	return contract.Staff{
		ID:        id,
		FirstName: value.FirstName,
		LastName:  value.LastName,
		Position:  value.Position,
		Role:      contract.StaffRole(value.Role),
		CreatedAt: value.CreatedAt,
		UpdatedAt: value.UpdatedAt,
		DeletedAt: value.DeletedAt,
	}, nil
}

func roleValue(role *contract.StaffRole) string {
	if role == nil {
		return "other"
	}
	return string(*role)
}

var _ internalapi.StrictServerInterface = (*Implementation)(nil)

func (i *Implementation) RegisterInternal(router *mux.Router, key string) {
	spec, err := internalapi.GetSwagger()
	if err != nil {
		panic(err)
	}
	validator := nethttpmiddleware.OapiRequestValidatorWithOptions(spec, &nethttpmiddleware.Options{DoNotValidateServers: true, Options: openapi3filter.Options{AuthenticationFunc: func(context.Context, *openapi3filter.AuthenticationInput) error { return nil }}, ErrorHandlerWithOpts: func(_ context.Context, _ error, w http.ResponseWriter, r *http.Request, _ nethttpmiddleware.ErrorHandlerOpts) {
		writeRequestError(w, r, nil)
	}})
	strict := internalapi.NewStrictHandlerWithOptions(i, nil, internalapi.StrictHTTPServerOptions{RequestErrorHandlerFunc: writeRequestError, ResponseErrorHandlerFunc: writeError})
	internalapi.HandlerWithOptions(strict, internalapi.GorillaServerOptions{BaseRouter: router, Middlewares: []internalapi.MiddlewareFunc{validator}, ErrorHandlerFunc: writeRequestError})
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/internal/") && (key == "" || subtle.ConstantTimeCompare([]byte(r.Header.Get("X-Service-Key")), []byte(key)) != 1) {
				writeError(w, r, models.ErrForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	})
}

func toStaffProfile(s models.Staff) (internalapi.ProfileJSONResponse, error) {
	id, err := uuid.Parse(s.ID)
	if err != nil {
		return internalapi.ProfileJSONResponse{}, err
	}
	return internalapi.ProfileJSONResponse{ID: id, FirstName: s.FirstName, LastName: s.LastName, Position: &s.Position, Role: internalapi.Role(s.Role), CreatedAt: &s.CreatedAt, UpdatedAt: &s.UpdatedAt, DeletedAt: s.DeletedAt, RoleRevision: s.RoleRevision}, nil
}
