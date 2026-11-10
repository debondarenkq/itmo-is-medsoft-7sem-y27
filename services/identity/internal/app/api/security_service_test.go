package api

import (
	"context"
	"encoding/base32"
	"encoding/json"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	auditcontract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/api/audit"
	internalapi "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/api/internalapi"
	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/api/security"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/internal/models"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/internal/pkg/modules"
	"github.com/gorilla/mux"
)

type contractStore struct {
	modules.SecurityStorage
	secret     string
	configured bool
	events     int
}

func (s *contractStore) Account(context.Context, string) (models.Account, error) {
	return models.Account{}, models.ErrUnauthorized
}

func (s *contractStore) SetupSecret(_ context.Context, candidate string) (bool, string, error) {
	if s.configured {
		return true, "", nil
	}
	if s.secret == "" {
		s.secret = candidate
	}
	return false, s.secret, nil
}
func (s *contractStore) AppendAudit(_ context.Context, _, _ string, _ models.AuditEventKind) error {
	s.events++
	return nil
}

func TestStrictContractsAndServerGeneratedTOTP(t *testing.T) {
	store := &contractStore{}
	router := mux.NewRouter()
	New(Deps{Module: &modules.Security{Store: store}, InternalKey: "test-service-key"}).Register(router)
	request := func(method, path, body, key string) *httptest.ResponseRecorder {
		r := httptest.NewRequest(method, path, strings.NewReader(body))
		r.Header.Set("Content-Type", "application/json")
		if key != "" {
			r.Header.Set("X-Service-Key", key)
		}
		response := httptest.NewRecorder()
		router.ServeHTTP(response, r)
		return response
	}
	first := request("GET", "/api/v1/auth/setup", "", "")
	if first.Code != 200 {
		t.Fatalf("setup status: %d", first.Code)
	}
	var setup contract.SetupStatus
	if err := json.Unmarshal(first.Body.Bytes(), &setup); err != nil {
		t.Fatal(err)
	}
	if setup.Configured || setup.Secret == nil || setup.ProvisioningURI == nil {
		t.Fatal("missing server-generated provisioning data")
	}
	key, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(*setup.Secret)
	if err != nil || len(key) != 20 {
		t.Fatal("TOTP secret must contain 160 random bits")
	}
	uri, err := url.Parse(*setup.ProvisioningURI)
	if err != nil {
		t.Fatal(err)
	}
	if uri.Scheme != "otpauth" || uri.Host != "totp" || uri.Path != "/HIS:admin" || uri.Query().Get("secret") != *setup.Secret || uri.Query().Get("issuer") != "HIS" || uri.Query().Get("algorithm") != "SHA1" || uri.Query().Get("period") != "30" || uri.Query().Get("digits") != "6" {
		t.Fatal("invalid TOTP provisioning URI")
	}
	again := request("GET", "/api/v1/auth/setup", "", "")
	var repeated contract.SetupStatus
	if err := json.Unmarshal(again.Body.Bytes(), &repeated); err != nil {
		t.Fatal(err)
	}
	if repeated.Secret == nil || *repeated.Secret != *setup.Secret || repeated.ProvisioningURI == nil || *repeated.ProvisioningURI != *setup.ProvisioningURI {
		t.Fatal("setup secret changed before confirmation")
	}
	for _, body := range []string{`{"username":"doctor"}`, `{"username":"doctor","password":"x","unexpected":true}`, `{"username":"doctor","password":"x"}{}`} {
		response := request("POST", "/api/v1/auth/login", body, "")
		var problem contract.SecurityError
		if err := json.Unmarshal(response.Body.Bytes(), &problem); err != nil {
			t.Fatal(err)
		}
		if response.Code != 400 || problem.Error.Code != contract.INVALIDREQUEST {
			t.Fatalf("malformed public request accepted: %d body=%q", response.Code, body)
		}
	}
	denied := request("POST", "/internal/audit", "not-json", "")
	var privateError internalapi.Error
	if err := json.Unmarshal(denied.Body.Bytes(), &privateError); err != nil {
		t.Fatal(err)
	}
	if denied.Code != 403 || privateError.Error.Code != internalapi.FORBIDDEN {
		t.Fatal("private authorization must precede parsing")
	}
	invalid := request("POST", "/internal/audit", `{"actor":"test","method":"GET","path":"/api/v1/patients"}`, "test-service-key")
	if invalid.Code != 400 || store.events != 0 {
		t.Fatal("invalid internal enum reached business logic")
	}
	valid := request("POST", "/internal/audit", `{"actor":"test","method":"POST","path":"/api/v1/patients","status":201,"event_type":"patient_created","title":"Создание пациента"}`, "test-service-key")
	var acknowledgement auditcontract.Acknowledgement
	if err := json.Unmarshal(valid.Body.Bytes(), &acknowledgement); err != nil {
		t.Fatal(err)
	}
	if valid.Code != 200 || !acknowledgement.Ok || store.events != 1 {
		t.Fatal("typed internal request not recorded")
	}
	unauthenticated := request("POST", "/internal/authorize", `{"method":"GET","path":"/api/v1/patients"}`, "test-service-key")
	if err := json.Unmarshal(unauthenticated.Body.Bytes(), &privateError); err != nil {
		t.Fatal(err)
	}
	if unauthenticated.Code != 401 || privateError.Error.Code != internalapi.UNAUTHORIZED {
		t.Fatal("missing Bearer must produce the documented error")
	}
	store.configured = true
	configured := request("GET", "/api/v1/auth/setup", "", "")
	var complete contract.SetupStatus
	if err := json.Unmarshal(configured.Body.Bytes(), &complete); err != nil {
		t.Fatal(err)
	}
	if !complete.Configured || complete.Secret != nil || complete.ProvisioningURI != nil {
		t.Fatal("configured admin leaked provisioning secrets")
	}
	if configured.Header().Get("Cache-Control") != "no-store" {
		t.Fatal("setup response may be cached")
	}
}
