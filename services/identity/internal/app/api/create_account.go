package api

import (
	"net/http"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/api/security"
)

func (api *SecurityAPI) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var input contract.AccountInput
	if !decodeSecurity(w, r, &input) {
		return
	}
	claims, _ := api.claims(r)
	a, password, err := api.Module.CreateAccount(r.Context(), input.Username, input.StaffId.String(), "account:"+claims.Subject)
	if err != nil {
		securityError(w, err)
		return
	}
	writeJSON(w, 201, map[string]any{"account": a, "password": password})
}
