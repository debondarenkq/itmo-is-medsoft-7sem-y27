package storage

import (
	"context"
	"errors"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/internal/models"
)

func (r *Repository) AccountByID(ctx context.Context, id string) (models.Account, error) {
	a, err := scanAccount(r.pool.QueryRow(ctx, "SELECT "+accountColumns+accountFrom+"WHERE a.id=$1", id))
	if errors.Is(err, models.ErrUnauthorized) {
		return a, models.ErrNotFound
	}
	return a, err
}
