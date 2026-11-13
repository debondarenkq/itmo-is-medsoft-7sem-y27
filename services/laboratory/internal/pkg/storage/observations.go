package storage

import (
	"context"
	"encoding/json"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/laboratory/internal/models"
)

func (r *Repository) Observations(ctx context.Context, patient, id string) ([]models.Observation, error) {
	rows, err := r.pool.Query(ctx, "SELECT resource FROM lab_observations WHERE ($1='' OR patient_id::text=$1) AND ($2='' OR id::text=$2) ORDER BY created_at,id", patient, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []models.Observation{}
	for rows.Next() {
		var data []byte
		var item models.Observation
		if err = rows.Scan(&data); err != nil {
			return nil, err
		}
		if err = json.Unmarshal(data, &item); err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return result, rows.Err()
}
