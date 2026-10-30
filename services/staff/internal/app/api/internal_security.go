package api

import (
	"net/http"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/internal/models"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/internal/pkg/modules"
)

func (api *SecurityAPI) internal(w http.ResponseWriter, r *http.Request) bool {
	if api.InternalKey == "" || r.Header.Get("X-Service-Key") != api.InternalKey {
		securityError(w, models.ErrForbidden)
		return false
	}
	return true
}

func (api *SecurityAPI) authorize(w http.ResponseWriter, r *http.Request) {
	if !api.internal(w, r) {
		return
	}
	var input struct {
		Method string `json:"method"`
		Path   string `json:"path"`
	}
	if !decodeSecurity(w, r, &input) {
		return
	}
	claims, err := api.claims(r)
	if err != nil {
		api.failed(w, r, err, "Неавторизованный запрос "+input.Method+" "+input.Path)
		return
	}
	if !modules.Allowed(claims, input.Method, input.Path) {
		if err := api.Module.Audit(r.Context(), "account:"+claims.Subject, "Отказ в доступе "+input.Method+" "+input.Path); err != nil {
			securityError(w, err)
			return
		}
		securityError(w, models.ErrForbidden)
		return
	}
	writeJSON(w, 200, claims)
}

func (api *SecurityAPI) audit(w http.ResponseWriter, r *http.Request) {
	if !api.internal(w, r) {
		return
	}
	var input struct {
		Actor   string `json:"actor"`
		Message string `json:"message"`
	}
	if !decodeSecurity(w, r, &input) {
		return
	}
	if input.Actor == "" || len(input.Message) > 2048 {
		securityError(w, models.ErrInvalidInput)
		return
	}
	if err := api.Module.Audit(r.Context(), input.Actor, input.Message); err != nil {
		securityError(w, err)
		return
	}
	writeJSON(w, 200, map[string]bool{"ok": true})
}
