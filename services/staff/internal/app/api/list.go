package api

import (
	"context"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/api"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/internal/pkg/platform"
)

func (i *Implementation) ListStaff(ctx context.Context, request contract.ListStaffRequestObject) (contract.ListStaffResponseObject, error) {
	page := platform.PageFromParams(request.Params.Limit, request.Params.Offset, request.Params.Q, request.Params.IncludeDeleted)
	values, total, err := i.Module.List(ctx, page)
	if err != nil {
		return nil, err
	}
	items := make([]contract.Staff, len(values))
	for index, value := range values {
		items[index] = toStaff(value)
	}
	return contract.ListStaff200JSONResponse{Items: items, Total: total, Limit: page.Limit, Offset: page.Offset}, nil
}
