package api

import (
	"context"
	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/api/security"
)

func (implementation *Implementation) CreateAccount(ctx context.Context, request contract.CreateAccountRequestObject) (contract.CreateAccountResponseObject, error) {
	claims, err := requestClaims(ctx)
	if err != nil {
		return nil, err
	}
	account, password, err := implementation.Module.CreateAccount(ctx, request.Body.Username, request.Body.StaffID.String(), "account:"+claims.Subject)
	if err != nil {
		return nil, err
	}
	response, err := toAccount(account)
	if err != nil {
		return nil, err
	}
	return contract.CreateAccount201JSONResponse{Account: response, Password: password}, nil
}
