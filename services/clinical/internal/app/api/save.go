package api

import (
	"context"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/api"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/pkg/platform"
)

func (i *Implementation) SaveRecord(ctx context.Context, request contract.SaveRecordRequestObject) (contract.SaveRecordResponseObject, error) {
	if err := platform.RequireID(request.ID); err != nil {
		return nil, err
	}
	value, err := i.Module.Save(ctx, request.ID.String(), fromChangesInput(*request.Body))
	if err != nil {
		return nil, err
	}
	return contract.SaveRecord200JSONResponse(toRecord(value)), nil
}
