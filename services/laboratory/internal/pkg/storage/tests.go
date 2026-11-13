package storage

import (
	"context"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/laboratory/internal/models"
)

func (r *Repository) Tests(ctx context.Context, inactive bool) ([]models.Test, error) {
	rows, err := r.pool.Query(ctx, "SELECT id,code,name,active FROM lab_tests WHERE active OR $1 ORDER BY name,id", inactive)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []models.Test{}
	for rows.Next() {
		var t models.Test
		if err = rows.Scan(&t.ID, &t.Code, &t.Name, &t.Active); err != nil {
			return nil, err
		}
		result = append(result, t)
	}
	return result, rows.Err()
}
