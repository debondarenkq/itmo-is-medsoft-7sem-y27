package storage

import (
	"context"
	"errors"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/models"
	"github.com/jackc/pgx/v5"
)

const patientColumns = "p.id,p.first_name,p.last_name,p.middle_name,to_char(p.birth_date,'YYYY-MM-DD'),p.administrative_sex,p.comment,EXISTS(SELECT 1 FROM records r WHERE r.patient_id=p.id),p.created_at,p.updated_at,(SELECT r.id FROM records r WHERE r.patient_id=p.id),(SELECT r.number FROM records r WHERE r.patient_id=p.id)"

func (repository *Repository) GetPatient(ctx context.Context, id string) (models.Patient, error) {
	return scanPatient(repository.pool.QueryRow(ctx, "SELECT "+patientColumns+" FROM patients p WHERE p.id=$1", id))
}

func mapNotFoundError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return models.ErrNotFound
	}
	return err
}

func scanPatient(row pgx.Row) (models.Patient, error) {
	var patient models.Patient
	err := row.Scan(&patient.ID, &patient.FirstName, &patient.LastName, &patient.MiddleName, &patient.BirthDate, &patient.AdministrativeSex, &patient.Comment, &patient.HasRecord, &patient.CreatedAt, &patient.UpdatedAt, &patient.RecordID, &patient.RecordNumber)
	return patient, mapNotFoundError(err)
}
