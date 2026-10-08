package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/pkg/modules"
)

func TestContractRejectsInvalidRequestsBeforeBusinessDependencies(t *testing.T) {
	id := uuid.NewString()
	tests := []struct {
		name   string
		method string
		path   string
		body   string
	}{
		{name: "malformed UUID", method: http.MethodGet, path: "/api/v1/patients/not-a-uuid"},
		{name: "nil UUID", method: http.MethodGet, path: "/api/v1/patients/" + uuid.Nil.String()},
		{name: "missing required date", method: http.MethodPost, path: "/api/v1/patients", body: `{"first_name":"Иван","last_name":"Иванов","administrative_sex":"M"}`},
		{name: "invalid calendar date", method: http.MethodPost, path: "/api/v1/patients", body: `{"first_name":"Иван","last_name":"Иванов","administrative_sex":"M","birth_date":"2001-02-29"}`},
		{name: "unknown field", method: http.MethodPost, path: "/api/v1/patients", body: `{"first_name":"Иван","last_name":"Иванов","administrative_sex":"M","birth_date":"2000-01-01","extra":true}`},
		{name: "unknown enum", method: http.MethodPost, path: "/api/v1/patients", body: `{"first_name":"Иван","last_name":"Иванов","administrative_sex":"other","birth_date":"2000-01-01"}`},
		{name: "null body", method: http.MethodPost, path: "/api/v1/patients", body: `null`},
		{name: "invalid limit", method: http.MethodGet, path: "/api/v1/patients?limit=0"},
		{name: "invalid offset", method: http.MethodGet, path: "/api/v1/patients?offset=-1"},
		{name: "invalid timestamp", method: http.MethodGet, path: "/api/v1/records/" + id + "/state?at=not-a-time"},
		{name: "missing state selector", method: http.MethodGet, path: "/api/v1/records/" + id + "/state"},
		{name: "two state selectors", method: http.MethodGet, path: "/api/v1/records/" + id + "/state?at=2026-10-08T00:00:00Z&version=1"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Any access to these absent dependencies would panic.
			implementation := New(Deps{Module: modules.New(modules.Deps{})})
			router := mux.NewRouter()
			implementation.Register(router)
			request := httptest.NewRequest(test.method, test.path, strings.NewReader(test.body))
			if test.method == http.MethodPost {
				request.Header.Set("Content-Type", "application/json")
			}
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)
			if response.Code != http.StatusBadRequest {
				t.Fatalf("status=%d, response=%s", response.Code, response.Body.String())
			}
			if !strings.Contains(response.Header().Get("Content-Type"), "application/json") {
				t.Fatal("contract error did not use the common JSON envelope")
			}
		})
	}
}
