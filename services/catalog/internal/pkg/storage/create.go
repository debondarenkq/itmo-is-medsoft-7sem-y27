package storage

import (
	"context"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/catalog/internal/models"
	"github.com/google/uuid"
)

func (repository *Repository) Create(ctx context.Context, diagnosis models.Diagnosis) (models.Diagnosis, error) {
	return scan(repository.pool.QueryRow(ctx, "INSERT INTO diagnoses(id,code,name) VALUES ($1, $2, $3) RETURNING "+columns, uuid.NewString(), diagnosis.Code, diagnosis.Name))
}
