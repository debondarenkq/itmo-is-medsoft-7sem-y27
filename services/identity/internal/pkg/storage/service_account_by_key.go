package storage

import (
	"context"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/internal/models"
)

func (r *Repository) ServiceAccountByKey(ctx context.Context, hash string) (models.Account, error) {
	return scanAccount(r.pool.QueryRow(ctx, "SELECT "+accountColumns+accountFrom+"WHERE a.is_service AND a.password_hash=$1", hash))
}
