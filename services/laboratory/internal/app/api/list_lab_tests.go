package api

import (
	"context"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/laboratory/api"
)

func (i *Implementation) ListLabTests(ctx context.Context, r contract.ListLabTestsRequestObject) (contract.ListLabTestsResponseObject, error) {
	tests, err := i.Module.Store.Tests(ctx, r.Params.IncludeInactive != nil && *r.Params.IncludeInactive)
	return contract.ListLabTests200JSONResponse(tests), err
}
