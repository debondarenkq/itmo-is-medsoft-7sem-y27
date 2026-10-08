package api

import (
	"net/http"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/catalog/internal/pkg/platform"
)

func (i *Implementation) List(w http.ResponseWriter, r *http.Request) error {
	p, err := platform.Pagination(r)
	if err != nil {
		return err
	}
	v, total, err := i.Module.List(r.Context(), p)
	if err != nil {
		return err
	}
	platform.List(w, v, total, p)
	return nil
}
