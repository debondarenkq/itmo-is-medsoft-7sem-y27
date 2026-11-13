package storage

import (
	"context"
	"time"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/internal/models"
)

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
	if !a.IsActive || !a.UpdatedAt.Equal(expected.UpdatedAt) || a.PasswordHash != expected.PasswordHash || a.SessionVersion != expected.SessionVersion {
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
	_, err = tx.Exec(ctx, "INSERT INTO refresh_tokens(token_hash,account_id,expires_at,staff_revision) VALUES($1,$2,$3,$4)", newHash, a.ID, time.Now().Add(24*time.Hour), expected.StaffRevision)
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
	kind := models.AuditEventKindLogin
	if oldHash != "" {
		kind = models.AuditEventKindTokenRefreshed
	}
	if err = addAudit(ctx, tx, "account:"+a.ID, message, kind); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
