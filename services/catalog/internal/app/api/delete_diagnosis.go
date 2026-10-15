package api

import (
	"context"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/catalog/api"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/catalog/internal/validation"
)

func (implementation *Implementation) DeleteDiagnosis(ctx context.Context, request contract.DeleteDiagnosisRequestObject) (contract.DeleteDiagnosisResponseObject, error) {
	if err := validation.RequireID(request.ID); err != nil {
		return nil, err
	}
	if err := implementation.Module.Delete(ctx, request.ID.String()); err != nil {
		return nil, err
	}
	return contract.DeleteDiagnosis204Response{}, nil
}
