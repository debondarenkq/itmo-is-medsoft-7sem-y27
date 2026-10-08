package api

import (
	"context"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/api"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/pkg/platform"
)

func (i *Implementation) GetRecord(ctx context.Context, request contract.GetRecordRequestObject) (contract.GetRecordResponseObject, error) {
	if err := platform.RequireID(request.ID); err != nil {
		return nil, err
	}
	value, err := i.Module.GetRecord(ctx, request.ID.String())
	if err != nil {
		return nil, err
	}
	return contract.GetRecord200JSONResponse(toRecord(value)), nil
}
