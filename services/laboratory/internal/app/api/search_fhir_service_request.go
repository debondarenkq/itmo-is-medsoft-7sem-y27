package api

import (
	"context"
	"strings"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/laboratory/api"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/laboratory/internal/models"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/laboratory/internal/pkg/modules"
	"github.com/google/uuid"
)

func (i *Implementation) SearchFHIRServiceRequest(ctx context.Context, r contract.SearchFHIRServiceRequestRequestObject) (contract.SearchFHIRServiceRequestResponseObject, error) {
	if err := validateSearch(r.Params.Patient, r.Params.UnderscoreID); err != nil {
		return nil, err
	}
	items, err := i.Module.Store.Requests(ctx, strings.TrimPrefix(value(r.Params.Patient), "Patient/"), value(r.Params.UnderscoreID))
	return contract.SearchFHIRServiceRequest200ApplicationFhirPlusJSONResponse(modules.Bundle(items)), err
}

func value(p *string) string {
	if p != nil {
		return *p
	}
	return ""
}

func validateSearch(patient, id *string) error {
	for _, v := range []*string{patient, id} {
		if v != nil {
			if _, err := uuid.Parse(strings.TrimPrefix(*v, "Patient/")); err != nil {
				return models.ErrInvalidInput
			}
		}
	}
	return nil
}
