package api

import (
	"context"

	internalapi "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/api/internalapi"
)

func (i *Implementation) UpdateStaffIdentityRole(ctx context.Context, request internalapi.UpdateStaffIdentityRoleRequestObject) (internalapi.UpdateStaffIdentityRoleResponseObject, error) {
	s, err := i.Module.UpdateRole(ctx, request.ID.String(), string(request.Body.Role))
	if err != nil {
		return nil, err
	}
	response, err := toStaffProfile(s)
	if err != nil {
		return nil, err
	}
	return internalapi.UpdateStaffIdentityRole200JSONResponse{ProfileJSONResponse: response}, nil
}
