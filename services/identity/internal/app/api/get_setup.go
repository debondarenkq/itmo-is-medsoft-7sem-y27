package api

import (
	"context"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/api/security"
)

func (implementation *Implementation) GetSetup(ctx context.Context, request contract.GetSetupRequestObject) (contract.GetSetupResponseObject, error) {
	setup, err := implementation.Module.SetupStatus(ctx)
	if err != nil {
		return nil, err
	}
	return contract.GetSetup200JSONResponse{Configured: setup.Configured, Secret: setup.Secret, ProvisioningURI: setup.ProvisioningURI}, nil
}
