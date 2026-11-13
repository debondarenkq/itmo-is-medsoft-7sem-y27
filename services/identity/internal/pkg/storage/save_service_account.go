package storage

import (
	"context"
	"errors"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/internal/models"
	"github.com/jackc/pgx/v5/pgconn"
)

func (r *Repository) SaveServiceAccount(ctx context.Context, a models.Account, create bool, actor string) (models.Account, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return a, err
	}
	defer tx.Rollback(ctx)
	if create {
		_, err = tx.Exec(ctx, "INSERT INTO accounts(id,username,password_hash,is_service) VALUES($1,$2,$3,true)", a.ID, a.Username, a.PasswordHash)
	} else {
		tag, e := tx.Exec(ctx, "UPDATE accounts SET username=$2,is_active=$3,password_hash=CASE WHEN $4='' THEN password_hash ELSE $4 END,session_version=session_version+1,updated_at=clock_timestamp() WHERE id=$1 AND is_service", a.ID, a.Username, a.IsActive, a.PasswordHash)
		err = e
		if err == nil && tag.RowsAffected() != 1 {
			return a, models.ErrNotFound
		}
		if err == nil {
			_, err = tx.Exec(ctx, "DELETE FROM refresh_tokens WHERE account_id=$1", a.ID)
		}
	}
	if err != nil {
		var pe *pgconn.PgError
		if errors.As(err, &pe) && pe.Code == "23505" {
			return a, models.ErrConflict
		}
		return a, err
	}
	a, err = scanAccount(tx.QueryRow(ctx, "SELECT "+accountColumns+accountFrom+"WHERE a.id=$1", a.ID))
	if err != nil {
		return a, err
	}
	if err = addAudit(ctx, tx, actor, "Сервисный аккаунт "+a.ID, models.AuditEventKindOther); err != nil {
		return a, err
	}
	return a, tx.Commit(ctx)
}
