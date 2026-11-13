package api

import (
	"context"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/api/security"
)

func (i *Implementation) LoginServiceAccount(ctx context.Context, r contract.LoginServiceAccountRequestObject) (contract.LoginServiceAccountResponseObject, error) {
	session, err := i.Module.LoginServiceAccount(ctx, r.Body.APIKey)
	if err != nil {
		return nil, err
	}
	result, err := toSession(session)
	return contract.LoginServiceAccount200JSONResponse(result), err
}
