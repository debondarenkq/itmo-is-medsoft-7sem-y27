package api

import (
	"context"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/laboratory/api"
)

func (i *Implementation) GetFHIRPractitioner(ctx context.Context, r contract.GetFHIRPractitionerRequestObject) (contract.GetFHIRPractitionerResponseObject, error) {
	p, err := i.Module.Directories.Practitioner(ctx, r.ID)
	return contract.GetFHIRPractitioner200ApplicationFhirPlusJSONResponse(p), err
}
