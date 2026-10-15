package api

import (
	"context"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/api"
)

func (implementation *Implementation) CreatePatient(ctx context.Context, request contract.CreatePatientRequestObject) (contract.CreatePatientResponseObject, error) {
	patient, err := implementation.Module.CreatePatient(ctx, fromPatientInput(*request.Body))
	if err != nil {
		return nil, err
	}
	response, err := toPatient(patient)
	if err != nil {
		return nil, err
	}
	return contract.CreatePatient201JSONResponse(response), nil
}
