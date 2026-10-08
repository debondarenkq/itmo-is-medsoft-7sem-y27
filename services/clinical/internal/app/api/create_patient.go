package api

import (
	"context"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/api"
)

func (i *Implementation) CreatePatient(ctx context.Context, request contract.CreatePatientRequestObject) (contract.CreatePatientResponseObject, error) {
	value, err := i.Module.CreatePatient(ctx, fromPatientInput(*request.Body))
	if err != nil {
		return nil, err
	}
	return contract.CreatePatient201JSONResponse(toPatient(value)), nil
}
