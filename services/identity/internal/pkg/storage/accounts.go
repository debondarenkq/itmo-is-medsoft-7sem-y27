package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/internal/models"
)

const accountColumns = `a.id,a.username,a.password_hash,a.created_at,a.updated_at,a.is_active,a.staff_id,a.is_admin,COALESCE(a.totp_secret,''),a.last_totp_step,a.session_version`

const accountFrom = ` FROM accounts a `

func scanAccount(row pgx.Row, extra ...any) (models.Account, error) {
	var a models.Account
	fields := []any{&a.ID, &a.Username, &a.PasswordHash, &a.CreatedAt, &a.UpdatedAt, &a.IsActive, &a.StaffID, &a.IsAdmin, &a.TOTPSecret, &a.LastTOTPStep, &a.SessionVersion}
	err := row.Scan(append(fields, extra...)...)
	if errors.Is(err, pgx.ErrNoRows) {
		return a, models.ErrUnauthorized
	}
	if a.IsAdmin {
		a.Role = "admin"
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

func (r *Repository) Accounts(ctx context.Context, limit, offset int) ([]models.Account, int, error) {
	return r.SearchAccounts(ctx, models.AccountFilter{Limit: limit, Offset: offset})
}

func (r *Repository) SearchAccounts(ctx context.Context, filter models.AccountFilter) ([]models.Account, int, error) {
	where := accountFrom + "WHERE ($1='' OR a.username ILIKE '%'||$1||'%' OR a.staff_id=ANY($3::uuid[])) AND ($2::boolean IS NULL OR (a.is_active AND NOT COALESCE(a.staff_id=ANY($4::uuid[]),false))=$2)"
	var total int
	if err := r.pool.QueryRow(ctx, "SELECT count(*)"+where, filter.Search, filter.IsActive, filter.StaffIDs, filter.DeletedStaffIDs).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := r.pool.Query(ctx, "SELECT "+accountColumns+where+" ORDER BY a.created_at,a.id LIMIT $5 OFFSET $6", filter.Search, filter.IsActive, filter.StaffIDs, filter.DeletedStaffIDs, filter.Limit, filter.Offset)
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

func (r *Repository) AccountByID(ctx context.Context, id string) (models.Account, error) {
	a, err := scanAccount(r.pool.QueryRow(ctx, "SELECT "+accountColumns+accountFrom+"WHERE a.id=$1", id))
	if errors.Is(err, models.ErrUnauthorized) {
		return a, models.ErrNotFound
	}
	return a, err
}
