package api

import (
	"context"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/api"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/pkg/platform"
)

func (i *Implementation) DeletePatient(ctx context.Context, request contract.DeletePatientRequestObject) (contract.DeletePatientResponseObject, error) {
	if err := platform.RequireID(request.ID); err != nil {
		return nil, err
	}
	if err := i.Module.DeletePatient(ctx, request.ID.String()); err != nil {
		return nil, err
	}
	return contract.DeletePatient204Response{}, nil
}
