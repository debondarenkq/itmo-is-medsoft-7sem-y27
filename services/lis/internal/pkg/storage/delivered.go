package storage

import (
	"context"

	"github.com/google/uuid"
)

func (r *Repository) Delivered(ctx context.Context, id uuid.UUID) error {
	_, err := r.Pool.Exec(ctx, "UPDATE results SET delivered=true WHERE id=$1", id)
	return err
}
