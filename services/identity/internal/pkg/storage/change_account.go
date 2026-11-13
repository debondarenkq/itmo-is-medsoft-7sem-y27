package storage

import (
	"context"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/internal/models"
)

func (r *Repository) ChangeAccount(ctx context.Context, id string, active *bool, passwordHash, actor string) (models.Account, error) {
	return r.UpdateAccount(ctx, id, models.AccountChanges{IsActive: active, PasswordHash: passwordHash}, actor)
}
