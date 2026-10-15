package api

import (
	"context"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/api"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/validation"
)

func (implementation *Implementation) GetPatient(ctx context.Context, request contract.GetPatientRequestObject) (contract.GetPatientResponseObject, error) {
	if err := validation.RequireID(request.ID); err != nil {
		return nil, err
	}
	patient, err := implementation.Module.GetPatient(ctx, request.ID.String())
	if err != nil {
		return nil, err
	}
	response, err := toPatient(patient)
	if err != nil {
		return nil, err
	}
	return contract.GetPatient200JSONResponse(response), nil
}
