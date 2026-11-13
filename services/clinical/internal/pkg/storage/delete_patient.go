package storage

import (
	"context"
	"errors"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/models"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

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
