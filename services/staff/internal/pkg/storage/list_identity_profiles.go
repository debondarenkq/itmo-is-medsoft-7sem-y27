package storage

import (
	"context"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/internal/models"
)

// ListIdentityProfiles lists staff.
func (r *Repository) ListIdentityProfiles(ctx context.Context) ([]models.Staff, error) {
	rows, err := r.pool.Query(ctx, "SELECT "+columns+" FROM staff ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []models.Staff{}
	for rows.Next() {
		s, err := scan(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, s)
	}
	return items, rows.Err()
}
