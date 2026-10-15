package api

import (
	"context"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/api"
)

func (implementation *Implementation) CreateStaff(ctx context.Context, request contract.CreateStaffRequestObject) (contract.CreateStaffResponseObject, error) {
	staff, err := implementation.Module.Create(ctx, fromStaffInput(*request.Body))
	if err != nil {
		return nil, err
	}
	response, err := toStaff(staff)
	if err != nil {
		return nil, err
	}
	return contract.CreateStaff201JSONResponse(response), nil
}
