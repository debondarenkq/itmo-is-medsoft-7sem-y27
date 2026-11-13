package api

import (
	"context"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/lis/api"
)

func (i *Implementation) ListLISHistory(ctx context.Context, r contract.ListLISHistoryRequestObject) (contract.ListLISHistoryResponseObject, error) {
	search := ""
	if r.Params.Search != nil {
		search = *r.Params.Search
	}
	v, e := i.Module.Store.History(ctx, search)
	return contract.ListLISHistory200JSONResponse(v), e
}
