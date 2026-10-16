package modules

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/catalog/internal/models"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/catalog/internal/validation"
)

const (
	diagnosisCodeLength    = 6
	diagnosisNameMinLength = 3
	diagnosisNameMaxLength = 24
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

func Validate(diagnosis models.Diagnosis) (models.Diagnosis, error) {
	diagnosis.Code = strings.TrimSpace(diagnosis.Code)
	if err := validation.Text("code", diagnosis.Code, diagnosisCodeLength, diagnosisCodeLength); err != nil {
		return diagnosis, err
	}
	diagnosis.Name = strings.TrimSpace(diagnosis.Name)
	if err := validation.Text("name", diagnosis.Name, diagnosisNameMinLength, diagnosisNameMaxLength); err != nil {
		return diagnosis, err
	}
	return diagnosis, nil
}

func (module *Module) List(ctx context.Context, page models.Page) ([]models.Diagnosis, int, error) {
	return module.Storage.List(ctx, page)
}

func (module *Module) Get(ctx context.Context, id string) (models.Diagnosis, error) {
	return module.Storage.Get(ctx, id)
}

func (module *Module) Create(ctx context.Context, diagnosis models.Diagnosis) (models.Diagnosis, error) {
	diagnosis, err := Validate(diagnosis)
	if err != nil {
		return diagnosis, err
	}
	created, err := module.Storage.Create(ctx, diagnosis)
	return created, diagnosisWriteError(err)
}

func (module *Module) Update(ctx context.Context, id string, diagnosis models.Diagnosis) (models.Diagnosis, error) {
	diagnosis, err := Validate(diagnosis)
	if err != nil {
		return diagnosis, err
	}
	existingDiagnosis, err := module.Storage.Get(ctx, id)
	if err != nil {
		return diagnosis, err
	}
	if existingDiagnosis.DeletedAt != nil {
		return diagnosis, fmt.Errorf("%w: Удалённую запись нельзя редактировать", models.ErrConflict)
	}
	updated, err := module.Storage.Update(ctx, id, diagnosis)
	return updated, diagnosisWriteError(err)
}

func (module *Module) Delete(ctx context.Context, id string) error {
	return module.Storage.Delete(ctx, id)
}

func diagnosisWriteError(err error) error {
	switch {
	case errors.Is(err, models.ErrDuplicateCode):
		return fmt.Errorf("%w: Диагноз с таким кодом уже существует", models.ErrConflict)
	case errors.Is(err, models.ErrDuplicateName):
		return fmt.Errorf("%w: Диагноз с таким наименованием уже существует", models.ErrConflict)
	default:
		return err
	}
}
