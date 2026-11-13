package storage

import (
	"context"
	"encoding/json"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/lis/internal/models"
)

func (r *Repository) History(ctx context.Context, search string) ([]models.History, error) {
	rows, err := r.Pool.Query(ctx, "SELECT snapshot,delivered FROM results WHERE $1='' OR snapshot::text ILIKE '%'||$1||'%' ORDER BY created_at DESC,id", search)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []models.History{}
	for rows.Next() {
		var h models.History
		var data []byte
		var delivered bool
		if err = rows.Scan(&data, &delivered); err != nil {
			return nil, err
		}
		if err = json.Unmarshal(data, &h); err != nil {
			return nil, err
		}
		h.Delivered = delivered
		result = append(result, h)
	}
	return result, rows.Err()
}
