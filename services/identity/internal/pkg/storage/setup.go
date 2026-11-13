package storage

import (
	"context"
	"errors"
	"time"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/internal/models"
	"github.com/jackc/pgx/v5/pgconn"
)

func (r *Repository) Setup(ctx context.Context, a models.Account, secret, refreshHash string, step int64) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	var stored string
	err = tx.QueryRow(ctx, "SELECT setup_secret FROM security_settings WHERE id FOR UPDATE").Scan(&stored)
	if err != nil || stored != secret {
		return models.ErrConflict
	}
	_, err = tx.Exec(ctx, "INSERT INTO accounts(id,username,password_hash,is_admin,totp_secret,last_totp_step) VALUES($1,'admin',$2,true,$3,$4)", a.ID, a.PasswordHash, secret, step)
	if err != nil {
		var pe *pgconn.PgError
		if errors.As(err, &pe) && pe.Code == "23505" {
			return models.ErrConflict
		}
		return err
	}
	_, err = tx.Exec(ctx, "INSERT INTO refresh_tokens(token_hash,account_id,expires_at) VALUES($1,$2,$3)", refreshHash, a.ID, time.Now().Add(24*time.Hour))
	if err != nil {
		return err
	}
	if err = addAudit(ctx, tx, "account:"+a.ID, "Настройка администратора и TOTP; вход администратора", models.AuditEventKindAdminSetup); err != nil {
		return err
	}
	_, err = tx.Exec(ctx, "DELETE FROM security_settings")
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}
