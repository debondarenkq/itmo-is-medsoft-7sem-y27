package modules

import (
	"context"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/models"
)

func (module *Module) GetRecord(ctx context.Context, id string) (models.Record, error) {
	return module.Storage.GetRecord(ctx, id)
}
