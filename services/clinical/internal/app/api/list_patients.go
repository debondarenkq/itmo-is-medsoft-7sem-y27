package api

import (
	"context"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/api"
)

func (implementation *Implementation) ListPatients(ctx context.Context, request contract.ListPatientsRequestObject) (contract.ListPatientsResponseObject, error) {
	page := pageFromParams(request.Params.Limit, request.Params.Offset, request.Params.Search)
	page.HasRecord = request.Params.HasRecord
	patients, total, err := implementation.Module.ListPatients(ctx, page)
	if err != nil {
		return nil, err
	}
	items := make([]contract.Patient, len(patients))
	for index, patient := range patients {
		items[index], err = toPatient(patient)
		if err != nil {
			return nil, err
		}
	}
	return contract.ListPatients200JSONResponse{Items: items, Total: total, Limit: page.Limit, Offset: page.Offset}, nil
}
