package api

import (
	"context"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/api"
)

func (implementation *Implementation) GetPatientIdentity(ctx context.Context, request contract.GetPatientIdentityRequestObject) (contract.GetPatientIdentityResponseObject, error) {
	patient, err := implementation.Module.GetPatient(ctx, request.ID.String())
	if err != nil {
		return nil, err
	}
	response, err := toPatient(patient)
	return contract.GetPatientIdentity200JSONResponse(response), err
}
