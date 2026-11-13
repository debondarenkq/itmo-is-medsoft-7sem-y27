package modules

import (
	"context"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/catalog/internal/models"
)

type Storage interface {
	List(context.Context, models.Page) ([]models.Diagnosis, int, error)
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
