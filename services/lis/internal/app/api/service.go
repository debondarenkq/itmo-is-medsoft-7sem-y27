package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/lis/api"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/lis/internal/models"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/lis/internal/pkg/modules"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/gorilla/mux"
	middleware "github.com/oapi-codegen/nethttp-middleware"
)

type Implementation struct{ Module *modules.Module }

var _ contract.StrictServerInterface = (*Implementation)(nil)

func writeError(w http.ResponseWriter, _ *http.Request, err error) {
	status, code := http.StatusBadGateway, "HIS_UNAVAILABLE"
	switch {
	case errors.Is(err, models.ErrInvalid):
		status, code = 400, "INVALID_REQUEST"
	case errors.Is(err, models.ErrConflict):
		status, code = 409, "CONFLICT"
	case errors.Is(err, models.ErrNotFound):
		status, code = 404, "NOT_FOUND"
	case errors.Is(err, models.ErrUnconfigured):
		status, code = 409, "NOT_CONFIGURED"
	}
	problem := contract.LabError{}
	problem.Error.Code = code
	problem.Error.Message = err.Error()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(problem)
}

func (i *Implementation) Register(router *mux.Router) {
	spec, err := contract.GetSwagger()
	if err != nil {
		panic(err)
	}
	validator := middleware.OapiRequestValidatorWithOptions(spec, &middleware.Options{DoNotValidateServers: true, Options: openapi3filter.Options{AuthenticationFunc: func(context.Context, *openapi3filter.AuthenticationInput) error { return nil }}, ErrorHandlerWithOpts: func(_ context.Context, _ error, w http.ResponseWriter, r *http.Request, _ middleware.ErrorHandlerOpts) {
		writeError(w, r, models.ErrInvalid)
	}})
	strict := contract.NewStrictHandlerWithOptions(i, nil, contract.StrictHTTPServerOptions{RequestErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, _ error) { writeError(w, r, models.ErrInvalid) }, ResponseErrorHandlerFunc: writeError})
	contract.HandlerWithOptions(strict, contract.GorillaServerOptions{BaseRouter: router, Middlewares: []contract.MiddlewareFunc{validator}})
}
