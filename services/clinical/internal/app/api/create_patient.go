package api

import (
	"net/http"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/models"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/pkg/platform"
)

func (i *Implementation) CreatePatient(w http.ResponseWriter, r *http.Request) error {
	var in models.PatientInput
	if err := platform.Decode(w, r, &in); err != nil {
		return err
	}
	v, err := i.Module.CreatePatient(r.Context(), in)
	if err != nil {
		return err
	}
	platform.JSON(w, 201, v)
	return nil
}
