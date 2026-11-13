package modules

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
	"sync"

	security "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/api/security"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/internal/models"
	laboratory "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/internal/pkg/adapters/laboratoryapi"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/routers"
	openapirouter "github.com/getkin/kin-openapi/routers/gorillamux"
)

var laboratoryPolicy = sync.OnceValue(func() []routers.Router {
	result := []routers.Router{}
	for _, load := range []func() (*openapi3.T, error){laboratory.GetSwagger, security.GetSwagger} {
		spec, err := load()
		if err != nil {
			panic(err)
		}
		spec.Servers = nil
		router, err := openapirouter.NewRouter(spec)
		if err != nil {
			panic(err)
		}
		result = append(result, router)
	}
	return result
})

func laboratoryAllowed(c models.Claims, method, path string) (bool, bool) {
	r, err := http.NewRequest(method, "http://policy"+path, nil)
	if err != nil {
		return false, true
	}
	for _, router := range laboratoryPolicy() {
		route, _, err := router.FindRoute(r)
		if err != nil {
			continue
		}
		value, ok := route.Operation.Extensions["x-roles"]
		if !ok {
			continue
		}
		var roles []string
		data, err := json.Marshal(value)
		if err != nil {
			return false, true
		}
		if json.Unmarshal(data, &roles) != nil {
			return false, true
		}
		return slices.Contains(roles, c.Role) && (c.Role == models.ServiceRole || c.Role == "admin" || c.StaffID != ""), true
	}
	return false, false
}

func Allowed(claims models.Claims, method, path string) bool {
	if allowed, matched := laboratoryAllowed(claims, method, path); matched {
		return allowed
	}
	if claims.Role == "admin" {
		return strings.HasPrefix(path, "/api/v1/staff") || strings.HasPrefix(path, "/api/v1/diagnoses") || strings.HasPrefix(path, "/api/v1/accounts") || path == "/api/v1/audit" || path == "/api/v1/auth/logout"
	}
	if claims.StaffID == "" {
		return false
	}
	if path == "/api/v1/auth/logout" {
		return true
	}
	if method == "GET" {
		return strings.HasPrefix(path, "/api/v1/patients") || strings.HasPrefix(path, "/api/v1/records") || strings.HasPrefix(path, "/api/v1/staff/") || strings.HasPrefix(path, "/api/v1/diagnoses")
	}
	if strings.HasPrefix(path, "/api/v1/records/") {
		return claims.Role == "doctor"
	}
	if strings.HasPrefix(path, "/api/v1/patients") {
		return claims.Role == "doctor" || claims.Role == "nurse"
	}
	return false
}
