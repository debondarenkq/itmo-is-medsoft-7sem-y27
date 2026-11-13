package storage

import (
	"context"
	"encoding/json"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/models"
	"github.com/jackc/pgx/v5"
)

const recordColumns = "id,patient_id,version,created_at,updated_at,state,number"

func (repository *Repository) GetRecord(ctx context.Context, id string) (models.Record, error) {
	return scanRecord(repository.pool.QueryRow(ctx, "SELECT "+recordColumns+" FROM records WHERE id=$1", id))
}

func scanRecord(row pgx.Row) (models.Record, error) {
	var record models.Record
	var state []byte
	if err := row.Scan(&record.ID, &record.PatientID, &record.Version, &record.CreatedAt, &record.UpdatedAt, &state, &record.Number); err != nil {
		return record, mapNotFoundError(err)
	}
	return record, json.Unmarshal(state, &record.State)
}
