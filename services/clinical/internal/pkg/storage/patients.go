package storage

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/models"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/pkg/platform"
)

type Repository struct{ pool *pgxpool.Pool }

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

const patientColumns = "p.id,p.first_name,p.last_name,p.middle_name,to_char(p.birth_date,'YYYY-MM-DD'),p.administrative_sex,p.comment,EXISTS(SELECT 1 FROM records r WHERE r.patient_id=p.id),p.created_at,p.updated_at"
const patientReturning = "id,first_name,last_name,middle_name,to_char(birth_date,'YYYY-MM-DD'),administrative_sex,comment,created_at,updated_at"

func scanPatient(row pgx.Row) (models.Patient, error) {
	var v models.Patient
	err := row.Scan(&v.ID, &v.FirstName, &v.LastName, &v.MiddleName, &v.BirthDate, &v.AdministrativeSex, &v.Comment, &v.HasRecord, &v.CreatedAt, &v.UpdatedAt)
	return v, err
}

func scanPatientWrite(row pgx.Row) (models.Patient, error) {
	var v models.Patient
	err := row.Scan(&v.ID, &v.FirstName, &v.LastName, &v.MiddleName, &v.BirthDate, &v.AdministrativeSex, &v.Comment, &v.CreatedAt, &v.UpdatedAt)
	return v, err
}

func (r *Repository) GetPatient(ctx context.Context, id string) (models.Patient, error) {
	return scanPatient(r.pool.QueryRow(ctx, "SELECT "+patientColumns+" FROM patients p WHERE p.id=$1", id))
}

func (r *Repository) ListPatients(ctx context.Context, p platform.Page) ([]models.Patient, int, error) {
	filter := " FROM patients p WHERE ($1='' OR (p.first_name || ' ' || p.last_name || ' ' || COALESCE(p.middle_name,'')) ILIKE '%' || $1 || '%')"
	var total int
	if err := r.pool.QueryRow(ctx, "SELECT count(*)"+filter, p.Query).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := r.pool.Query(ctx, "SELECT "+patientColumns+filter+" ORDER BY p.last_name,p.first_name,p.id LIMIT $2 OFFSET $3", p.Query, p.Limit, p.Offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	items := []models.Patient{}
	for rows.Next() {
		v, err := scanPatient(rows)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, v)
	}
	return items, total, rows.Err()
}

func (r *Repository) CreatePatient(ctx context.Context, v models.PatientInput) (models.Patient, error) {
	return scanPatientWrite(r.pool.QueryRow(ctx, "INSERT INTO patients(id,first_name,last_name,middle_name,birth_date,administrative_sex,comment) VALUES($1,$2,$3,$4,$5,$6,$7) RETURNING "+patientReturning, uuid.NewString(), v.FirstName, v.LastName, v.MiddleName, v.BirthDate, v.AdministrativeSex, v.Comment))
}

func (r *Repository) UpdatePatient(ctx context.Context, id string, v models.PatientInput) (models.Patient, error) {
	// Locking the patient also serializes this snapshot with record creation/deletion.
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return models.Patient{}, err
	}
	defer tx.Rollback(ctx)
	patient, err := scanPatientWrite(tx.QueryRow(ctx, "UPDATE patients SET first_name=$2,last_name=$3,middle_name=$4,birth_date=$5,administrative_sex=$6,comment=$7,updated_at=clock_timestamp() WHERE id=$1 RETURNING "+patientReturning, id, v.FirstName, v.LastName, v.MiddleName, v.BirthDate, v.AdministrativeSex, v.Comment))
	if err != nil {
		return patient, err
	}
	if err = tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM records WHERE patient_id=$1)", id).Scan(&patient.HasRecord); err != nil {
		return patient, err
	}
	return patient, tx.Commit(ctx)
}

func (r *Repository) DeletePatient(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, "DELETE FROM patients WHERE id=$1", id)
	var db *pgconn.PgError
	if errors.As(err, &db) && db.Code == "23503" {
		return platform.Conflict("PATIENT_HAS_RECORD", "Нельзя удалить пациента, для которого создана ЭМК")
	}
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
