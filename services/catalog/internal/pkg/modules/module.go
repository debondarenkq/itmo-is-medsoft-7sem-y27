package modules

import (
	"context"
	"strings"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/catalog/internal/models"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/catalog/internal/pkg/platform"
)

type Storage interface {
	List(context.Context, platform.Page) ([]models.Diagnosis, int, error)
	Get(context.Context, string) (models.Diagnosis, error)
	Create(context.Context, models.Diagnosis) (models.Diagnosis, error)
	Update(context.Context, string, models.Diagnosis) (models.Diagnosis, error)
	Delete(context.Context, string) error
	UpsertBatch(context.Context, []models.Diagnosis) ([]models.Diagnosis, error)
}
type Deps struct{ Storage Storage }
type Module struct {
	Deps
}

func New(deps Deps) *Module {
	return &Module{Deps: deps}
}

func Validate(v models.Diagnosis) (models.Diagnosis, error) {
	v.Code = strings.TrimSpace(v.Code)
	if err := platform.Text("code", v.Code, 6, 6); err != nil {
		return v, err
	}
	v.Name = strings.TrimSpace(v.Name)
	if err := platform.Text("name", v.Name, 3, 24); err != nil {
		return v, err
	}
	return v, nil
}

func (m *Module) List(ctx context.Context, p platform.Page) ([]models.Diagnosis, int, error) {
	return m.Storage.List(ctx, p)
}

func (m *Module) Get(ctx context.Context, id string) (models.Diagnosis, error) {
	return m.Storage.Get(ctx, id)
}

func (m *Module) Create(ctx context.Context, v models.Diagnosis) (models.Diagnosis, error) {
	v, err := Validate(v)
	if err != nil {
		return v, err
	}
	return m.Storage.Create(ctx, v)
}

func (m *Module) Update(ctx context.Context, id string, v models.Diagnosis) (models.Diagnosis, error) {
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
