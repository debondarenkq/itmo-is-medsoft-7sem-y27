package security

import (
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUnexpectedIdentityResponseCannotAuthorizeMutation(t *testing.T) {
	for _, mode := range []string{"missing_claims", "invalid_role", "missing_audit_ack", "unauthorized"} {
		t.Run(mode, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				if r.URL.Path == "/internal/audit" {
					_, _ = w.Write([]byte(`{}`))
					return
				}
				if mode == "unauthorized" {
					w.WriteHeader(401)
					_ = json.NewEncoder(w).Encode(Error{Error: ErrorDetail{Code: UNAUTHORIZED, Message: "Token rejected"}})
					return
				}
				if mode == "missing_claims" {
					_, _ = w.Write([]byte(`{}`))
					return
				}
				role := ClaimsRole("doctor")
				if mode == "invalid_role" {
					role = "invented-role"
				}
				_ = json.NewEncoder(w).Encode(Claims{Sub: uuid.New(), Jti: uuid.New(), Role: role, Iss: "his-identity", Aud: "his-api", Exp: 9999999999, Iat: 1})
			}))
			defer server.Close()
			client, err := NewClientWithResponses(server.URL)
			if err != nil {
				t.Fatal(err)
			}
			authorizer := &Authorizer{client: client}
			mutations := 0
			handler := authorizer.Middleware(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { mutations++ }))
			recorder := httptest.NewRecorder()
			handler.ServeHTTP(recorder, httptest.NewRequest("POST", "/api/v1/staff", nil))
			expected := 503
			code := DEPENDENCYUNAVAILABLE
			if mode == "unauthorized" {
				expected = 401
				code = UNAUTHORIZED
			}
			var problem Error
			if err := json.Unmarshal(recorder.Body.Bytes(), &problem); err != nil {
				t.Fatal(err)
			}
			if recorder.Code != expected || problem.Error.Code != code || mutations != 0 {
				t.Fatalf("unexpected dependency response allowed a mutation: status=%d code=%s mutations=%d", recorder.Code, problem.Error.Code, mutations)
			}
		})
	}
}
