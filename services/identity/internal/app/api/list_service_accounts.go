package api

import (
	"context"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/api/security"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/internal/models"
	"github.com/google/uuid"
)

func (i *Implementation) ListServiceAccounts(ctx context.Context, _ contract.ListServiceAccountsRequestObject) (contract.ListServiceAccountsResponseObject, error) {
	accounts, err := i.Module.ServiceAccounts(ctx)
	if err != nil {
		return nil, err
	}
	result := contract.ListServiceAccounts200JSONResponse{}
	for _, a := range accounts {
		result = append(result, serviceAccount(a))
	}
	return result, nil
}

func serviceAccount(a models.Account) contract.ServiceAccount {
	return contract.ServiceAccount{ID: uuid.MustParse(a.ID), Name: a.Username, Role: contract.Lis, IsActive: a.IsActive, CreatedAt: a.CreatedAt}
}
