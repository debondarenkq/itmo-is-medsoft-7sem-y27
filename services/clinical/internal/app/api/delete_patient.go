package api

import (
	"net/http"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/pkg/platform"
)

func (i *Implementation) DeletePatient(w http.ResponseWriter, r *http.Request) error {
	id, err := platform.ID(r)
	if err != nil {
		return err
	}
	if err = i.Module.DeletePatient(r.Context(), id); err != nil {
		return err
	}
	w.WriteHeader(http.StatusNoContent)
	return nil
}
