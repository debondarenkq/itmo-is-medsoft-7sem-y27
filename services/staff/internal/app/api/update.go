package api

import (
	"context"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/api"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/internal/pkg/platform"
)

func (i *Implementation) UpdateStaff(ctx context.Context, request contract.UpdateStaffRequestObject) (contract.UpdateStaffResponseObject, error) {
	if err := platform.RequireID(request.ID); err != nil {
		return nil, err
	}
	value, err := i.Module.Update(ctx, request.ID.String(), fromStaffInput(*request.Body))
	if err != nil {
		return nil, err
	}
	return contract.UpdateStaff200JSONResponse(toStaff(value)), nil
}
