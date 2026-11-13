package storage

import (
	"context"
)

func (r *Repository) SetupSecret(ctx context.Context, candidate string) (bool, string, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return false, "", err
	}
	defer tx.Rollback(ctx)
	var exists bool
	err = tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM accounts WHERE is_admin)").Scan(&exists)
	if err != nil || exists {
		return exists, "", err
	}
	_, err = tx.Exec(ctx, "INSERT INTO security_settings(id,setup_secret) VALUES(true,$1) ON CONFLICT DO NOTHING", candidate)
	if err != nil {
		return false, "", err
	}
	var secret string
	err = tx.QueryRow(ctx, "SELECT setup_secret FROM security_settings WHERE id").Scan(&secret)
	if err != nil {
		return false, "", err
	}
	return false, secret, tx.Commit(ctx)
}
