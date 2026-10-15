package api

import (
	"context"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/catalog/api"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/catalog/internal/validation"
)

func (implementation *Implementation) UpdateDiagnosis(ctx context.Context, request contract.UpdateDiagnosisRequestObject) (contract.UpdateDiagnosisResponseObject, error) {
	if err := validation.RequireID(request.ID); err != nil {
		return nil, err
	}
	diagnosis, err := implementation.Module.Update(ctx, request.ID.String(), fromDiagnosisInput(*request.Body))
	if err != nil {
		return nil, err
	}
	response, err := toDiagnosis(diagnosis)
	if err != nil {
		return nil, err
	}
	return contract.UpdateDiagnosis200JSONResponse(response), nil
}
