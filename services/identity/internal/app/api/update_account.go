package api

import (
	"context"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/api/security"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/internal/models"
)

func (implementation *Implementation) UpdateAccount(ctx context.Context, request contract.UpdateAccountRequestObject) (contract.UpdateAccountResponseObject, error) {
	input := request.Body
	changes := models.AccountChanges{IsActive: &input.IsActive, Username: input.Username}
	if input.Role != nil {
		role := string(*input.Role)
		changes.Role = &role
	}
	claims, err := requestClaims(ctx)
	if err != nil {
		return nil, err
	}
	account, err := implementation.Module.UpdateAccount(ctx, request.ID.String(), changes, "account:"+claims.Subject)
	if err != nil {
		return nil, err
	}
	response, err := toAccount(account)
	if err != nil {
		return nil, err
	}
	return contract.UpdateAccount200JSONResponse(response), nil
}
