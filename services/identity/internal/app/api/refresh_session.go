package api

import (
	"context"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/api/security"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/internal/models"
)

func (implementation *Implementation) RefreshSession(ctx context.Context, request contract.RefreshSessionRequestObject) (contract.RefreshSessionResponseObject, error) {
	session, err := implementation.Module.Refresh(ctx, request.Body.RefreshToken)
	if err != nil {
		return nil, implementation.auditFailure(ctx, err, models.AuditEventKindAccessDenied, "Отклонено обновление пары токенов")
	}
	response, err := toSession(session)
	if err != nil {
		return nil, err
	}
	return contract.RefreshSession200JSONResponse(response), nil
}
