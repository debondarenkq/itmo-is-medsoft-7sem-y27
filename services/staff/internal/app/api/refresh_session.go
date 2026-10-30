package api

import (
	"net/http"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/api/security"
)

func (api *SecurityAPI) RefreshSession(w http.ResponseWriter, r *http.Request) {
	var input contract.TokenInput
	if !decodeSecurity(w, r, &input) {
		return
	}
	session, err := api.Module.Refresh(r.Context(), input.RefreshToken)
	if err != nil {
		api.failed(w, r, err, "Отклонено обновление пары токенов")
		return
	}
	writeJSON(w, 200, session)
}
