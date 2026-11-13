package modules

import (
	"context"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/internal/models"
)

func (s *Security) ServiceAccounts(ctx context.Context) ([]models.Account, error) {
	return s.Store.ServiceAccounts(ctx)
}
