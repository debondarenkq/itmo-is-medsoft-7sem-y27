package storage

import (
	"context"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func addAudit(ctx context.Context, tx pgx.Tx, actor, message, eventType string) error {
	_, err := tx.Exec(ctx, "INSERT INTO audit_events(event_id,actor,message,event_type) VALUES($1,$2,$3,$4)", uuid.NewString(), actor, message, eventType)
	return err
}

func (r *Repository) AppendAudit(ctx context.Context, actor, message, eventType string) error {
	_, err := r.pool.Exec(ctx, "INSERT INTO audit_events(event_id,actor,message,event_type) VALUES($1,$2,$3,$4)", uuid.NewString(), actor, message, eventType)
	return err
}

func (r *Repository) AuditEvents(ctx context.Context, limit, offset int) ([]models.AuditEvent, int, error) {
	return r.SearchAuditEvents(ctx, models.AuditFilter{Limit: limit, Offset: offset})
}

func (r *Repository) SearchAuditEvents(ctx context.Context, filter models.AuditFilter) ([]models.AuditEvent, int, error) {
	where := " FROM audit_events e WHERE ($1='' OR e.event_type=$1) AND ($2::timestamptz IS NULL OR e.event_time >= $2) AND ($3::timestamptz IS NULL OR e.event_time <= $3)"
	var total int
	if err := r.pool.QueryRow(ctx, "SELECT count(*)"+where, filter.EventType, filter.From, filter.To).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := r.pool.Query(ctx, `SELECT e.event_id,e.event_time,e.actor,e.message,e.event_type,
 COALESCE((SELECT COALESCE(NULLIF(concat_ws(' ',s.last_name,s.first_name),''),a.username) FROM accounts a LEFT JOIN staff s ON s.id=a.staff_id WHERE e.actor='account:'||a.id::text),
 (SELECT 'Система: '||concat_ws(' ',s.last_name,s.first_name) FROM staff s WHERE e.actor='system:staff:'||s.id::text),
 CASE WHEN e.actor LIKE 'anonymous:%' THEN 'Неавторизованный пользователь' ELSE e.actor END)`+where+" ORDER BY e.event_time,e.event_id LIMIT $4 OFFSET $5", filter.EventType, filter.From, filter.To, filter.Limit, filter.Offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	items := []models.AuditEvent{}
	for rows.Next() {
		var event models.AuditEvent
		if err := rows.Scan(&event.EventID, &event.EventTime, &event.Actor, &event.Message, &event.EventType, &event.ActorName); err != nil {
			return nil, 0, err
		}
		items = append(items, event)
	}
	return items, total, rows.Err()
}
