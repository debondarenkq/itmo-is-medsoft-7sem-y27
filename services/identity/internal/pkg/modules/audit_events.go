package modules

import (
	"context"
	"fmt"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/internal/models"
)

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
