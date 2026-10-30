package api

import (
	"net/http"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/api/security"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/internal/pkg/modules"
)

func (api *SecurityAPI) Logout(w http.ResponseWriter, r *http.Request) {
	var input contract.TokenInput
	if !decodeSecurity(w, r, &input) {
		return
	}
	claims, err := api.claims(r)
	if err == nil {
		err = api.Module.Store.Logout(r.Context(), claims.Subject, modules.TokenHash(input.RefreshToken))
	}
	if err != nil {
		securityError(w, err)
		return
	}
	writeJSON(w, 200, map[string]bool{"configured": true})
}
