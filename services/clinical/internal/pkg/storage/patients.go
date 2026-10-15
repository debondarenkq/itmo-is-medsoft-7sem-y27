package storage

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/models"
)

type Repository struct{ pool *pgxpool.Pool }

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

const patientColumns = "p.id,p.first_name,p.last_name,p.middle_name,to_char(p.birth_date,'YYYY-MM-DD'),p.administrative_sex,p.comment,EXISTS(SELECT 1 FROM records r WHERE r.patient_id=p.id),p.created_at,p.updated_at"
const patientReturning = "id,first_name,last_name,middle_name,to_char(birth_date,'YYYY-MM-DD'),administrative_sex,comment,created_at,updated_at"

func scanPatient(row pgx.Row) (models.Patient, error) {
	var patient models.Patient
	err := row.Scan(&patient.ID, &patient.FirstName, &patient.LastName, &patient.MiddleName, &patient.BirthDate, &patient.AdministrativeSex, &patient.Comment, &patient.HasRecord, &patient.CreatedAt, &patient.UpdatedAt)
	return patient, mapNotFoundError(err)
}

func scanPatientWrite(row pgx.Row) (models.Patient, error) {
	var patient models.Patient
	err := row.Scan(&patient.ID, &patient.FirstName, &patient.LastName, &patient.MiddleName, &patient.BirthDate, &patient.AdministrativeSex, &patient.Comment, &patient.CreatedAt, &patient.UpdatedAt)
	return patient, mapNotFoundError(err)
}

func (repository *Repository) GetPatient(ctx context.Context, id string) (models.Patient, error) {
	return scanPatient(repository.pool.QueryRow(ctx, "SELECT "+patientColumns+" FROM patients p WHERE p.id=$1", id))
}

func (repository *Repository) ListPatients(ctx context.Context, page models.Page) ([]models.Patient, int, error) {
	filter := " FROM patients p WHERE ($1='' OR (p.first_name || ' ' || p.last_name || ' ' || COALESCE(p.middle_name,'')) ILIKE '%' || $1 || '%')"
	var total int
	if err := repository.pool.QueryRow(ctx, "SELECT count(*)"+filter, page.Search).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := repository.pool.Query(ctx, "SELECT "+patientColumns+filter+" ORDER BY p.last_name,p.first_name,p.id LIMIT $2 OFFSET $3", page.Search, page.Limit, page.Offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	items := []models.Patient{}
	for rows.Next() {
		patient, err := scanPatient(rows)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, patient)
	}
	return items, total, rows.Err()
}

func (repository *Repository) CreatePatient(ctx context.Context, patientData models.PatientData) (models.Patient, error) {
	return scanPatientWrite(repository.pool.QueryRow(ctx, "INSERT INTO patients(id,first_name,last_name,middle_name,birth_date,administrative_sex,comment) VALUES($1,$2,$3,$4,$5,$6,$7) RETURNING "+patientReturning, uuid.NewString(), patientData.FirstName, patientData.LastName, patientData.MiddleName, patientData.BirthDate, patientData.AdministrativeSex, patientData.Comment))
}

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
	if err = tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM records WHERE patient_id=$1)", id).Scan(&patient.HasRecord); err != nil {
		return patient, err
	}
	return patient, tx.Commit(ctx)
}

func (repository *Repository) DeletePatient(ctx context.Context, id string) error {
	tag, err := repository.pool.Exec(ctx, "DELETE FROM patients WHERE id=$1", id)
	var databaseError *pgconn.PgError
	if errors.As(err, &databaseError) && databaseError.Code == pgerrcode.ForeignKeyViolation && databaseError.ConstraintName == recordPatientForeignKey {
		return models.ErrReferenced
	}
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return models.ErrNotFound
	}
	return nil
}
