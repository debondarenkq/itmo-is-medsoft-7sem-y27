package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	nethttpmiddleware "github.com/oapi-codegen/nethttp-middleware"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/catalog/api"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/catalog/internal/models"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/catalog/internal/pkg/modules"
)

type Deps struct{ Module *modules.Module }
type Implementation struct{ Deps }

var _ contract.StrictServerInterface = (*Implementation)(nil)

func New(deps Deps) *Implementation { return &Implementation{Deps: deps} }

func (implementation *Implementation) Register(router *mux.Router) {
	router.HandleFunc("/openapi.json", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, json.RawMessage(contract.Specification()))
	}).Methods(http.MethodGet)
	spec, err := contract.GetSwagger()
	if err != nil {
		panic(err)
	}
	validator := nethttpmiddleware.OapiRequestValidatorWithOptions(spec, &nethttpmiddleware.Options{
		DoNotValidateServers: true,
		ErrorHandlerWithOpts: func(_ context.Context, err error, w http.ResponseWriter, request *http.Request, options nethttpmiddleware.ErrorHandlerOpts) {
			writeJSON(w, options.StatusCode, contract.Error{Error: contract.ErrorDetail{Code: "INVALID_REQUEST", Message: err.Error()}})
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
	status, code := http.StatusInternalServerError, "INTERNAL"
	switch {
	case errors.Is(err, models.ErrInvalidInput):
		status, code = http.StatusBadRequest, "INVALID_REQUEST"
	case errors.Is(err, models.ErrNotFound):
		status, code = http.StatusNotFound, "NOT_FOUND"
	case errors.Is(err, models.ErrConflict):
		status, code = http.StatusConflict, "CONFLICT"
	}
	message := err.Error()
	if status == http.StatusInternalServerError {
		slog.Error("request failed", "method", request.Method, "path", request.URL.Path, "error", err)
		message = "Внутренняя ошибка сервиса"
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

func writeRequestError(w http.ResponseWriter, request *http.Request, err error) {
	writeError(w, request, fmt.Errorf("%w: %s", models.ErrInvalidInput, err))
}

func fromDiagnosisInput(value contract.DiagnosisInput) models.Diagnosis {
	return models.Diagnosis{
		Code: value.Code,
		Name: value.Name,
	}
}

func toDiagnosis(value models.Diagnosis) (contract.Diagnosis, error) {
	id, err := uuid.Parse(value.ID)
	if err != nil {
		return contract.Diagnosis{}, fmt.Errorf("map diagnosis id: %w", err)
	}
	return contract.Diagnosis{
		ID:        id,
		Code:      value.Code,
		Name:      value.Name,
		CreatedAt: value.CreatedAt,
		UpdatedAt: value.UpdatedAt,
		DeletedAt: value.DeletedAt,
	}, nil
}
