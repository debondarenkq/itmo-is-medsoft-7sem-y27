package api

import (
	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/api/security"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/internal/models"
	"net/http"
)

func (api *SecurityAPI) ListAudit(w http.ResponseWriter, r *http.Request, params contract.ListAuditParams) {
	limit, offset := pagination(params.Limit, params.Offset)
	filter := models.AuditFilter{Limit: limit, Offset: offset, From: params.From, To: params.To}
	if params.EventType != nil {
		filter.EventType = string(*params.EventType)
	}
	items, total, err := api.Module.AuditEvents(r.Context(), filter)
	if err != nil {
		securityError(w, err)
		return
	}
	writeJSON(w, 200, map[string]any{"items": items, "total": total})
}
