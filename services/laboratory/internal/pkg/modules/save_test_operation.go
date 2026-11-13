package modules

import (
	"context"
	"regexp"
	"strings"
	"unicode/utf8"

	api "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/laboratory/api"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/laboratory/internal/models"
	"github.com/google/uuid"
)

var loincCode = regexp.MustCompile(`^\d{1,8}-\d$`)

func (m *Module) SaveTest(ctx context.Context, id uuid.UUID, input api.LabTestInput, create bool) (models.Test, error) {
	input.Name = strings.TrimSpace(input.Name)
	const minNameLength, maxNameLength = 3, 160
	if !loincCode.MatchString(input.Code) || utf8.RuneCountInString(input.Name) < minNameLength || utf8.RuneCountInString(input.Name) > maxNameLength {
		return models.Test{}, models.ErrInvalidInput
	}
	if create {
		id = uuid.New()
	}
	return m.Store.SaveTest(ctx, models.Test{ID: id, Code: input.Code, Name: input.Name, Active: true}, create)
}
