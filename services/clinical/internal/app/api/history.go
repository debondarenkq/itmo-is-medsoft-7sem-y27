package api

import (
	"net/http"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/pkg/platform"
)

func (i *Implementation) History(w http.ResponseWriter, r *http.Request) error {
	id, err := platform.ID(r)
	if err != nil {
		return err
	}
	p, err := platform.Pagination(r)
	if err != nil {
		return err
	}
	v, total, err := i.Module.History(r.Context(), id, p)
	if err != nil {
		return err
	}
	platform.List(w, v, total, p)
	return nil
}
