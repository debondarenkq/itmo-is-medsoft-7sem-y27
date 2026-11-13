package modules

import (
	"context"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/catalog/internal/models"
)

func (module *Module) Get(ctx context.Context, id string) (models.Diagnosis, error) {
	return module.Storage.Get(ctx, id)
}
