package api

import (
	"context"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/laboratory/api"
)

func (i *Implementation) GetFHIRPatient(ctx context.Context, r contract.GetFHIRPatientRequestObject) (contract.GetFHIRPatientResponseObject, error) {
	p, err := i.Module.Directories.Patient(ctx, r.ID, false)
	return contract.GetFHIRPatient200ApplicationFhirPlusJSONResponse(p), err
}
