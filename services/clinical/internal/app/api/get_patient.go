package api

import (
	"context"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/api"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/pkg/platform"
)

func (i *Implementation) GetPatient(ctx context.Context, request contract.GetPatientRequestObject) (contract.GetPatientResponseObject, error) {
	if err := platform.RequireID(request.ID); err != nil {
		return nil, err
	}
	value, err := i.Module.GetPatient(ctx, request.ID.String())
	if err != nil {
		return nil, err
	}
	return contract.GetPatient200JSONResponse(toPatient(value)), nil
}
