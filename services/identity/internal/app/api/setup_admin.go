package api

import (
	"context"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/api/security"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/internal/models"
)

func (implementation *Implementation) SetupAdmin(ctx context.Context, request contract.SetupAdminRequestObject) (contract.SetupAdminResponseObject, error) {
	if implementation.limited(ctx) {
		return nil, implementation.auditFailure(ctx, models.ErrUnauthorized, models.AuditEventKindLoginFailed, "Превышен лимит настройки администратора")
	}
	input := request.Body
	session, err := implementation.Module.Setup(ctx, input.Password, input.Confirmation, input.Secret, input.Totp)
	if err != nil {
		return nil, implementation.auditFailure(ctx, err, models.AuditEventKindLoginFailed, "Неуспешная настройка администратора")
	}
	response, err := toSession(session)
	if err != nil {
		return nil, err
	}
	return contract.SetupAdmin200JSONResponse(response), nil
}
