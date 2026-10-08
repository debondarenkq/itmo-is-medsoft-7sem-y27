package api

import (
	"context"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/catalog/api"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/catalog/internal/pkg/platform"
)

func (i *Implementation) GetDiagnosis(ctx context.Context, request contract.GetDiagnosisRequestObject) (contract.GetDiagnosisResponseObject, error) {
	if err := platform.RequireID(request.ID); err != nil {
		return nil, err
	}
	value, err := i.Module.Get(ctx, request.ID.String())
	if err != nil {
		return nil, err
	}
	return contract.GetDiagnosis200JSONResponse(toDiagnosis(value)), nil
}
