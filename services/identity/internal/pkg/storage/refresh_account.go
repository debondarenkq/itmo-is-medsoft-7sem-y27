package storage

import (
	"context"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/internal/models"
)

func (r *Repository) RefreshAccount(ctx context.Context, hash string) (models.Account, error) {
	var revision int64
	a, err := scanAccount(r.pool.QueryRow(ctx, "SELECT "+accountColumns+",t.staff_revision"+accountFrom+"JOIN refresh_tokens t ON t.account_id=a.id WHERE t.token_hash=$1 AND t.expires_at>clock_timestamp()", hash), &revision)
	if err != nil {
		return a, err
	}
	a.StaffRevision = revision
	return a, nil
}
