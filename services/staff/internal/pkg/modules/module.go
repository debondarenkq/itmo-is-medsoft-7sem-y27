package modules

import (
	"context"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/internal/models"
)

type Storage interface {
	List(context.Context, models.Page) ([]models.Staff, int, error)
	Get(context.Context, string) (models.Staff, error)
	Create(context.Context, models.Staff) (models.Staff, error)
	Update(context.Context, string, models.Staff) (models.Staff, error)
	Delete(context.Context, string) error
	ListIdentityProfiles(context.Context) ([]models.Staff, error)
	UpdateRole(context.Context, string, string) (models.Staff, error)
}

type Deps struct{ Storage Storage }

type Module struct {
	Deps
}

func New(deps Deps) *Module {
	return &Module{Deps: deps}
}
