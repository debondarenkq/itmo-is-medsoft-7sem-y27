package api

import (
	"context"
	"strings"

	internalapi "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/api/internalapi"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/internal/models"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/internal/pkg/modules"
	"github.com/google/uuid"
)

func (implementation *Implementation) Authorize(ctx context.Context, request internalapi.AuthorizeRequestObject) (internalapi.AuthorizeResponseObject, error) {
	token := ""
	if request.Params.Authorization != nil {
		token = strings.TrimPrefix(*request.Params.Authorization, "Bearer ")
	}
	claims, err := implementation.Module.VerifyAccess(ctx, token)
	input := request.Body
	if err != nil {
		return nil, implementation.auditFailure(ctx, err, models.AuditEventKindAccessDenied, "Неавторизованный запрос "+string(input.Method)+" "+input.Path)
	}
	if !modules.Allowed(claims, string(input.Method), input.Path) {
		if err := implementation.Module.Audit(ctx, "account:"+claims.Subject, "Отказ в доступе "+string(input.Method)+" "+input.Path, models.AuditEventKindAccessDenied); err != nil {
			return nil, err
		}
		return nil, models.ErrForbidden
	}
	subject, err := uuid.Parse(claims.Subject)
	if err != nil {
		return nil, err
	}
	jti, err := uuid.Parse(claims.ID)
	if err != nil {
		return nil, err
	}
	var staffID *uuid.UUID
	if claims.StaffID != "" {
		staffID, err = toStaffID(&claims.StaffID)
		if err != nil {
			return nil, err
		}
	}
	return internalapi.Authorize200JSONResponse{Sub: subject, StaffID: staffID, Role: internalapi.ClaimsRole(claims.Role), Iss: claims.Issuer, Aud: claims.Audience, Exp: claims.Expires, Iat: claims.Issued, Jti: jti}, nil
}
