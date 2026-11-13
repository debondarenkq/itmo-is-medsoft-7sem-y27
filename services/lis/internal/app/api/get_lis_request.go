package api

import (
	"context"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/lis/api"
)

func (i *Implementation) GetLISRequest(ctx context.Context, r contract.GetLISRequestRequestObject) (contract.GetLISRequestResponseObject, error) {
	v, e := i.Module.Details(ctx, r.ID)
	return contract.GetLISRequest200JSONResponse(v), e
}
