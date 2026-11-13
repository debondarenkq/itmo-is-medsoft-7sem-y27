package modules

import (
	"context"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/internal/models"
)

// Audit records events.
func (s *Security) Audit(ctx context.Context, actor, message string, kind models.AuditEventKind) error {
	return s.Store.AppendAudit(ctx, actor, message, kind)
}
