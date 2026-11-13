package storage

import (
	"context"
	"errors"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/catalog/internal/models"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	diagnosisCodeConstraint = "diagnoses_code_key"
	diagnosisNameConstraint = "diagnoses_name_key"
)

const columns = "id, code, name, created_at, updated_at, deleted_at"

func (repository *Repository) Get(ctx context.Context, id string) (models.Diagnosis, error) {
	return scan(repository.pool.QueryRow(ctx, "SELECT "+columns+" FROM diagnoses WHERE id=$1", id))
}

func mapDiagnosisError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return models.ErrNotFound
	}
	var databaseError *pgconn.PgError
	if !errors.As(err, &databaseError) || databaseError.Code != pgerrcode.UniqueViolation {
		return err
	}
	switch databaseError.ConstraintName {
	case diagnosisCodeConstraint:
		return models.ErrDuplicateCode
	case diagnosisNameConstraint:
		return models.ErrDuplicateName
	default:
		return err
	}
}

func scan(row pgx.Row) (models.Diagnosis, error) {
	var diagnosis models.Diagnosis
	err := row.Scan(&diagnosis.ID, &diagnosis.Code, &diagnosis.Name, &diagnosis.CreatedAt, &diagnosis.UpdatedAt, &diagnosis.DeletedAt)
	return diagnosis, mapDiagnosisError(err)
}
