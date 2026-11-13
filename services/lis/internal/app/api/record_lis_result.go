package api

import (
	"context"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/lis/api"
)

func (i *Implementation) RecordLISResult(ctx context.Context, r contract.RecordLISResultRequestObject) (contract.RecordLISResultResponseObject, error) {
	v, e := i.Module.Record(ctx, *r.Body)
	return contract.RecordLISResult201JSONResponse(v), e
}
