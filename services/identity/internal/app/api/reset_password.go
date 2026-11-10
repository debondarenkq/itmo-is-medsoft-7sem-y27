package api

import (
	"context"
	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/api/security"
)

func (implementation *Implementation) ResetPassword(ctx context.Context, request contract.ResetPasswordRequestObject) (contract.ResetPasswordResponseObject, error) {
	claims, err := requestClaims(ctx)
	if err != nil {
		return nil, err
	}
	account, password, err := implementation.Module.ResetPassword(ctx, request.ID.String(), "account:"+claims.Subject)
	if err != nil {
		return nil, err
	}
	response, err := toAccount(account)
	if err != nil {
		return nil, err
	}
	return contract.ResetPassword200JSONResponse{Account: response, Password: password}, nil
}
