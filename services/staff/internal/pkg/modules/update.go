package modules

import (
	"context"
	"fmt"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/internal/models"
)

func (module *Module) Update(ctx context.Context, id string, staff models.Staff) (models.Staff, error) {
	staff, err := Validate(staff)
	if err != nil {
		return staff, err
	}
	existingStaff, err := module.Storage.Get(ctx, id)
	if err != nil {
		return staff, err
	}
	if existingStaff.DeletedAt != nil {
		return staff, fmt.Errorf("%w: Удалённую запись нельзя редактировать", models.ErrConflict)
	}
	return module.Storage.Update(ctx, id, staff)
}
