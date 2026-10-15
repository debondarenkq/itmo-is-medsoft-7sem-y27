package api

import (
	"context"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/catalog/api"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/catalog/internal/models"
)

func (implementation *Implementation) ImportDiagnoses(ctx context.Context, request contract.ImportDiagnosesRequestObject) (contract.ImportDiagnosesResponseObject, error) {
	diagnoses := make([]models.Diagnosis, len(request.Body.Diagnoses))
	for index, diagnosis := range request.Body.Diagnoses {
		diagnoses[index] = fromDiagnosisInput(diagnosis)
	}
	imported, err := implementation.Module.Import(ctx, diagnoses)
	if err != nil {
		return nil, err
	}
	items := make([]contract.Diagnosis, len(imported))
	for index, diagnosis := range imported {
		items[index], err = toDiagnosis(diagnosis)
		if err != nil {
			return nil, err
		}
	}
	return contract.ImportDiagnoses200JSONResponse{Items: items}, nil
}
