package api

import (
	"net"
	"net/http"
	"time"

	"github.com/google/uuid"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/api/security"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/internal/models"
)

func (api *SecurityAPI) limited(r *http.Request) bool {
	api.mu.Lock()
	defer api.mu.Unlock()
	now := time.Now()
	for key, window := range api.attempts {
		if now.Sub(window.start) > time.Minute {
			delete(api.attempts, key)
		}
	}
	host, _, _ := net.SplitHostPort(r.RemoteAddr)
	// Global cap also covers repeated requests behind the gateway.
	window := api.attempts[host]
	if window.start.IsZero() {
		window.start = now
	}
	window.count++
	api.attempts[host] = window
	return window.count > 20
}

func (api *SecurityAPI) failed(w http.ResponseWriter, r *http.Request, err error, message string) {
	if auditErr := api.Module.Audit(r.Context(), "anonymous:"+r.RemoteAddr+":request:"+uuid.NewString(), message); auditErr != nil {
		securityError(w, auditErr)
		return
	}
	securityError(w, err)
}

func (api *SecurityAPI) Login(w http.ResponseWriter, r *http.Request) {
	if api.limited(r) {
		api.failed(w, r, models.ErrUnauthorized, "Превышен лимит попыток входа")
		return
	}
	var input contract.LoginInput
	if !decodeSecurity(w, r, &input) {
		return
	}
	code := ""
	if input.Totp != nil {
		code = *input.Totp
	}
	session, err := api.Module.Login(r.Context(), input.Username, input.Password, code)
	if err != nil {
		api.failed(w, r, err, "Неуспешная попытка входа: "+input.Username)
		return
	}
	writeJSON(w, 200, session)
}
