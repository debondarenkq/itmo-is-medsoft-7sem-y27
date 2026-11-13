package storage

import (
	"context"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/catalog/internal/models"
)

func (repository *Repository) Delete(ctx context.Context, id string) error {
	tag, err := repository.pool.Exec(ctx, "UPDATE diagnoses SET deleted_at=COALESCE(deleted_at,clock_timestamp()), updated_at=clock_timestamp() WHERE id=$1", id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return models.ErrNotFound
	}
	return nil
}
