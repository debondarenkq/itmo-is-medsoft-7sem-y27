package api

import (
	"net/http"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/catalog/internal/pkg/platform"
)

func (i *Implementation) Get(w http.ResponseWriter, r *http.Request) error {
	id, err := platform.ID(r)
	if err != nil {
		return err
	}
	v, err := i.Module.Get(r.Context(), id)
	if err != nil {
		return err
	}
	platform.JSON(w, 200, v)
	return nil
}
