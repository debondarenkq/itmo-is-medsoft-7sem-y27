package api

import (
	"context"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/api"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/pkg/platform"
)

func (i *Implementation) ListPatients(ctx context.Context, request contract.ListPatientsRequestObject) (contract.ListPatientsResponseObject, error) {
	page := platform.PageFromParams(request.Params.Limit, request.Params.Offset, request.Params.Q, nil)
	values, total, err := i.Module.ListPatients(ctx, page)
	if err != nil {
		return nil, err
	}
	items := make([]contract.Patient, len(values))
	for index, value := range values {
		items[index] = toPatient(value)
	}
	return contract.ListPatients200JSONResponse{Items: items, Total: total, Limit: page.Limit, Offset: page.Offset}, nil
}
