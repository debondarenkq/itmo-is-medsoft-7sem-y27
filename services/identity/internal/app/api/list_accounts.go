package api

import (
	"context"
	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/api/security"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/internal/models"
	"strings"
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
