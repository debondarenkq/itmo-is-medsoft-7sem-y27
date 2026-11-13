package storage

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/lis/internal/models"
	"github.com/jackc/pgx/v5/pgconn"
)

func (r *Repository) Record(ctx context.Context, h models.History) error {
	data, err := json.Marshal(h)
	if err != nil {
		return err
	}
	_, err = r.Pool.Exec(ctx, "INSERT INTO results(id,request_id,snapshot) VALUES($1,$2,$3)", h.ID, h.Request.ID, data)
	var pe *pgconn.PgError
	if errors.As(err, &pe) && pe.Code == "23505" {
		return models.ErrConflict
	}
	return err
}
