package storage

import (
	"context"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/catalog/internal/models"
)

func (repository *Repository) Update(ctx context.Context, id string, diagnosis models.Diagnosis) (models.Diagnosis, error) {
	return scan(repository.pool.QueryRow(ctx, "UPDATE diagnoses SET code=$2, name=$3, updated_at=clock_timestamp() WHERE id=$1 AND deleted_at IS NULL RETURNING "+columns, id, diagnosis.Code, diagnosis.Name))
}
