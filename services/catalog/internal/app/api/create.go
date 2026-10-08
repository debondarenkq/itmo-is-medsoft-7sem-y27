package api

import (
	"context"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/catalog/api"
)

func (i *Implementation) CreateDiagnosis(ctx context.Context, request contract.CreateDiagnosisRequestObject) (contract.CreateDiagnosisResponseObject, error) {
	value, err := i.Module.Create(ctx, fromDiagnosisInput(*request.Body))
	if err != nil {
		return nil, err
	}
	return contract.CreateDiagnosis201JSONResponse(toDiagnosis(value)), nil
}
