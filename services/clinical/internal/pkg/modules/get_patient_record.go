package modules

import (
	"context"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/models"
)

func (module *Module) GetPatientRecord(ctx context.Context, id string) (models.Record, error) {
	return module.Storage.GetPatientRecord(ctx, id)
}
