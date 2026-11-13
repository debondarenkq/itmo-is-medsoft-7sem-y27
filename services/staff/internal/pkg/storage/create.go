package storage

import (
	"context"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/internal/models"
	"github.com/google/uuid"
)

func (repository *Repository) Create(ctx context.Context, staff models.Staff) (models.Staff, error) {
	return scan(repository.pool.QueryRow(ctx, "INSERT INTO staff(id,first_name,last_name,position,role) VALUES ($1, $2, $3, $4, $5) RETURNING "+columns, uuid.NewString(), staff.FirstName, staff.LastName, staff.Position, staff.Role))
}
