package api

import (
	"net/http"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/catalog/internal/models"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/catalog/internal/pkg/platform"
)

type importInput struct {
	Diagnoses []input `json:"diagnoses"`
}

// Import accepts the contents of the JSON file selected by the desktop client.
func (i *Implementation) Import(w http.ResponseWriter, r *http.Request) error {
	var request importInput
	if err := platform.Decode(w, r, &request); err != nil {
		return err
	}

	diagnoses := make([]models.Diagnosis, len(request.Diagnoses))
	for index, diagnosis := range request.Diagnoses {
		diagnoses[index] = diagnosis.model()
	}
	imported, err := i.Module.Import(r.Context(), diagnoses)
	if err != nil {
		return err
	}

	platform.JSON(w, http.StatusOK, map[string]any{"items": imported})
	return nil
}
