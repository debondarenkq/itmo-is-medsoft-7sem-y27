package api

import (
	"context"
	"strings"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/api/security"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/internal/models"
	"github.com/google/uuid"
)

func (implementation *Implementation) ListAccounts(ctx context.Context, request contract.ListAccountsRequestObject) (contract.ListAccountsResponseObject, error) {
	params := request.Params
	limit, offset := pagination(params.Limit, params.Offset)
	filter := models.AccountFilter{Limit: limit, Offset: offset, IsActive: params.IsActive}
	if params.Search != nil {
		filter.Search = strings.TrimSpace(*params.Search)
	}
	items, total, err := implementation.Module.SearchAccounts(ctx, filter)
	if err != nil {
		return nil, err
	}
	response := make([]contract.Account, 0, len(items))
	for _, item := range items {
		account, err := toAccount(item)
		if err != nil {
			return nil, err
		}
		response = append(response, account)
	}
	return contract.ListAccounts200JSONResponse{Items: response, Total: total}, nil
}

func pagination(limit, offset *int) (int, int) {
	l, o := 50, 0
	if limit != nil {
		l = *limit
	}
	if offset != nil {
		o = *offset
	}
	if l < 1 {
		l = 1
	}
	if l > 200 {
		l = 200
	}
	if o < 0 {
		o = 0
	}
	return l, o
}

func toAccount(a models.Account) (contract.Account, error) {
	id, err := uuid.Parse(a.ID)
	if err != nil {
		return contract.Account{}, err
	}
	staffID, err := toStaffID(a.StaffID)
	if err != nil {
		return contract.Account{}, err
	}
	return contract.Account{ID: id, Username: a.Username, CreatedAt: a.CreatedAt, UpdatedAt: a.UpdatedAt, IsActive: a.IsActive, StaffID: staffID, Role: contract.AccountRole(a.Role), StaffName: a.StaffName}, nil
}
