package api

import (
	"net/http"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/pkg/platform"
)

func (i *Implementation) GetPatientRecord(w http.ResponseWriter, r *http.Request) error {
	id, err := platform.ID(r)
	if err != nil {
		return err
	}
	v, err := i.Module.GetPatientRecord(r.Context(), id)
	if err != nil {
		return err
	}
	platform.JSON(w, 200, v)
	return nil
}
