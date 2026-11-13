package storage

import (
	"context"
	"errors"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/internal/models"
	"github.com/jackc/pgx/v5"
)

const accountColumns = `a.id,a.username,a.password_hash,a.created_at,a.updated_at,a.is_active,a.staff_id,a.is_admin,COALESCE(a.totp_secret,''),a.last_totp_step,a.session_version,a.is_service`

const accountFrom = ` FROM accounts a `

func (r *Repository) Account(ctx context.Context, username string) (models.Account, error) {
	return scanAccount(r.pool.QueryRow(ctx, "SELECT "+accountColumns+accountFrom+"WHERE a.username=$1 AND NOT a.is_service", username))
}

func scanAccount(row pgx.Row, extra ...any) (models.Account, error) {
	var a models.Account
	fields := []any{&a.ID, &a.Username, &a.PasswordHash, &a.CreatedAt, &a.UpdatedAt, &a.IsActive, &a.StaffID, &a.IsAdmin, &a.TOTPSecret, &a.LastTOTPStep, &a.SessionVersion, &a.IsService}
	err := row.Scan(append(fields, extra...)...)
	if errors.Is(err, pgx.ErrNoRows) {
		return a, models.ErrUnauthorized
	}
	if a.IsService {
		a.Role = models.ServiceRole
	}
	if a.IsAdmin {
		a.Role = "admin"
	}
	return a, err
}
