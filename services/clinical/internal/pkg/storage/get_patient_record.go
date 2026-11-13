package storage

import (
	"context"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/models"
)

func (repository *Repository) GetPatientRecord(ctx context.Context, id string) (models.Record, error) {
	return scanRecord(repository.pool.QueryRow(ctx, "SELECT "+recordColumns+" FROM records WHERE patient_id=$1", id))
}
