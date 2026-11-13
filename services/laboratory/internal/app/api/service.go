package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/laboratory/api"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/laboratory/internal/models"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/laboratory/internal/pkg/modules"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/gorilla/mux"
	middleware "github.com/oapi-codegen/nethttp-middleware"
)

type Implementation struct{ Module *modules.Module }

var _ contract.StrictServerInterface = (*Implementation)(nil)

func writeError(w http.ResponseWriter, r *http.Request, err error) {
	status, code, issue := http.StatusInternalServerError, "INTERNAL", contract.Exception
	switch {
	case errors.Is(err, models.ErrInvalidInput):
		status, code, issue = 400, "INVALID_REQUEST", contract.Invalid
	case errors.Is(err, models.ErrConflict):
		status, code, issue = 409, "CONFLICT", contract.Conflict
	case errors.Is(err, models.ErrNotFound):
		status, code, issue = 404, "NOT_FOUND", contract.NotFound
	case errors.Is(err, models.ErrUnavailable):
		status, code, issue = 503, "DEPENDENCY_UNAVAILABLE", contract.Transient
	}
	message := err.Error()
	if status == 500 {
		message = "Внутренняя ошибка лабораторного сервиса"
	}
	var result any
	if strings.Contains(r.URL.Path, "/fhir/") {
		w.Header().Set("Content-Type", models.FHIRMediaType)
		result = models.Outcome(issue, message)
	} else {
		w.Header().Set("Content-Type", "application/json")
		problem := contract.LabError{}
		problem.Error.Code = code
		problem.Error.Message = message
		result = problem
	}
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(result)
}

func init() {
	openapi3filter.RegisterBodyDecoder(models.FHIRMediaType, openapi3filter.JSONBodyDecoder)
}

func (i *Implementation) Register(router *mux.Router) {
	spec, err := contract.GetSwagger()
	if err != nil {
		panic(err)
	}
	validator := middleware.OapiRequestValidatorWithOptions(spec, &middleware.Options{DoNotValidateServers: true, Options: openapi3filter.Options{AuthenticationFunc: func(context.Context, *openapi3filter.AuthenticationInput) error { return nil }}, ErrorHandlerWithOpts: func(_ context.Context, _ error, w http.ResponseWriter, r *http.Request, _ middleware.ErrorHandlerOpts) {
		writeError(w, r, models.ErrInvalidInput)
	}})
	strict := contract.NewStrictHandlerWithOptions(i, nil, contract.StrictHTTPServerOptions{RequestErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, _ error) { writeError(w, r, models.ErrInvalidInput) }, ResponseErrorHandlerFunc: writeError})
	contract.HandlerWithOptions(strict, contract.GorillaServerOptions{BaseRouter: router, Middlewares: []contract.MiddlewareFunc{validator}})
}
