package storage

import (
	"context"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/models"
)

func (repository *Repository) UpdatePatient(ctx context.Context, id string, patientData models.PatientData) (models.Patient, error) {
	// Locking the patient also serializes this snapshot with record creation/deletion.
	tx, err := repository.pool.Begin(ctx)
	if err != nil {
		return models.Patient{}, err
	}
	defer tx.Rollback(ctx)
	patient, err := scanPatientWrite(tx.QueryRow(ctx, "UPDATE patients SET first_name=$2,last_name=$3,middle_name=$4,birth_date=$5,administrative_sex=$6,comment=$7,updated_at=clock_timestamp() WHERE id=$1 RETURNING "+patientReturning, id, patientData.FirstName, patientData.LastName, patientData.MiddleName, patientData.BirthDate, patientData.AdministrativeSex, patientData.Comment))
	if err != nil {
		return patient, err
	}
	if err = tx.QueryRow(ctx, "SELECT (SELECT id FROM records WHERE patient_id=$1),(SELECT number FROM records WHERE patient_id=$1)", id).Scan(&patient.RecordID, &patient.RecordNumber); err != nil {
		return patient, err
	}
	patient.HasRecord = patient.RecordID != nil
	return patient, tx.Commit(ctx)
}
