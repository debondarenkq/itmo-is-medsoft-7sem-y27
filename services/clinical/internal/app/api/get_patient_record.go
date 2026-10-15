package api

import (
	"context"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/api"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/validation"
)

func (implementation *Implementation) GetPatientRecord(ctx context.Context, request contract.GetPatientRecordRequestObject) (contract.GetPatientRecordResponseObject, error) {
	if err := validation.RequireID(request.ID); err != nil {
		return nil, err
	}
	record, err := implementation.Module.GetPatientRecord(ctx, request.ID.String())
	if err != nil {
		return nil, err
	}
	response, err := toRecord(record)
	if err != nil {
		return nil, err
	}
	return contract.GetPatientRecord200JSONResponse(response), nil
}
