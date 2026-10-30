package api

import (
	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/api/security"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/internal/models"
	"net/http"
	"strings"
)

func (api *SecurityAPI) ListAccounts(w http.ResponseWriter, r *http.Request, params contract.ListAccountsParams) {
	limit, offset := pagination(params.Limit, params.Offset)
	filter := models.AccountFilter{Limit: limit, Offset: offset, IsActive: params.IsActive}
	if params.Search != nil {
		filter.Search = strings.TrimSpace(*params.Search)
	}
	items, total, err := api.Module.Store.SearchAccounts(r.Context(), filter)
	if err != nil {
		securityError(w, err)
		return
	}
	writeJSON(w, 200, map[string]any{"items": items, "total": total})
}
