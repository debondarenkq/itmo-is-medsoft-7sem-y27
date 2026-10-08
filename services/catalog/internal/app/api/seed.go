package api

import (
	"net/http"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/catalog/internal/pkg/platform"
)

func (i *Implementation) Seed(w http.ResponseWriter, r *http.Request) error {
	v, err := i.Module.Seed(r.Context())
	if err != nil {
		return err
	}
	platform.JSON(w, 200, map[string]any{"items": v})
	return nil
}
