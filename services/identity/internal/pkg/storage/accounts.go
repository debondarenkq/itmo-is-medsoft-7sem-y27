package storage

import (
	"context"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/internal/models"
)

func (r *Repository) Accounts(ctx context.Context, limit, offset int) ([]models.Account, int, error) {
	return r.SearchAccounts(ctx, models.AccountFilter{Limit: limit, Offset: offset})
}
