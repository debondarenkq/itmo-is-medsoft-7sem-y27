package api

import (
	"context"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/api"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/validation"
)

func (implementation *Implementation) UpdatePatient(ctx context.Context, request contract.UpdatePatientRequestObject) (contract.UpdatePatientResponseObject, error) {
	if err := validation.RequireID(request.ID); err != nil {
		return nil, err
	}
	patient, err := implementation.Module.UpdatePatient(ctx, request.ID.String(), fromPatientInput(*request.Body))
	if err != nil {
		return nil, err
	}
	response, err := toPatient(patient)
	if err != nil {
		return nil, err
	}
	return contract.UpdatePatient200JSONResponse(response), nil
}
