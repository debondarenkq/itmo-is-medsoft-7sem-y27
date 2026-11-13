package api

import (
	"context"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/laboratory/api"
	"github.com/google/uuid"
)

func (i *Implementation) CreateLabTest(ctx context.Context, r contract.CreateLabTestRequestObject) (contract.CreateLabTestResponseObject, error) {
	t, err := i.Module.SaveTest(ctx, uuid.Nil, *r.Body, true)
	return contract.CreateLabTest201JSONResponse(t), err
}
