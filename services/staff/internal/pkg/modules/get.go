package modules

import (
	"context"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/internal/models"
)

func (module *Module) Get(ctx context.Context, id string) (models.Staff, error) {
	return module.Storage.Get(ctx, id)
}
