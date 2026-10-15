package api

import (
	"context"
	"fmt"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/api"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/models"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/validation"
)

func (implementation *Implementation) GetRecordState(ctx context.Context, request contract.GetRecordStateRequestObject) (contract.GetRecordStateResponseObject, error) {
	if err := validation.RequireID(request.ID); err != nil {
		return nil, err
	}
	if (request.Params.At == nil) == (request.Params.Version == nil) {
		return nil, fmt.Errorf("%w: Укажите ровно один параметр: at или version", models.ErrInvalidInput)
	}
	query := models.StateQuery{At: request.Params.At, Version: request.Params.Version}
	record, err := implementation.Module.StateAt(ctx, request.ID.String(), query)
	if err != nil {
		return nil, err
	}
	response, err := toRecord(record)
	if err != nil {
		return nil, err
	}
	return contract.GetRecordState200JSONResponse(response), nil
}
