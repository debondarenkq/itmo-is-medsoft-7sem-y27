package modules

import (
	"context"
	"fmt"
	"strings"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/internal/models"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/internal/validation"
)

const (
	staffNameMinLength     = 1
	staffNameMaxLength     = 24
	staffPositionMinLength = 3
	staffPositionMaxLength = 24
)

func (module *Module) Create(ctx context.Context, staff models.Staff) (models.Staff, error) {
	staff, err := Validate(staff)
	if err != nil {
		return staff, err
	}
	return module.Storage.Create(ctx, staff)
}

func Validate(staff models.Staff) (models.Staff, error) {
	if staff.Role == "" {
		staff.Role = "other"
	}
	if err := validateRole(staff.Role); err != nil {
		return staff, err
	}
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

func validateRole(role string) error {
	switch role {
	case "doctor", "nurse", "laboratory-technician", "other":
		return nil
	default:
		return fmt.Errorf("%w: Неизвестная роль", models.ErrInvalidInput)
	}
}
