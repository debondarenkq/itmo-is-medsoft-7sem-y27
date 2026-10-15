package api

import (
	"context"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/api"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/internal/validation"
)

func (implementation *Implementation) DeleteStaff(ctx context.Context, request contract.DeleteStaffRequestObject) (contract.DeleteStaffResponseObject, error) {
	if err := validation.RequireID(request.ID); err != nil {
		return nil, err
	}
	if err := implementation.Module.Delete(ctx, request.ID.String()); err != nil {
		return nil, err
	}
	return contract.DeleteStaff204Response{}, nil
}
