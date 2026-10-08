package api

import (
	"context"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/catalog/api"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/catalog/internal/pkg/platform"
)

func (i *Implementation) ListDiagnoses(ctx context.Context, request contract.ListDiagnosesRequestObject) (contract.ListDiagnosesResponseObject, error) {
	page := platform.PageFromParams(request.Params.Limit, request.Params.Offset, request.Params.Q, request.Params.IncludeDeleted)
	values, total, err := i.Module.List(ctx, page)
	if err != nil {
		return nil, err
	}
	items := make([]contract.Diagnosis, len(values))
	for index, value := range values {
		items[index] = toDiagnosis(value)
	}
	return contract.ListDiagnoses200JSONResponse{Items: items, Total: total, Limit: page.Limit, Offset: page.Offset}, nil
}
