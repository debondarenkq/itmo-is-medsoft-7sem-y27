package modules

import (
	"context"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/internal/models"
)

func (m *Module) ListIdentityProfiles(ctx context.Context) ([]models.Staff, error) {
	return m.Storage.ListIdentityProfiles(ctx)
}
