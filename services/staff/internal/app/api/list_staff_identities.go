package api

import (
	"context"

	internalapi "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/api/internalapi"
)

func (i *Implementation) ListStaffIdentities(ctx context.Context, _ internalapi.ListStaffIdentitiesRequestObject) (internalapi.ListStaffIdentitiesResponseObject, error) {
	items, err := i.Module.ListIdentityProfiles(ctx)
	if err != nil {
		return nil, err
	}
	response := make(internalapi.ListStaffIdentities200JSONResponse, 0, len(items))
	for _, item := range items {
		profile, err := toStaffProfile(item)
		if err != nil {
			return nil, err
		}
		response = append(response, internalapi.StaffProfile(profile))
	}
	return response, nil
}
