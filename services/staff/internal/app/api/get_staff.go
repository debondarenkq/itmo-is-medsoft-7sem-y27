package api

import (
	"context"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/api"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/internal/validation"
)

func (implementation *Implementation) GetStaff(ctx context.Context, request contract.GetStaffRequestObject) (contract.GetStaffResponseObject, error) {
	if err := validation.RequireID(request.ID); err != nil {
		return nil, err
	}
	staff, err := implementation.Module.Get(ctx, request.ID.String())
	if err != nil {
		return nil, err
	}
	response, err := toStaff(staff)
	if err != nil {
		return nil, err
	}
	return contract.GetStaff200JSONResponse(response), nil
}
