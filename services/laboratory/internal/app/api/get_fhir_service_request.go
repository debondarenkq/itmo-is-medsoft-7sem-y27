package api

import (
	"context"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/laboratory/api"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/laboratory/internal/models"
)

func (i *Implementation) GetFHIRServiceRequest(ctx context.Context, r contract.GetFHIRServiceRequestRequestObject) (contract.GetFHIRServiceRequestResponseObject, error) {
	items, err := i.Module.Store.Requests(ctx, "", r.ID.String())
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, models.ErrNotFound
	}
	return contract.GetFHIRServiceRequest200ApplicationFhirPlusJSONResponse(items[0]), nil
}
