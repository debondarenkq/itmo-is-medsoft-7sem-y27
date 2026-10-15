package api

import (
	"context"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/api"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/internal/validation"
)

func (implementation *Implementation) UpdateStaff(ctx context.Context, request contract.UpdateStaffRequestObject) (contract.UpdateStaffResponseObject, error) {
	if err := validation.RequireID(request.ID); err != nil {
		return nil, err
	}
	staff, err := implementation.Module.Update(ctx, request.ID.String(), fromStaffInput(*request.Body))
	if err != nil {
		return nil, err
	}
	response, err := toStaff(staff)
	if err != nil {
		return nil, err
	}
	return contract.UpdateStaff200JSONResponse(response), nil
}
