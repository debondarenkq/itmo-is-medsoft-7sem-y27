package api

import (
	"context"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/api/security"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/internal/models"
)

func (i *Implementation) RotateServiceAccountKey(ctx context.Context, r contract.RotateServiceAccountKeyRequestObject) (contract.RotateServiceAccountKeyResponseObject, error) {
	c, err := requestClaims(ctx)
	if err != nil {
		return nil, err
	}
	a, err := i.Module.Store.AccountByID(ctx, r.ID.String())
	if err != nil {
		return nil, err
	}
	if !a.IsService {
		return nil, models.ErrNotFound
	}
	a, key, err := i.Module.SaveServiceAccount(ctx, a.ID, a.Username, a.IsActive, true, "account:"+c.Subject)
	if err != nil {
		return nil, err
	}
	return contract.RotateServiceAccountKey200JSONResponse{Account: serviceAccount(a), APIKey: key}, nil
}
