package storage

import (
	"context"
	"encoding/json"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/laboratory/internal/models"
	"github.com/google/uuid"
)

func (r *Repository) CreateRequests(ctx context.Context, patient uuid.UUID, requests []models.Request) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	for _, request := range requests {
		data, err := json.Marshal(request)
		if err != nil {
			return err
		}
		if _, err = tx.Exec(ctx, "INSERT INTO lab_requests(id,patient_id,resource) VALUES($1,$2,$3)", request.ID, patient, data); err != nil {
			return databaseError(err)
		}
	}
	return tx.Commit(ctx)
}
