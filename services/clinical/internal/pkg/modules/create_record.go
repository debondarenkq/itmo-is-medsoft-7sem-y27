package modules

import (
	"context"
	"errors"
	"fmt"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/models"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/validation"
)

func (module *Module) CreateRecord(ctx context.Context, patientID, staffID string) (models.Record, error) {
	staffID, err := validation.ParseID(staffID)
	if err != nil {
		return models.Record{}, fmt.Errorf("%w: Укажите медработника, создающего ЭМК", models.ErrInvalidInput)
	}
	patient, err := module.Storage.GetPatient(ctx, patientID)
	if err != nil {
		return models.Record{}, err
	}
	if patient.HasRecord {
		return models.Record{}, errRecordExists
	}
	actor, err := module.Staff.GetActor(ctx, staffID)
	if err != nil {
		return models.Record{}, err
	}
	record, err := module.Storage.CreateRecord(ctx, patientID, actor)
	if errors.Is(err, models.ErrAlreadyExists) {
		return models.Record{}, errRecordExists
	}
	return record, err
}
