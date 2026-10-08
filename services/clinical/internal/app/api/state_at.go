package api

import (
	"context"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/api"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/models"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/pkg/platform"
)

func (i *Implementation) GetRecordState(ctx context.Context, request contract.GetRecordStateRequestObject) (contract.GetRecordStateResponseObject, error) {
	if err := platform.RequireID(request.ID); err != nil {
		return nil, err
	}
	if (request.Params.At == nil) == (request.Params.Version == nil) {
		return nil, platform.Bad("STATE_QUERY", "Укажите ровно один параметр: at или version")
	}
	query := models.StateQuery{At: request.Params.At, Version: request.Params.Version}
	value, err := i.Module.StateAt(ctx, request.ID.String(), query)
	if err != nil {
		return nil, err
	}
	return contract.GetRecordState200JSONResponse(toRecord(value)), nil
}
