package storage

import (
	"context"

	"github.com/google/uuid"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/catalog/internal/models"
)

func (r *Repository) Seed(ctx context.Context) ([]models.Diagnosis, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	samples := []models.Diagnosis{
		{Code: "J00.00", Name: "Острый назофарингит"},
		{Code: "I10.00", Name: "Гипертензия"},
		{Code: "E11.00", Name: "Сахарный диабет 2 типа"},
		{Code: "J45.00", Name: "Бронхиальная астма"},
		{Code: "K29.00", Name: "Гастрит"},
	}
	result := []models.Diagnosis{}
	for _, v := range samples {
		item, err := scan(tx.QueryRow(ctx, "INSERT INTO diagnoses(id,code,name) VALUES($1,$2,$3) ON CONFLICT(code) DO UPDATE SET name=EXCLUDED.name,deleted_at=NULL,updated_at=clock_timestamp() RETURNING "+columns, uuid.NewString(), v.Code, v.Name))
		if err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}
	return result, nil
}
