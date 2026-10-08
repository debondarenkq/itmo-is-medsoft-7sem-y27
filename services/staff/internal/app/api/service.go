package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/api"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/internal/models"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/internal/pkg/modules"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/internal/pkg/platform"
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
	r.HandleFunc("/api/v1/staff", platform.Handle(i.List)).Methods("GET")
	r.HandleFunc("/api/v1/staff", platform.Handle(i.Create)).Methods("POST")
	r.HandleFunc("/api/v1/staff/{id}", platform.Handle(i.Get)).Methods("GET")
	r.HandleFunc("/api/v1/staff/{id}", platform.Handle(i.Update)).Methods("PUT")
	r.HandleFunc("/api/v1/staff/{id}", platform.Handle(i.Delete)).Methods("DELETE")
}

type input struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Position  string `json:"position"`
}

func (v input) model() models.Staff {
	return models.Staff{FirstName: v.FirstName, LastName: v.LastName, Position: v.Position}
}
