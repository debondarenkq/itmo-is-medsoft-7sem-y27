package api

import (
	"context"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/catalog/api"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/catalog/internal/validation"
)

func (implementation *Implementation) GetDiagnosis(ctx context.Context, request contract.GetDiagnosisRequestObject) (contract.GetDiagnosisResponseObject, error) {
	if err := validation.RequireID(request.ID); err != nil {
		return nil, err
	}
	diagnosis, err := implementation.Module.Get(ctx, request.ID.String())
	if err != nil {
		return nil, err
	}
	response, err := toDiagnosis(diagnosis)
	if err != nil {
		return nil, err
	}
	return contract.GetDiagnosis200JSONResponse(response), nil
}
