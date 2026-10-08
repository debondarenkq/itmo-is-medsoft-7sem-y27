package api

import (
	"github.com/google/uuid"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/api"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/internal/models"
)

func fromStaffInput(value contract.StaffInput) models.Staff {
	return models.Staff{FirstName: value.FirstName, LastName: value.LastName, Position: value.Position}
}
func toStaff(value models.Staff) contract.Staff {
	return contract.Staff{ID: uuid.MustParse(value.ID), FirstName: value.FirstName, LastName: value.LastName, Position: value.Position, CreatedAt: value.CreatedAt, UpdatedAt: value.UpdatedAt, DeletedAt: value.DeletedAt}
}
