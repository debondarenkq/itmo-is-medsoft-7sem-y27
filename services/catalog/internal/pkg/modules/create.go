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

func (module *Module) Create(ctx context.Context, diagnosis models.Diagnosis) (models.Diagnosis, error) {
	diagnosis, err := Validate(diagnosis)
	if err != nil {
		return diagnosis, err
	}
	created, err := module.Storage.Create(ctx, diagnosis)
	return created, diagnosisWriteError(err)
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
