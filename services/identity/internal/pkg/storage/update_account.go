package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/internal/models"
	"github.com/jackc/pgx/v5/pgconn"
)

func (r *Repository) UpdateAccount(ctx context.Context, id string, changes models.AccountChanges, actor string) (models.Account, error) {
	active, passwordHash := changes.IsActive, changes.PasswordHash
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return models.Account{}, err
	}
	defer tx.Rollback(ctx)
	a, err := scanAccount(tx.QueryRow(ctx, "SELECT "+accountColumns+accountFrom+"WHERE a.id=$1 FOR UPDATE OF a", id))
	if err != nil {
		return a, models.ErrNotFound
	}
	if a.IsAdmin {
		return a, models.ErrForbidden
	}
	if changes.Username != nil {
		a.Username = *changes.Username
	}

	if active != nil {
		a.IsActive = *active
	}
	if passwordHash != "" {
		a.PasswordHash = passwordHash
	}
	_, err = tx.Exec(ctx, "UPDATE accounts SET is_active=$2,password_hash=$3,username=$4,updated_at=clock_timestamp(),session_version=session_version+1 WHERE id=$1", a.ID, a.IsActive, a.PasswordHash, a.Username)
	if err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) && pgError.Code == "23505" {
			return a, models.ErrConflict
		}
		return a, err
	}
	revoked := int64(0)
	if !a.IsActive || passwordHash != "" || changes.RevokeSessions {
		tag, err := tx.Exec(ctx, "DELETE FROM refresh_tokens WHERE account_id=$1", a.ID)
		if err != nil {
			return a, err
		}
		revoked = tag.RowsAffected()
	}
	if err = addAudit(ctx, tx, actor, fmt.Sprintf("Изменение аккаунта %s; отозвано refresh токенов: %d", a.ID, revoked), models.AuditEventKindAccountUpdated); err != nil {
		return a, err
	}
	a, err = scanAccount(tx.QueryRow(ctx, "SELECT "+accountColumns+accountFrom+"WHERE a.id=$1", id))
	if err != nil {
		return a, err
	}
	return a, tx.Commit(ctx)
}
