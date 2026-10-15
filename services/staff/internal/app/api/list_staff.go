package api

import (
	"context"
	"strings"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/api"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/internal/models"
)

func (implementation *Implementation) ListStaff(ctx context.Context, request contract.ListStaffRequestObject) (contract.ListStaffResponseObject, error) {
	page := pageFromParams(request.Params.Limit, request.Params.Offset, request.Params.Search, request.Params.IncludeDeleted)
	staffMembers, total, err := implementation.Module.List(ctx, page)
	if err != nil {
		return nil, err
	}
	items := make([]contract.Staff, len(staffMembers))
	for index, staff := range staffMembers {
		items[index], err = toStaff(staff)
		if err != nil {
			return nil, err
		}
	}
	return contract.ListStaff200JSONResponse{Items: items, Total: total, Limit: page.Limit, Offset: page.Offset}, nil
}

func pageFromParams(limit, offset *int, search *string, includeDeleted *bool) models.Page {
	page := models.Page{Limit: 50}
	if limit != nil {
		page.Limit = *limit
	}
	if offset != nil {
		page.Offset = *offset
	}
	if search != nil {
		page.Search = strings.TrimSpace(*search)
	}
	if includeDeleted != nil {
		page.IncludeDeleted = *includeDeleted
	}
	return page
}
