package api

import (
	"context"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/api/security"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/internal/models"
	"github.com/google/uuid"
)

func (implementation *Implementation) ListAudit(ctx context.Context, request contract.ListAuditRequestObject) (contract.ListAuditResponseObject, error) {
	params := request.Params
	limit, offset := pagination(params.Limit, params.Offset)
	filter := models.AuditFilter{Limit: limit, Offset: offset, From: params.From, To: params.To}
	if params.EventType != nil {
		filter.EventType = models.AuditEventKind(*params.EventType)
	}
	items, total, err := implementation.Module.AuditEvents(ctx, filter)
	if err != nil {
		return nil, err
	}
	response := make([]contract.AuditEvent, 0, len(items))
	for _, item := range items {
		id, err := uuid.Parse(item.EventID)
		if err != nil {
			return nil, err
		}
		response = append(response, contract.AuditEvent{EventID: id, EventTime: item.EventTime, Actor: item.Actor, ActorName: item.ActorName, Message: item.Message, EventType: contract.AuditEventKind(item.EventType)})
	}
	return contract.ListAudit200JSONResponse{Items: response, Total: total}, nil
}
