package api

import (
	"net/http"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/api/security"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/internal/models"
)

func (api *SecurityAPI) GetSetup(w http.ResponseWriter, r *http.Request) {
	configured, secret, err := api.Module.SetupStatus(r.Context())
	if err != nil {
		securityError(w, err)
		return
	}
	writeJSON(w, 200, map[string]any{"configured": configured, "secret": secret})
}

func (api *SecurityAPI) SetupAdmin(w http.ResponseWriter, r *http.Request) {
	if api.limited(r) {
		api.failed(w, r, models.ErrUnauthorized, "Превышен лимит настройки администратора")
		return
	}
	var input contract.SetupInput
	if !decodeSecurity(w, r, &input) {
		return
	}
	session, err := api.Module.Setup(r.Context(), input.Password, input.Confirmation, input.Secret, input.Totp)
	if err != nil {
		api.failed(w, r, err, "Неуспешная настройка администратора")
		return
	}
	writeJSON(w, 200, session)
}
