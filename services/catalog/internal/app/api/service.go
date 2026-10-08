package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/catalog/api"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/catalog/internal/models"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/catalog/internal/pkg/modules"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/catalog/internal/pkg/platform"
)

type Deps struct{ Module *modules.Module }
type Implementation struct{ Deps }

func New(deps Deps) *Implementation {
	return &Implementation{Deps: deps}
}

func (i *Implementation) Register(r *mux.Router) {
	r.HandleFunc("/openapi.json", func(w http.ResponseWriter, _ *http.Request) {
		platform.JSON(w, http.StatusOK, json.RawMessage(contract.Specification()))
	}).Methods("GET")
	r.HandleFunc("/api/v1/diagnoses/seed", platform.Handle(i.Seed)).Methods("POST")
	r.HandleFunc("/api/v1/diagnoses", platform.Handle(i.List)).Methods("GET")
	r.HandleFunc("/api/v1/diagnoses", platform.Handle(i.Create)).Methods("POST")
	r.HandleFunc("/api/v1/diagnoses/{id}", platform.Handle(i.Get)).Methods("GET")
	r.HandleFunc("/api/v1/diagnoses/{id}", platform.Handle(i.Update)).Methods("PUT")
	r.HandleFunc("/api/v1/diagnoses/{id}", platform.Handle(i.Delete)).Methods("DELETE")
}

type input struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

func (v input) model() models.Diagnosis {
	return models.Diagnosis{Code: v.Code, Name: v.Name}
}
