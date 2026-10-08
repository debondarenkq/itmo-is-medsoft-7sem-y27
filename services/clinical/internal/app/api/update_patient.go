package api

import (
	"context"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/api"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/pkg/platform"
)

func (i *Implementation) UpdatePatient(ctx context.Context, request contract.UpdatePatientRequestObject) (contract.UpdatePatientResponseObject, error) {
	if err := platform.RequireID(request.ID); err != nil {
		return nil, err
	}
	value, err := i.Module.UpdatePatient(ctx, request.ID.String(), fromPatientInput(*request.Body))
	if err != nil {
		return nil, err
	}
	return contract.UpdatePatient200JSONResponse(toPatient(value)), nil
}
