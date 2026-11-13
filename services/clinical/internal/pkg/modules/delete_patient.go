package modules

import (
	"context"
	"errors"
	"fmt"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/models"
)

var (
	errPatientHasRecord = fmt.Errorf("%w: Нельзя удалить пациента, для которого создана ЭМК", models.ErrConflict)
	errRecordExists     = fmt.Errorf("%w: У пациента уже есть ЭМК", models.ErrConflict)
)

func (module *Module) DeletePatient(ctx context.Context, id string) error {
	patient, err := module.Storage.GetPatient(ctx, id)
	if err != nil {
		return err
	}
	if patient.HasRecord {
		return errPatientHasRecord
	}
	err = module.Storage.DeletePatient(ctx, id)
	if errors.Is(err, models.ErrReferenced) {
		return errPatientHasRecord
	}
	return err
}
