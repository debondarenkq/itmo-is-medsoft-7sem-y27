package api

import (
	"context"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/api"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/internal/pkg/platform"
)

func (i *Implementation) GetStaff(ctx context.Context, request contract.GetStaffRequestObject) (contract.GetStaffResponseObject, error) {
	if err := platform.RequireID(request.ID); err != nil {
		return nil, err
	}
	value, err := i.Module.Get(ctx, request.ID.String())
	if err != nil {
		return nil, err
	}
	return contract.GetStaff200JSONResponse(toStaff(value)), nil
}
