package api

import (
	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/api/security"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/internal/models"
	"github.com/google/uuid"
	"net/http"
)

func (api *SecurityAPI) UpdateAccount(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	var input contract.AccountUpdate
	if !decodeSecurity(w, r, &input) {
		return
	}
	changes := models.AccountChanges{IsActive: &input.IsActive, Username: input.Username}
	if input.Role != nil {
		role := string(*input.Role)
		changes.Role = &role
	}
	claims, _ := api.claims(r)
	account, err := api.Module.UpdateAccount(r.Context(), id.String(), changes, "account:"+claims.Subject)
	if err != nil {
		securityError(w, err)
		return
	}
	writeJSON(w, 200, account)
}
