package api

import (
	"context"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/api/security"
)

func (i *Implementation) UpdateServiceAccount(ctx context.Context, r contract.UpdateServiceAccountRequestObject) (contract.UpdateServiceAccountResponseObject, error) {
	c, err := requestClaims(ctx)
	if err != nil {
		return nil, err
	}
	a, _, err := i.Module.SaveServiceAccount(ctx, r.ID.String(), r.Body.Name, r.Body.IsActive, false, "account:"+c.Subject)
	if err != nil {
		return nil, err
	}
	return contract.UpdateServiceAccount200JSONResponse(serviceAccount(a)), nil
}
