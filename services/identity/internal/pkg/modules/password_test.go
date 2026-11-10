package modules

import (
	"context"
	"errors"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/internal/models"
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

type auditCapture struct {
	SecurityStorage
	kind    models.AuditEventKind
	message string
	err     error
}

func (s *auditCapture) AppendAudit(_ context.Context, _, message string, kind models.AuditEventKind) error {
	s.kind = kind
	s.message = message
	return s.err
}

func TestAuditTypeDoesNotDependOnMessage(t *testing.T) {
	store := &auditCapture{}
	security := &Security{Store: store}
	for _, message := range []string{"Неуспешная попытка входа", "Authentication failed", "Текст полностью изменён"} {
		if err := security.Audit(context.Background(), "anonymous:test", message, models.AuditEventKindLoginFailed); err != nil {
			t.Fatal(err)
		}
		if store.kind != models.AuditEventKindLoginFailed || store.message != message {
			t.Fatal("event type or message changed")
		}
	}
}

func TestAuditHTTPPreservesOperationAndOutcome(t *testing.T) {
	store := &auditCapture{}
	security := &Security{Store: store}
	id := "00000000-0000-4000-8000-000000000001"
	for _, test := range []struct {
		method, path string
		status       int
		kind         models.AuditEventKind
		title        string
		message      string
	}{
		{"POST", "/api/v1/auth/logout", 200, models.AuditEventKindLogout, "Выход из системы", "Выполнено: Выход из системы"},
		{"POST", "/api/v1/patients", 0, models.AuditEventKindPatientCreated, "Создание пациента", "Запрошено: Создание пациента"},
		{"POST", "/api/v1/patients", 201, models.AuditEventKindPatientCreated, "Создание пациента", "Выполнено: Создание пациента"},
		{"POST", "/api/v1/patients", 409, models.AuditEventKindPatientCreated, "Создание пациента", "Отклонено: Создание пациента; код 409"},
		{"POST", "/api/v1/patients/" + id + "/record", 201, models.AuditEventKindRecordCreated, "Создание ЭМК", "Выполнено: Создание ЭМК; объект " + id},
		{"POST", "/api/v1/records/" + id + "/changes", 200, models.AuditEventKindRecordUpdated, "Изменение ЭМК", "Выполнено: Изменение ЭМК; объект " + id},
		{"POST", "/api/v1/accounts/" + id + "/password", 200, models.AuditEventKindAccountUpdated, "Изменение аккаунта", "Выполнено: Изменение аккаунта; объект " + id},
		{"POST", "/api/v1/staffing", 404, models.AuditEventKindOther, "POST /api/v1/staffing", "Отклонено: POST /api/v1/staffing; код 404"},
	} {
		if err := security.AuditHTTP(context.Background(), models.HTTPAudit{Actor: "account:test", Method: models.HTTPAuditMethod(test.method), Path: test.path, Status: &test.status, EventType: test.kind, Title: test.title}); err != nil {
			t.Fatal(err)
		}
		if store.kind != test.kind || store.message != test.message {
			t.Fatalf("%s %s: kind=%s message=%q", test.method, test.path, store.kind, store.message)
		}
	}
	store.err = errors.New("audit unavailable")
	if err := security.AuditHTTP(context.Background(), models.HTTPAudit{Actor: "account:test", Method: models.HTTPAuditMethodPOST, Path: "/api/v1/patients", EventType: models.AuditEventKindPatientCreated, Title: "Создание пациента"}); !errors.Is(err, store.err) {
		t.Fatal("audit failure ignored")
	}
}
