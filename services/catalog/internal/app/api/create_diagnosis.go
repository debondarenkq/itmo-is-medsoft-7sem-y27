package api

import (
	"context"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/catalog/api"
)

func (implementation *Implementation) CreateDiagnosis(ctx context.Context, request contract.CreateDiagnosisRequestObject) (contract.CreateDiagnosisResponseObject, error) {
	diagnosis, err := implementation.Module.Create(ctx, fromDiagnosisInput(*request.Body))
	if err != nil {
		return nil, err
	}
	response, err := toDiagnosis(diagnosis)
	if err != nil {
		return nil, err
	}
	return contract.CreateDiagnosis201JSONResponse(response), nil
}
