package storage

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/lis/internal/models"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) Connection(ctx context.Context) (models.Connection, error) {
	var c models.Connection
	var data []byte
	err := r.Pool.QueryRow(ctx, "SELECT session FROM connection WHERE id").Scan(&data)
	if errors.Is(err, pgx.ErrNoRows) {
		return c, models.ErrUnconfigured
	}
	if err != nil {
		return c, err
	}
	err = json.Unmarshal(data, &c.Session)
	return c, err
}
