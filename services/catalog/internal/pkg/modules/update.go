package modules

import (
	"context"
	"fmt"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/catalog/internal/models"
)

func (module *Module) Update(ctx context.Context, id string, diagnosis models.Diagnosis) (models.Diagnosis, error) {
	diagnosis, err := Validate(diagnosis)
	if err != nil {
		return diagnosis, err
	}
	existingDiagnosis, err := module.Storage.Get(ctx, id)
	if err != nil {
		return diagnosis, err
	}
	if existingDiagnosis.DeletedAt != nil {
		return diagnosis, fmt.Errorf("%w: Удалённую запись нельзя редактировать", models.ErrConflict)
	}
	updated, err := module.Storage.Update(ctx, id, diagnosis)
	return updated, diagnosisWriteError(err)
}
