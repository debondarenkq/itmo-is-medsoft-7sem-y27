package modules

import (
	"context"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/models"
)

func (module *Module) ListPatients(ctx context.Context, page models.Page) ([]models.Patient, int, error) {
	return module.Storage.ListPatients(ctx, page)
}
