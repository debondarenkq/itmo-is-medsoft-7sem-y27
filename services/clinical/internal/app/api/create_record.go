package api

import (
	"context"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/api"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/pkg/platform"
)

func (i *Implementation) CreateRecord(ctx context.Context, request contract.CreateRecordRequestObject) (contract.CreateRecordResponseObject, error) {
	if err := platform.RequireID(request.ID); err != nil {
		return nil, err
	}
	value, err := i.Module.CreateRecord(ctx, request.ID.String(), request.Body.StaffID.String())
	if err != nil {
		return nil, err
	}
	return contract.CreateRecord201JSONResponse(toRecord(value)), nil
}
