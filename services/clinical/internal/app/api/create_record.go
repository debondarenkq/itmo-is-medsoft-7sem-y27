package api

import (
	"context"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/api"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/validation"
)

func (implementation *Implementation) CreateRecord(ctx context.Context, request contract.CreateRecordRequestObject) (contract.CreateRecordResponseObject, error) {
	if err := validation.RequireID(request.ID); err != nil {
		return nil, err
	}
	record, err := implementation.Module.CreateRecord(ctx, request.ID.String(), request.Body.StaffID.String())
	if err != nil {
		return nil, err
	}
	response, err := toRecord(record)
	if err != nil {
		return nil, err
	}
	return contract.CreateRecord201JSONResponse(response), nil
}
