package storage

import (
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/catalog/internal/models"
)

const (
	diagnosisCodeConstraint = "diagnoses_code_key"
	diagnosisNameConstraint = "diagnoses_name_key"
)

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
