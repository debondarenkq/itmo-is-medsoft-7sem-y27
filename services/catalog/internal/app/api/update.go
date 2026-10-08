package api

import (
	"context"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/catalog/api"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/catalog/internal/pkg/platform"
)

func (i *Implementation) UpdateDiagnosis(ctx context.Context, request contract.UpdateDiagnosisRequestObject) (contract.UpdateDiagnosisResponseObject, error) {
	if err := platform.RequireID(request.ID); err != nil {
		return nil, err
	}
	value, err := i.Module.Update(ctx, request.ID.String(), fromDiagnosisInput(*request.Body))
	if err != nil {
		return nil, err
	}
	return contract.UpdateDiagnosis200JSONResponse(toDiagnosis(value)), nil
}
