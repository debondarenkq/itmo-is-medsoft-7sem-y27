package api

import (
	"net/http"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/models"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/pkg/platform"
)

func (i *Implementation) CreateRecord(w http.ResponseWriter, r *http.Request) error {
	id, err := platform.ID(r)
	if err != nil {
		return err
	}
	var in models.CreateRecordInput
	if err = platform.Decode(w, r, &in); err != nil {
		return err
	}
	v, err := i.Module.CreateRecord(r.Context(), id, in.StaffID)
	if err != nil {
		return err
	}
	platform.JSON(w, 201, v)
	return nil
}
