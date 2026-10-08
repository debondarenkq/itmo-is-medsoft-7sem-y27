package modules

import (
	"context"
	"strings"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/internal/models"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/internal/pkg/platform"
)

type Storage interface {
	List(context.Context, platform.Page) ([]models.Staff, int, error)
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

func Validate(v models.Staff) (models.Staff, error) {
	v.FirstName = strings.TrimSpace(v.FirstName)
	if err := platform.Text("first_name", v.FirstName, 3, 24); err != nil {
		return v, err
	}
	v.LastName = strings.TrimSpace(v.LastName)
	if err := platform.Text("last_name", v.LastName, 3, 24); err != nil {
		return v, err
	}
	v.Position = strings.TrimSpace(v.Position)
	if err := platform.Text("position", v.Position, 3, 24); err != nil {
		return v, err
	}
	return v, nil
}

func (m *Module) List(ctx context.Context, p platform.Page) ([]models.Staff, int, error) {
	return m.Storage.List(ctx, p)
}

func (m *Module) Get(ctx context.Context, id string) (models.Staff, error) {
	return m.Storage.Get(ctx, id)
}

func (m *Module) Create(ctx context.Context, v models.Staff) (models.Staff, error) {
	v, err := Validate(v)
	if err != nil {
		return v, err
	}
	return m.Storage.Create(ctx, v)
}

func (m *Module) Update(ctx context.Context, id string, v models.Staff) (models.Staff, error) {
	v, err := Validate(v)
	if err != nil {
		return v, err
	}
	old, err := m.Storage.Get(ctx, id)
	if err != nil {
		return v, err
	}
	if old.DeletedAt != nil {
		return v, platform.Conflict("DELETED", "Удалённую запись нельзя редактировать")
	}
	return m.Storage.Update(ctx, id, v)
}

func (m *Module) Delete(ctx context.Context, id string) error {
	return m.Storage.Delete(ctx, id)
}
