package modules

import (
	"context"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/internal/models"
)

func (module *Module) List(ctx context.Context, page models.Page) ([]models.Staff, int, error) {
	return module.Storage.List(ctx, page)
}
