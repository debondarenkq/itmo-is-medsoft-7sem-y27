package modules

import (
	"context"
	"fmt"
	"strings"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/internal/models"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/internal/validation"
)

const (
	staffNameMinLength     = 3
	staffNameMaxLength     = 24
	staffPositionMinLength = 3
	staffPositionMaxLength = 24
)

type Storage interface {
	List(context.Context, models.Page) ([]models.Staff, int, error)
	Get(context.Context, string) (models.Staff, error)
	Create(context.Context, models.Staff) (models.Staff, error)
	Update(context.Context, string, models.Staff) (models.Staff, error)
	Delete(context.Context, string) error
}
type Deps struct{ Storage Storage }
type Module struct {
	Deps
}

func New(deps Deps) *Module {
	return &Module{Deps: deps}
}

func Validate(staff models.Staff) (models.Staff, error) {
	staff.FirstName = strings.TrimSpace(staff.FirstName)
	if err := validation.Text("first_name", staff.FirstName, staffNameMinLength, staffNameMaxLength); err != nil {
		return staff, err
	}
	staff.LastName = strings.TrimSpace(staff.LastName)
	if err := validation.Text("last_name", staff.LastName, staffNameMinLength, staffNameMaxLength); err != nil {
		return staff, err
	}
	staff.Position = strings.TrimSpace(staff.Position)
	if err := validation.Text("position", staff.Position, staffPositionMinLength, staffPositionMaxLength); err != nil {
		return staff, err
	}
	return staff, nil
}

func (module *Module) List(ctx context.Context, page models.Page) ([]models.Staff, int, error) {
	return module.Storage.List(ctx, page)
}

func (module *Module) Get(ctx context.Context, id string) (models.Staff, error) {
	return module.Storage.Get(ctx, id)
}

func (module *Module) Create(ctx context.Context, staff models.Staff) (models.Staff, error) {
	staff, err := Validate(staff)
	if err != nil {
		return staff, err
	}
	return module.Storage.Create(ctx, staff)
}

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

func (module *Module) Delete(ctx context.Context, id string) error {
	return module.Storage.Delete(ctx, id)
}
