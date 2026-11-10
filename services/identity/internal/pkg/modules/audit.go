package modules

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/internal/models"
	"github.com/google/uuid"
)

// Audit records events.
func (s *Security) Audit(ctx context.Context, actor, message string, kind models.AuditEventKind) error {
	return s.Store.AppendAudit(ctx, actor, message, kind)
}

// AuditHTTP records requests.
func (s *Security) AuditHTTP(ctx context.Context, input models.HTTPAudit) error {
	status := 0
	if input.Status != nil {
		status = *input.Status
	}
	if !input.Method.Valid() || !input.EventType.Valid() || input.Title == "" || input.Actor == "" {
		return models.ErrInvalidInput
	}
	title := input.Title
	for _, segment := range strings.Split(input.Path, "/") {
		if _, err := uuid.Parse(segment); err == nil {
			title += "; объект " + segment
			break
		}
	}
	message := "Запрошено: " + title
	if status >= http.StatusBadRequest {
		message = fmt.Sprintf("Отклонено: %s; код %d", title, status)
	} else if status != 0 {
		message = "Выполнено: " + title
	}
	return s.Audit(ctx, input.Actor, message, input.EventType)
}

func (s *Security) AuditEvents(ctx context.Context, filter models.AuditFilter) ([]models.AuditEvent, int, error) {
	if filter.From != nil && filter.To != nil && filter.From.After(*filter.To) {
		return nil, 0, fmt.Errorf("%w: Начало периода позже окончания", models.ErrInvalidInput)
	}
	events, total, err := s.Store.SearchAuditEvents(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	names := map[string]string{}
	for i := range events {
		event := &events[i]
		id := ""
		if event.ActorStaffID != nil {
			id = *event.ActorStaffID
		}
		if id == "" {
			continue
		}
		name, ok := names[id]
		if !ok {
			p, err := s.Directory.Get(ctx, id)
			if err != nil {
				return nil, 0, err
			}
			name = staffName(p)
			names[id] = name
		}
		event.ActorName = name
	}
	return events, total, nil
}
