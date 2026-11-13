package modules

import (
	"context"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/internal/models"
)

func (m *Module) UpdateRole(ctx context.Context, id, role string) (models.Staff, error) {
	if err := validateRole(role); err != nil {
		return models.Staff{}, err
	}
	s, err := m.Get(ctx, id)
	if err != nil {
		return s, err
	}
	if s.DeletedAt != nil {
		return s, models.ErrConflict
	}
	return m.Storage.UpdateRole(ctx, id, role)
}
