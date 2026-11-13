package api

import (
	"context"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/laboratory/api"
)

func (i *Implementation) CreateFHIRObservation(ctx context.Context, r contract.CreateFHIRObservationRequestObject) (contract.CreateFHIRObservationResponseObject, error) {
	result, err := i.Module.SaveObservation(ctx, *r.Body)
	return contract.CreateFHIRObservation201ApplicationFhirPlusJSONResponse(result), err
}
