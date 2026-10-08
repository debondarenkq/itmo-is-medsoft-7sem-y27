package api

import (
	"net/http"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/internal/pkg/platform"
)

func (i *Implementation) Create(w http.ResponseWriter, r *http.Request) error {
	var in input
	if err := platform.Decode(w, r, &in); err != nil {
		return err
	}
	v, err := i.Module.Create(r.Context(), in.model())
	if err != nil {
		return err
	}
	platform.JSON(w, 201, v)
	return nil
}
