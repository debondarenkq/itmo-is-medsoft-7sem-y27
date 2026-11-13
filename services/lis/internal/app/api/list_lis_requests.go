package api

import (
	"context"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/lis/api"
)

func (i *Implementation) ListLISRequests(ctx context.Context, _ contract.ListLISRequestsRequestObject) (contract.ListLISRequestsResponseObject, error) {
	v, e := i.Module.Requests(ctx)
	return contract.ListLISRequests200JSONResponse(v), e
}
