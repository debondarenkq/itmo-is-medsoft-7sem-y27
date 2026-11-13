package storage

import (
	"context"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/internal/models"
)

func (r *Repository) SearchAuditEvents(ctx context.Context, filter models.AuditFilter) ([]models.AuditEvent, int, error) {
	where := " FROM audit_events e WHERE ($1='' OR e.event_type=$1) AND ($2::timestamptz IS NULL OR e.event_time >= $2) AND ($3::timestamptz IS NULL OR e.event_time <= $3)"
	var total int
	if err := r.pool.QueryRow(ctx, "SELECT count(*)"+where, filter.EventType, filter.From, filter.To).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := r.pool.Query(ctx, `SELECT e.event_id,e.event_time,e.actor,e.message,e.event_type,
 COALESCE((SELECT a.username FROM accounts a WHERE e.actor='account:'||a.id::text),
 CASE WHEN e.actor LIKE 'anonymous:%' THEN 'Неавторизованный пользователь' ELSE e.actor END),
 (SELECT a.staff_id FROM accounts a WHERE e.actor='account:'||a.id::text)`+where+" ORDER BY e.event_time,e.event_id LIMIT $4 OFFSET $5", filter.EventType, filter.From, filter.To, filter.Limit, filter.Offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	items := []models.AuditEvent{}
	for rows.Next() {
		var event models.AuditEvent
		if err := rows.Scan(&event.EventID, &event.EventTime, &event.Actor, &event.Message, &event.EventType, &event.ActorName, &event.ActorStaffID); err != nil {
			return nil, 0, err
		}
		items = append(items, event)
	}
	return items, total, rows.Err()
}
