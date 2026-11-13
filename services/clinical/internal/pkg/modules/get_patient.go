package modules

import (
	"context"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/models"
)

func (module *Module) GetPatient(ctx context.Context, id string) (models.Patient, error) {
	return module.Storage.GetPatient(ctx, id)
}
