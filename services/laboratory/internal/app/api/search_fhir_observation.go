package api

import (
	"context"
	"strings"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/laboratory/api"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/laboratory/internal/pkg/modules"
)

func (i *Implementation) SearchFHIRObservation(ctx context.Context, r contract.SearchFHIRObservationRequestObject) (contract.SearchFHIRObservationResponseObject, error) {
	if err := validateSearch(r.Params.Patient, r.Params.UnderscoreID); err != nil {
		return nil, err
	}
	items, err := i.Module.Store.Observations(ctx, strings.TrimPrefix(value(r.Params.Patient), "Patient/"), value(r.Params.UnderscoreID))
	return contract.SearchFHIRObservation200ApplicationFhirPlusJSONResponse(modules.Bundle(items)), err
}
