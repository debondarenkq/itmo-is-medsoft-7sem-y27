package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/api"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/internal/pkg/modules"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/internal/pkg/platform"
)

type Deps struct{ Module *modules.Module }
type Implementation struct{ Deps }

var _ contract.StrictServerInterface = (*Implementation)(nil)

func New(deps Deps) *Implementation { return &Implementation{Deps: deps} }

func (i *Implementation) Register(r *mux.Router) {
	r.HandleFunc("/openapi.json", func(w http.ResponseWriter, _ *http.Request) {
		platform.JSON(w, http.StatusOK, json.RawMessage(contract.Specification()))
	}).Methods(http.MethodGet)
	spec, err := contract.GetSwagger()
	if err != nil {
		panic(err)
	}
	r.Use(platform.ContractMiddleware(spec))
	strict := contract.NewStrictHandlerWithOptions(i, nil, contract.StrictHTTPServerOptions{
		RequestErrorHandlerFunc:  platform.WriteRequestError,
		ResponseErrorHandlerFunc: platform.WriteError,
	})
	contract.HandlerWithOptions(strict, contract.GorillaServerOptions{
		BaseRouter:       r,
		ErrorHandlerFunc: platform.WriteRequestError,
	})
}
