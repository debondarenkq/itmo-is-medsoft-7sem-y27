package api

import (
	"context"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/laboratory/api"
)

func (i *Implementation) UpdateLabTest(ctx context.Context, r contract.UpdateLabTestRequestObject) (contract.UpdateLabTestResponseObject, error) {
	t, err := i.Module.SaveTest(ctx, r.ID, *r.Body, false)
	return contract.UpdateLabTest200JSONResponse(t), err
}
