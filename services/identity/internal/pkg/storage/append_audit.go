package storage

import (
	"context"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) AppendAudit(ctx context.Context, actor, message string, eventType models.AuditEventKind) error {
	_, err := r.pool.Exec(ctx, "INSERT INTO audit_events(event_id,actor,message,event_type) VALUES($1,$2,$3,$4)", uuid.NewString(), actor, message, eventType)
	return err
}

func addAudit(ctx context.Context, tx pgx.Tx, actor, message string, eventType models.AuditEventKind) error {
	_, err := tx.Exec(ctx, "INSERT INTO audit_events(event_id,actor,message,event_type) VALUES($1,$2,$3,$4)", uuid.NewString(), actor, message, eventType)
	return err
}
