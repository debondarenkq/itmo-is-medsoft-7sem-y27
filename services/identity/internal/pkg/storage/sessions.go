package storage

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/internal/models"
)

func (r *Repository) RefreshAccount(ctx context.Context, hash string) (models.Account, error) {
	return scanAccount(r.pool.QueryRow(ctx, "SELECT "+accountColumns+accountFrom+"JOIN refresh_tokens t ON t.account_id=a.id WHERE t.token_hash=$1 AND t.expires_at>clock_timestamp()", hash))
}

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
	_, err = tx.Exec(ctx, "INSERT INTO refresh_tokens VALUES($1,$2,$3)", refreshHash, a.ID, time.Now().Add(24*time.Hour))
	if err != nil {
		return err
	}
	if err = addAudit(ctx, tx, "account:"+a.ID, "Настройка администратора и TOTP; вход администратора", "admin_setup"); err != nil {
		return err
	}
	_, err = tx.Exec(ctx, "DELETE FROM security_settings")
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *Repository) Issue(ctx context.Context, expected models.Account, oldHash, newHash string, step int64) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	a, err := scanAccount(tx.QueryRow(ctx, "SELECT "+accountColumns+accountFrom+"WHERE a.id=$1 FOR UPDATE OF a", expected.ID))
	if err != nil {
		return err
	}
	if !a.IsActive || !a.UpdatedAt.Equal(expected.UpdatedAt) || a.PasswordHash != expected.PasswordHash || a.Role != expected.Role {
		return models.ErrUnauthorized
	}
	if oldHash != "" {
		tag, err := tx.Exec(ctx, "DELETE FROM refresh_tokens WHERE token_hash=$1 AND account_id=$2 AND expires_at>clock_timestamp()", oldHash, a.ID)
		if err != nil {
			return err
		}
		if tag.RowsAffected() != 1 {
			return models.ErrUnauthorized
		}
	}
	if step >= 0 {
		if step <= a.LastTOTPStep {
			return models.ErrUnauthorized
		}
		_, err = tx.Exec(ctx, "UPDATE accounts SET last_totp_step=$2 WHERE id=$1", a.ID, step)
		if err != nil {
			return err
		}
	}
	_, err = tx.Exec(ctx, "INSERT INTO refresh_tokens VALUES($1,$2,$3)", newHash, a.ID, time.Now().Add(24*time.Hour))
	if err != nil {
		return err
	}
	message := "Вход пользователя"
	if a.IsAdmin {
		message = "Вход администратора"
	}
	if oldHash != "" {
		message = "Обновление пары токенов"
	}
	kind := "login"
	if oldHash != "" {
		kind = "token_refreshed"
	}
	if err = addAudit(ctx, tx, "account:"+a.ID, message, kind); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

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
	if err = addAudit(ctx, tx, "account:"+accountID, "Выход из системы", "logout"); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
