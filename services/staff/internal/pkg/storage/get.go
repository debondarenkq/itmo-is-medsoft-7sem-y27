package storage

import (
	"context"
	"errors"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/internal/models"
	"github.com/jackc/pgx/v5"
)

const columns = "id, first_name, last_name, position, role, created_at, updated_at, deleted_at, role_revision"

func (repository *Repository) Get(ctx context.Context, id string) (models.Staff, error) {
	return scan(repository.pool.QueryRow(ctx, "SELECT "+columns+" FROM staff WHERE id=$1", id))
}

func scan(row pgx.Row) (models.Staff, error) {
	var staff models.Staff
	err := row.Scan(&staff.ID, &staff.FirstName, &staff.LastName, &staff.Position, &staff.Role, &staff.CreatedAt, &staff.UpdatedAt, &staff.DeletedAt, &staff.RoleRevision)
	if errors.Is(err, pgx.ErrNoRows) {
		return staff, models.ErrNotFound
	}
	return staff, err
}
