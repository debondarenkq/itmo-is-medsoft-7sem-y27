package storage

import (
	"context"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/internal/models"
)

func (r *Repository) Logout(ctx context.Context, accountID, hash string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, "DELETE FROM refresh_tokens WHERE account_id=$1 AND token_hash=$2", accountID, hash)
	if err != nil {
		return err
	}
	if err = addAudit(ctx, tx, "account:"+accountID, "Выход из системы", models.AuditEventKindLogout); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
