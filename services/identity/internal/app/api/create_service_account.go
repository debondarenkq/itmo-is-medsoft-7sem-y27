package api

import (
	"context"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/api/security"
)

func (i *Implementation) CreateServiceAccount(ctx context.Context, r contract.CreateServiceAccountRequestObject) (contract.CreateServiceAccountResponseObject, error) {
	c, err := requestClaims(ctx)
	if err != nil {
		return nil, err
	}
	a, key, err := i.Module.SaveServiceAccount(ctx, "", r.Body.Name, true, true, "account:"+c.Subject)
	if err != nil {
		return nil, err
	}
	return contract.CreateServiceAccount201JSONResponse{Account: serviceAccount(a), APIKey: key}, nil
}
