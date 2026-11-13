package storage

import (
	"context"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/internal/models"
)

func (r *Repository) ServiceAccounts(ctx context.Context) ([]models.Account, error) {
	rows, err := r.pool.Query(ctx, "SELECT "+accountColumns+accountFrom+"WHERE a.is_service ORDER BY a.created_at,a.id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []models.Account{}
	for rows.Next() {
		a, err := scanAccount(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, a)
	}
	return result, rows.Err()
}
