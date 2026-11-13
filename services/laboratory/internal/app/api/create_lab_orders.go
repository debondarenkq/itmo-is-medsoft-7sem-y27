package api

import (
	"context"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/laboratory/api"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/laboratory/internal/models"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/laboratory/internal/pkg/adapters/security"
	"github.com/google/uuid"
)

func (i *Implementation) CreateLabOrders(ctx context.Context, r contract.CreateLabOrdersRequestObject) (contract.CreateLabOrdersResponseObject, error) {
	id, err := uuid.Parse(security.RequestStaffID(ctx))
	if err != nil {
		return nil, models.ErrInvalidInput
	}
	bundle, err := i.Module.CreateOrders(ctx, id, *r.Body)
	return contract.CreateLabOrders201JSONResponse(bundle), err
}
