package storage

import (
	"context"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

const patientReturning = "id,first_name,last_name,middle_name,to_char(birth_date,'YYYY-MM-DD'),administrative_sex,comment,created_at,updated_at"

func (repository *Repository) CreatePatient(ctx context.Context, patientData models.PatientData) (models.Patient, error) {
	return scanPatientWrite(repository.pool.QueryRow(ctx, "INSERT INTO patients(id,first_name,last_name,middle_name,birth_date,administrative_sex,comment) VALUES($1,$2,$3,$4,$5,$6,$7) RETURNING "+patientReturning, uuid.NewString(), patientData.FirstName, patientData.LastName, patientData.MiddleName, patientData.BirthDate, patientData.AdministrativeSex, patientData.Comment))
}

func scanPatientWrite(row pgx.Row) (models.Patient, error) {
	var patient models.Patient
	err := row.Scan(&patient.ID, &patient.FirstName, &patient.LastName, &patient.MiddleName, &patient.BirthDate, &patient.AdministrativeSex, &patient.Comment, &patient.CreatedAt, &patient.UpdatedAt)
	return patient, mapNotFoundError(err)
}
