package modules

import (
	"context"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/catalog/internal/models"
)

func (m *Module) Seed(ctx context.Context) ([]models.Diagnosis, error) {
	diagnoses := []models.Diagnosis{
		{Code: "J00.00", Name: "Острый назофарингит"},
		{Code: "I10.00", Name: "Гипертензия"},
		{Code: "E11.00", Name: "Сахарный диабет 2 типа"},
		{Code: "J45.00", Name: "Бронхиальная астма"},
		{Code: "K29.00", Name: "Гастрит"},
	}
	for index, diagnosis := range diagnoses {
		validated, err := Validate(diagnosis)
		if err != nil {
			return nil, err
		}
		diagnoses[index] = validated
	}

	return m.Storage.UpsertBatch(ctx, diagnoses)
}
