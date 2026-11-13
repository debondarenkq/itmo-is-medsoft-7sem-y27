package api

import (
	"context"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/lis/api"
)

func (i *Implementation) GetLISConnection(ctx context.Context, _ contract.GetLISConnectionRequestObject) (contract.GetLISConnectionResponseObject, error) {
	v, e := i.Module.Status(ctx)
	return contract.GetLISConnection200JSONResponse(v), e
}
