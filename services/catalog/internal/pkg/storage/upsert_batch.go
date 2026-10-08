package storage

import (
	"context"

	"github.com/google/uuid"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/catalog/internal/models"
)

func (r *Repository) UpsertBatch(ctx context.Context, diagnoses []models.Diagnosis) ([]models.Diagnosis, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	result := []models.Diagnosis{}
	for _, diagnosis := range diagnoses {
		item, err := scan(tx.QueryRow(ctx, "INSERT INTO diagnoses(id,code,name) VALUES($1,$2,$3) ON CONFLICT(code) DO UPDATE SET name=EXCLUDED.name,deleted_at=NULL,updated_at=clock_timestamp() RETURNING "+columns, uuid.NewString(), diagnosis.Code, diagnosis.Name))
		if err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}
	return result, nil
}
