package storage

import (
	"context"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/internal/models"
)

func (repository *Repository) Update(ctx context.Context, id string, staff models.Staff) (models.Staff, error) {
	return scan(repository.pool.QueryRow(ctx, "UPDATE staff SET first_name=$2, last_name=$3, position=$4, role_revision=role_revision+CASE WHEN role IS DISTINCT FROM $5 THEN 1 ELSE 0 END, role=$5, updated_at=clock_timestamp() WHERE id=$1 AND deleted_at IS NULL RETURNING "+columns, id, staff.FirstName, staff.LastName, staff.Position, staff.Role))
}
