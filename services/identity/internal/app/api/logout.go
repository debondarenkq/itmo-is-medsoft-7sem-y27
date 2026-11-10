package api

import (
	"context"
	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/api/security"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/internal/pkg/modules"
)

func (implementation *Implementation) Logout(ctx context.Context, request contract.LogoutRequestObject) (contract.LogoutResponseObject, error) {
	claims, err := requestClaims(ctx)
	if err != nil {
		return nil, err
	}
	if err = implementation.Module.Store.Logout(ctx, claims.Subject, modules.TokenHash(request.Body.RefreshToken)); err != nil {
		return nil, err
	}
	return contract.Logout200JSONResponse{Configured: true}, nil
}
