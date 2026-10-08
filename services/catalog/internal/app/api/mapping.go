package api

import (
	"github.com/google/uuid"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/catalog/api"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/catalog/internal/models"
)

func fromDiagnosisInput(value contract.DiagnosisInput) models.Diagnosis {
	return models.Diagnosis{Code: value.Code, Name: value.Name}
}
func toDiagnosis(value models.Diagnosis) contract.Diagnosis {
	return contract.Diagnosis{ID: uuid.MustParse(value.ID), Code: value.Code, Name: value.Name, CreatedAt: value.CreatedAt, UpdatedAt: value.UpdatedAt, DeletedAt: value.DeletedAt}
}
