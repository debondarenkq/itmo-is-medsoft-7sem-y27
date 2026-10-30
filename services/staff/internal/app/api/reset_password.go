package api

import (
	"net/http"

	"github.com/google/uuid"
)

func (api *SecurityAPI) ResetPassword(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	claims, _ := api.claims(r)
	a, password, err := api.Module.ResetPassword(r.Context(), id.String(), "account:"+claims.Subject)
	if err != nil {
		securityError(w, err)
		return
	}
	writeJSON(w, 200, map[string]any{"account": a, "password": password})
}
