package api

import (
	"context"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/api"
)

func (i *Implementation) CreateStaff(ctx context.Context, request contract.CreateStaffRequestObject) (contract.CreateStaffResponseObject, error) {
	value, err := i.Module.Create(ctx, fromStaffInput(*request.Body))
	if err != nil {
		return nil, err
	}
	return contract.CreateStaff201JSONResponse(toStaff(value)), nil
}
