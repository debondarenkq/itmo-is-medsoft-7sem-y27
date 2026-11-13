package api

import (
	"context"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/lis/api"
)

func (i *Implementation) ConfigureLISConnection(ctx context.Context, r contract.ConfigureLISConnectionRequestObject) (contract.ConfigureLISConnectionResponseObject, error) {
	v, e := i.Module.Configure(ctx, *r.Body)
	return contract.ConfigureLISConnection200JSONResponse(v), e
}
