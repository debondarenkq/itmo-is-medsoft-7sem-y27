package modules

import (
	"context"
	"fmt"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/catalog/internal/models"
)

func (module *Module) Import(ctx context.Context, diagnoses []models.Diagnosis) ([]models.Diagnosis, error) {
	if len(diagnoses) == 0 || len(diagnoses) > 1000 {
		return nil, fmt.Errorf("%w: Файл должен содержать от 1 до 1000 диагнозов", models.ErrInvalidInput)
	}

	validated := make([]models.Diagnosis, len(diagnoses))
	codes := make(map[string]struct{}, len(diagnoses))
	names := make(map[string]struct{}, len(diagnoses))
	for index, diagnosis := range diagnoses {
		entry, err := Validate(diagnosis)
		if err != nil {
			return nil, fmt.Errorf("запись %d: %w", index+1, err)
		}
		if _, exists := codes[entry.Code]; exists {
			return nil, fmt.Errorf("%w: %s", models.ErrInvalidInput, "В файле повторяется код: "+entry.Code)
		}
		if _, exists := names[entry.Name]; exists {
			return nil, fmt.Errorf("%w: %s", models.ErrInvalidInput, "В файле повторяется наименование: "+entry.Name)
		}
		codes[entry.Code] = struct{}{}
		names[entry.Name] = struct{}{}
		validated[index] = entry
	}

	imported, err := module.Storage.UpsertBatch(ctx, validated)
	return imported, diagnosisWriteError(err)
}
