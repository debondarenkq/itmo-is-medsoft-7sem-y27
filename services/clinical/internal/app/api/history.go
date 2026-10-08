package api

import (
	"context"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/api"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/pkg/platform"
)

func (i *Implementation) GetRecordHistory(ctx context.Context, request contract.GetRecordHistoryRequestObject) (contract.GetRecordHistoryResponseObject, error) {
	if err := platform.RequireID(request.ID); err != nil {
		return nil, err
	}
	page := platform.PageFromParams(request.Params.Limit, request.Params.Offset, nil, nil)
	values, total, err := i.Module.History(ctx, request.ID.String(), page)
	if err != nil {
		return nil, err
	}
	items := make([]contract.Event, len(values))
	for index, value := range values {
		items[index] = toEvent(value)
	}
	return contract.GetRecordHistory200JSONResponse{Items: items, Total: total, Limit: page.Limit, Offset: page.Offset}, nil
}
