package storage

import (
	"context"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/laboratory/internal/models"
	"github.com/google/uuid"
)

func (r *Repository) DeactivateTest(ctx context.Context, id uuid.UUID) (models.Test, error) {
	var t models.Test
	err := r.pool.QueryRow(ctx, "UPDATE lab_tests SET active=false WHERE id=$1 RETURNING id,code,name,active", id).Scan(&t.ID, &t.Code, &t.Name, &t.Active)
	return t, databaseError(err)
}
