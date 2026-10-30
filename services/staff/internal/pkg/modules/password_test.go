package modules

import (
	"errors"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/internal/models"
	"testing"
)

func TestAdminPasswordRequirements(t *testing.T) {
	for _, test := range []struct {
		password, confirmation string
		valid                  bool
	}{
		{"Abcdefg1", "Abcdefg1", true},
		{"Пароль12", "Пароль12", true},
		{"abcdefg1", "abcdefg1", false},
		{"ABCDEFG1", "ABCDEFG1", false},
		{"Abcdefgh", "Abcdefgh", false},
		{"Abc1234", "Abc1234", false},
		{"Abcdefg1", "Abcdefg2", false},
	} {
		err := ValidateAdminPassword(test.password, test.confirmation)
		if test.valid && err != nil {
			t.Fatalf("valid password rejected: %v", err)
		}
		if !test.valid && !errors.Is(err, models.ErrInvalidInput) {
			t.Fatalf("invalid password accepted: %v", err)
		}
	}
}

func TestAuditEventClassification(t *testing.T) {
	for message, expected := range map[string]string{
		"Неуспешная попытка входа: user":                        "login_failed",
		"Результат POST /api/v1/patients: HTTP 201":             "patient_created",
		"Запрос POST /api/v1/patients/id/record":                "record_created",
		"Результат POST /api/v1/records/id/changes: HTTP 200":   "record_updated",
		"Результат POST /api/v1/accounts/id/password: HTTP 200": "account_updated",
		"Отказ в доступе GET /api/v1/audit":                     "access_denied",
	} {
		if actual := auditEventType(message); actual != expected {
			t.Fatalf("%q: got %s, expected %s", message, actual, expected)
		}
	}
}
