package modules

import (
	"context"
	"errors"
	"fmt"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/catalog/internal/models"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/catalog/internal/pkg/platform"
)

func (m *Module) Import(ctx context.Context, diagnoses []models.Diagnosis) ([]models.Diagnosis, error) {
	if len(diagnoses) == 0 || len(diagnoses) > 1000 {
		return nil, platform.Bad("IMPORT_SIZE", "Файл должен содержать от 1 до 1000 диагнозов")
	}

	validated := make([]models.Diagnosis, len(diagnoses))
	codes := make(map[string]struct{}, len(diagnoses))
	names := make(map[string]struct{}, len(diagnoses))
	for index, diagnosis := range diagnoses {
		entry, err := Validate(diagnosis)
		if err != nil {
			var validation *platform.Error
			if errors.As(err, &validation) {
				fields := make(map[string]string, len(validation.Fields))
				for field, message := range validation.Fields {
					fields[fmt.Sprintf("diagnoses[%d].%s", index, field)] = message
				}
				return nil, &platform.Error{
					Status:  validation.Status,
					Code:    validation.Code,
					Message: fmt.Sprintf("Некорректный диагноз в записи %d", index+1),
					Fields:  fields,
				}
			}
			return nil, err
		}
		if _, exists := codes[entry.Code]; exists {
			return nil, platform.Bad("IMPORT_DUPLICATE_CODE", "В файле повторяется код: "+entry.Code)
		}
		if _, exists := names[entry.Name]; exists {
			return nil, platform.Bad("IMPORT_DUPLICATE_NAME", "В файле повторяется наименование: "+entry.Name)
		}
		codes[entry.Code] = struct{}{}
		names[entry.Name] = struct{}{}
		validated[index] = entry
	}

	return m.Storage.UpsertBatch(ctx, validated)
}
