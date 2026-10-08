package api

import (
	"context"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/catalog/api"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/catalog/internal/pkg/platform"
)

func (i *Implementation) DeleteDiagnosis(ctx context.Context, request contract.DeleteDiagnosisRequestObject) (contract.DeleteDiagnosisResponseObject, error) {
	if err := platform.RequireID(request.ID); err != nil {
		return nil, err
	}
	if err := i.Module.Delete(ctx, request.ID.String()); err != nil {
		return nil, err
	}
	return contract.DeleteDiagnosis204Response{}, nil
}
