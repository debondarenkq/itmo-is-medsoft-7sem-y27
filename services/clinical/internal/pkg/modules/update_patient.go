package modules

import (
	"context"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/models"
)

func (module *Module) UpdatePatient(ctx context.Context, id string, patientData models.PatientData) (models.Patient, error) {
	patientData, err := ValidatePatient(patientData)
	if err != nil {
		return models.Patient{}, err
	}
	return module.Storage.UpdatePatient(ctx, id, patientData)
}
