package api

import (
	"context"

	internalapi "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/api/internalapi"
)

func (i *Implementation) GetStaffIdentity(ctx context.Context, request internalapi.GetStaffIdentityRequestObject) (internalapi.GetStaffIdentityResponseObject, error) {
	s, err := i.Module.Get(ctx, request.ID.String())
	if err != nil {
		return nil, err
	}
	response, err := toStaffProfile(s)
	if err != nil {
		return nil, err
	}
	return internalapi.GetStaffIdentity200JSONResponse{ProfileJSONResponse: response}, nil
}
