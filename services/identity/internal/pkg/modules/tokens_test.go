package modules

import (
	"crypto/ed25519"
	"crypto/rand"
	"strings"
	"testing"
	"time"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/internal/models"
)

func TestTOTPReference(t *testing.T) {
	// RFC 6238 Appendix B, SHA-1, reduced to six digits.
	secret := "GEZDGNBVGY3TQOJQGEZDGNBVGY3TQOJQ"
	for _, test := range []struct {
		seconds int64
		code    string
	}{{59, "287082"}, {1111111109, "081804"}, {1111111111, "050471"}, {1234567890, "005924"}, {2000000000, "279037"}, {20000000000, "353130"}} {
		if step, ok := TOTPStep(secret, test.code, time.Unix(test.seconds, 0)); !ok || step != test.seconds/30 {
			t.Fatalf("RFC vector %d: step=%d valid=%v", test.seconds, step, ok)
		}
	}
	for _, code := range []string{"", "12345", "abcdef", "0000000"} {
		if _, ok := TOTPStep(secret, code, time.Unix(59, 0)); ok {
			t.Fatalf("accepted invalid code %q", code)
		}
	}
}
func TestSignedAccessTokens(t *testing.T) {
	_, key, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	security := &Security{key: key}
	session, err := security.sign(models.Account{ID: "account", Role: "doctor"})
	if err != nil {
		t.Fatal(err)
	}
	claims, err := security.Verify(session.AccessToken)
	if err != nil || claims.Subject != "account" || claims.Role != "doctor" {
		t.Fatalf("verify: %v %+v", err, claims)
	}
	parts := strings.Split(session.AccessToken, ".")
	parts[1] = "e30"
	if _, err := security.Verify(strings.Join(parts, ".")); err == nil {
		t.Fatal("tampered token accepted")
	}
	_, other, _ := ed25519.GenerateKey(rand.Reader)
	if _, err := (&Security{key: other}).Verify(session.AccessToken); err == nil {
		t.Fatal("foreign signature accepted")
	}
}
func TestRolePermissions(t *testing.T) {
	for _, role := range []string{"doctor", "nurse", "laboratory-technician", "other"} {
		c := models.Claims{Subject: "account", StaffID: "staff", Role: role}
		for _, path := range []string{"/api/v1/patients", "/api/v1/records/id", "/api/v1/records/id/history"} {
			if !Allowed(c, "GET", path) {
				t.Fatalf("%s cannot view %s", role, path)
			}
		}
		if Allowed(c, "POST", "/api/v1/patients") != (role == "doctor" || role == "nurse") {
			t.Fatalf("patient creation: %s", role)
		}
		if Allowed(c, "POST", "/api/v1/patients/id/record") != (role == "doctor" || role == "nurse") {
			t.Fatalf("record creation: %s", role)
		}
		if Allowed(c, "POST", "/api/v1/records/id/changes") != (role == "doctor") {
			t.Fatalf("record changes: %s", role)
		}
		if Allowed(c, "GET", "/api/v1/audit") || Allowed(c, "POST", "/api/v1/staff") {
			t.Fatalf("admin access: %s", role)
		}
	}
	if Allowed(models.Claims{Role: "doctor"}, "GET", "/api/v1/patients") {
		t.Fatal("unlinked account accepted")
	}
	if Allowed(models.Claims{Role: "admin"}, "GET", "/api/v1/patients") {
		t.Fatal("administrator entered clinical UI")
	}
}

func TestLISRolePermissions(t *testing.T) {
	id := "30000000-0000-4000-8000-000000000001"
	for _, test := range []struct {
		name, role, method, path string
		allowed                  bool
	}{
		{"patient", "lis", "GET", "/api/v1/fhir/Patient/" + id, true},
		{"practitioner", "lis", "GET", "/api/v1/fhir/Practitioner/" + id, true},
		{"requests", "lis", "GET", "/api/v1/fhir/ServiceRequest", true},
		{"results", "lis", "GET", "/api/v1/fhir/Observation", true},
		{"write result", "lis", "POST", "/api/v1/fhir/Observation", true},
		{"deny request creation", "lis", "POST", "/api/v1/lab/orders", false},
		{"deny result editing", "lis", "PUT", "/api/v1/fhir/Observation/" + id, false},
		{"deny patient editing", "lis", "PUT", "/api/v1/fhir/Patient/" + id, false},
		{"deny raw clinical", "lis", "GET", "/api/v1/patients", false},
		{"deny account administration", "lis", "GET", "/api/v1/service-accounts", false},
		{"admin service accounts", "admin", "GET", "/api/v1/service-accounts", true},
		{"admin directory", "admin", "POST", "/api/v1/lab/tests", true},
		{"deny admin medical access", "admin", "GET", "/api/v1/fhir/ServiceRequest", false},
		{"doctor orders", "doctor", "POST", "/api/v1/lab/orders", true},
		{"nurse cannot order", "nurse", "POST", "/api/v1/lab/orders", false},
		{"worker reads results", "other", "GET", "/api/v1/fhir/Observation", true},
		{"doctor cannot submit results", "doctor", "POST", "/api/v1/fhir/Observation", false},
	} {
		t.Run(test.name, func(t *testing.T) {
			c := models.Claims{Role: test.role}
			if test.role != "lis" && test.role != "admin" {
				c.StaffID = id
			}
			if got := Allowed(c, test.method, test.path); got != test.allowed {
				t.Fatalf("allowed=%v, want %v", got, test.allowed)
			}
		})
	}
}
