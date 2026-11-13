package storage

import (
	"context"
	"errors"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/internal/models"
	"github.com/jackc/pgx/v5/pgconn"
)

func (r *Repository) CreateAccount(ctx context.Context, a models.Account, actor string) (models.Account, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return a, err
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, "INSERT INTO accounts(id,username,password_hash,staff_id) VALUES($1,$2,$3,$4)", a.ID, a.Username, a.PasswordHash, a.StaffID)
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
	if err = addAudit(ctx, tx, actor, "Создание аккаунта "+a.ID, models.AuditEventKindAccountCreated); err != nil {
		return a, err
	}
	return a, tx.Commit(ctx)
}
