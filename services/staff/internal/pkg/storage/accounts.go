package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/internal/models"
)

const accountColumns = `a.id,a.username,a.password_hash,a.created_at,a.updated_at,a.is_active,a.staff_id,a.is_admin,COALESCE(a.totp_secret,''),a.last_totp_step,CASE WHEN a.is_admin THEN 'admin' ELSE COALESCE(s.role,'other') END,concat_ws(' ',s.last_name,s.first_name)`

const accountFrom = ` FROM accounts a LEFT JOIN staff s ON s.id=a.staff_id `

func scanAccount(row pgx.Row) (models.Account, error) {
	var a models.Account
	err := row.Scan(&a.ID, &a.Username, &a.PasswordHash, &a.CreatedAt, &a.UpdatedAt, &a.IsActive, &a.StaffID, &a.IsAdmin, &a.TOTPSecret, &a.LastTOTPStep, &a.Role, &a.StaffName)
	if errors.Is(err, pgx.ErrNoRows) {
		return a, models.ErrUnauthorized
	}
	return a, err
}

func (r *Repository) Account(ctx context.Context, username string) (models.Account, error) {
	return scanAccount(r.pool.QueryRow(ctx, "SELECT "+accountColumns+accountFrom+"WHERE a.username=$1", username))
}

func (r *Repository) CreateAccount(ctx context.Context, a models.Account, actor string) (models.Account, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return a, err
	}
	defer tx.Rollback(ctx)
	var deleted *time.Time
	err = tx.QueryRow(ctx, "SELECT deleted_at FROM staff WHERE id=$1 FOR UPDATE", a.StaffID).Scan(&deleted)
	if errors.Is(err, pgx.ErrNoRows) || deleted != nil {
		return a, models.ErrNotFound
	}
	if err != nil {
		return a, err
	}
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
	if err = addAudit(ctx, tx, actor, "Создание аккаунта "+a.ID, "account_created"); err != nil {
		return a, err
	}
	return a, tx.Commit(ctx)
}

func (r *Repository) ChangeAccount(ctx context.Context, id string, active *bool, passwordHash, actor string) (models.Account, error) {
	return r.UpdateAccount(ctx, id, models.AccountChanges{IsActive: active, PasswordHash: passwordHash}, actor)
}

func (r *Repository) UpdateAccount(ctx context.Context, id string, changes models.AccountChanges, actor string) (models.Account, error) {
	active, passwordHash := changes.IsActive, changes.PasswordHash
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return models.Account{}, err
	}
	defer tx.Rollback(ctx)
	// Lock staff first, in the same order as staff deletion and its trigger.
	_, err = tx.Exec(ctx, "SELECT id FROM staff WHERE id=(SELECT staff_id FROM accounts WHERE id=$1) FOR UPDATE", id)
	if err != nil {
		return models.Account{}, err
	}
	a, err := scanAccount(tx.QueryRow(ctx, "SELECT "+accountColumns+accountFrom+"WHERE a.id=$1 FOR UPDATE OF a", id))
	if err != nil {
		return a, models.ErrNotFound
	}
	if a.IsAdmin {
		return a, models.ErrForbidden
	}
	if active != nil && *active {
		var deleted *time.Time
		err = tx.QueryRow(ctx, "SELECT deleted_at FROM staff WHERE id=$1", a.StaffID).Scan(&deleted)
		if err != nil {
			return a, err
		}
		if deleted != nil {
			return a, models.ErrConflict
		}
	}
	if changes.Username != nil {
		a.Username = *changes.Username
	}
	if changes.Role != nil {
		tag, err := tx.Exec(ctx, "UPDATE staff SET role=$2,updated_at=clock_timestamp() WHERE id=$1 AND deleted_at IS NULL", a.StaffID, *changes.Role)
		if err != nil {
			return a, err
		}
		if tag.RowsAffected() != 1 {
			return a, models.ErrConflict
		}
	}
	if active != nil {
		a.IsActive = *active
	}
	if passwordHash != "" {
		a.PasswordHash = passwordHash
	}
	_, err = tx.Exec(ctx, "UPDATE accounts SET is_active=$2,password_hash=$3,username=$4,updated_at=clock_timestamp() WHERE id=$1", a.ID, a.IsActive, a.PasswordHash, a.Username)
	if err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) && pgError.Code == "23505" {
			return a, models.ErrConflict
		}
		return a, err
	}
	revoked := int64(0)
	if !a.IsActive || passwordHash != "" {
		tag, err := tx.Exec(ctx, "DELETE FROM refresh_tokens WHERE account_id=$1", a.ID)
		if err != nil {
			return a, err
		}
		revoked = tag.RowsAffected()
	}
	if err = addAudit(ctx, tx, actor, fmt.Sprintf("Изменение аккаунта %s; отозвано refresh токенов: %d", a.ID, revoked), "account_updated"); err != nil {
		return a, err
	}
	a, err = scanAccount(tx.QueryRow(ctx, "SELECT "+accountColumns+accountFrom+"WHERE a.id=$1", id))
	if err != nil {
		return a, err
	}
	return a, tx.Commit(ctx)
}

func (r *Repository) Accounts(ctx context.Context, limit, offset int) ([]models.Account, int, error) {
	return r.SearchAccounts(ctx, models.AccountFilter{Limit: limit, Offset: offset})
}

func (r *Repository) SearchAccounts(ctx context.Context, filter models.AccountFilter) ([]models.Account, int, error) {
	where := accountFrom + "WHERE ($1='' OR a.username ILIKE '%'||$1||'%' OR concat_ws(' ',s.last_name,s.first_name) ILIKE '%'||$1||'%') AND ($2::boolean IS NULL OR a.is_active=$2)"
	var total int
	if err := r.pool.QueryRow(ctx, "SELECT count(*)"+where, filter.Search, filter.IsActive).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := r.pool.Query(ctx, "SELECT "+accountColumns+where+" ORDER BY a.created_at,a.id LIMIT $3 OFFSET $4", filter.Search, filter.IsActive, filter.Limit, filter.Offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	items := []models.Account{}
	for rows.Next() {
		account, err := scanAccount(rows)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, account)
	}
	return items, total, rows.Err()
}
