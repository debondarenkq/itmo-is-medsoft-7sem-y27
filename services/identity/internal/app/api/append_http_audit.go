package api

import (
	"context"

	auditcontract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/api/audit"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/internal/models"
)

func (implementation *Implementation) AppendHTTPAudit(ctx context.Context, request auditcontract.AppendHTTPAuditRequestObject) (auditcontract.AppendHTTPAuditResponseObject, error) {
	input := request.Body
	event := models.HTTPAudit{Actor: input.Actor, Method: models.HTTPAuditMethod(input.Method), Path: input.Path, Status: input.Status, EventType: models.AuditEventKind(input.EventType), Title: input.Title}
	if err := implementation.Module.AuditHTTP(ctx, event); err != nil {
		return nil, err
	}
	return auditcontract.AppendHTTPAudit200JSONResponse{Ok: true}, nil
}
