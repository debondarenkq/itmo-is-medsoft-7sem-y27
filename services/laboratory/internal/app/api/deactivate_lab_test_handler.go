package api

import (
	"context"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/laboratory/api"
)

func (i *Implementation) DeactivateLabTest(ctx context.Context, r contract.DeactivateLabTestRequestObject) (contract.DeactivateLabTestResponseObject, error) {
	t, err := i.Module.Store.DeactivateTest(ctx, r.ID)
	return contract.DeactivateLabTest200JSONResponse(t), err
}
