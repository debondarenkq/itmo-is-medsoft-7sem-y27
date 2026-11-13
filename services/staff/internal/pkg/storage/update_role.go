package storage

import (
	"context"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/internal/models"
)

func (r *Repository) UpdateRole(ctx context.Context, id, role string) (models.Staff, error) {
	return scan(r.pool.QueryRow(ctx, "UPDATE staff SET role_revision=role_revision+CASE WHEN role IS DISTINCT FROM $2 THEN 1 ELSE 0 END,role=$2,updated_at=clock_timestamp() WHERE id=$1 AND deleted_at IS NULL RETURNING "+columns, id, role))
}
