package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/api"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/models"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/pkg/modules"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/pkg/platform"
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
	r.HandleFunc("/api/v1/patients", platform.Handle(i.ListPatients)).Methods("GET")
	r.HandleFunc("/api/v1/patients", platform.Handle(i.CreatePatient)).Methods("POST")
	r.HandleFunc("/api/v1/patients/{id}", platform.Handle(i.GetPatient)).Methods("GET")
	r.HandleFunc("/api/v1/patients/{id}", platform.Handle(i.UpdatePatient)).Methods("PUT")
	r.HandleFunc("/api/v1/patients/{id}", platform.Handle(i.DeletePatient)).Methods("DELETE")
	r.HandleFunc("/api/v1/patients/{id}/record", platform.Handle(i.GetPatientRecord)).Methods("GET")
	r.HandleFunc("/api/v1/patients/{id}/record", platform.Handle(i.CreateRecord)).Methods("POST")
	r.HandleFunc("/api/v1/records/{id}", platform.Handle(i.GetRecord)).Methods("GET")
	r.HandleFunc("/api/v1/records/{id}/changes", platform.Handle(i.Save)).Methods("POST")
	r.HandleFunc("/api/v1/records/{id}/history", platform.Handle(i.History)).Methods("GET")
	r.HandleFunc("/api/v1/records/{id}/state", platform.Handle(i.StateAt)).Methods("GET")
}

func parseStateQuery(r *http.Request) (models.StateQuery, error) {
	var q models.StateQuery
	values := r.URL.Query()
	_, hasAt := values["at"]
	_, hasVersion := values["version"]
	if hasAt == hasVersion {
		return q, platform.Bad("STATE_QUERY", "Укажите ровно один параметр: at или version")
	}
	if hasAt {
		t, err := time.Parse(time.RFC3339Nano, values.Get("at"))
		if err != nil {
			return q, platform.Bad("STATE_QUERY", "at должен быть датой и временем RFC 3339")
		}
		t = t.UTC()
		q.At = &t
	} else {
		v, err := strconv.ParseInt(values.Get("version"), 10, 64)
		if err != nil || v < 1 {
			return q, platform.Bad("STATE_QUERY", "version должен быть положительным целым числом")
		}
		q.Version = &v
	}
	return q, nil
}
