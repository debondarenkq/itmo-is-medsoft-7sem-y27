package storage

import (
	"context"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/internal/models"
)

func (r *Repository) AuditEvents(ctx context.Context, limit, offset int) ([]models.AuditEvent, int, error) {
	return r.SearchAuditEvents(ctx, models.AuditFilter{Limit: limit, Offset: offset})
}
