package storage

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/lis/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) ByRequest(ctx context.Context, id uuid.UUID) (*models.History, error) {
	var data []byte
	var delivered bool
	err := r.Pool.QueryRow(ctx, "SELECT snapshot,delivered FROM results WHERE request_id=$1", id).Scan(&data, &delivered)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var h models.History
	if err = json.Unmarshal(data, &h); err != nil {
		return nil, err
	}
	h.Delivered = delivered
	return &h, nil
}
