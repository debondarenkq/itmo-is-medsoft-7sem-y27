package modules

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/internal/models"
	"github.com/google/uuid"
)

func auditRequest(message string) (method, path string) {
	fields := strings.Fields(message)
	if len(fields) >= 3 && (fields[0] == "Запрос" || fields[0] == "Результат") {
		return fields[1], strings.TrimSuffix(fields[2], ":")
	}
	return "", ""
}

func auditEventType(message string) string {
	switch {
	case strings.HasPrefix(message, "Настройка администратора"):
		return "admin_setup"
	case strings.Contains(message, "Неуспеш") || strings.Contains(message, "Превышен лимит"):
		return "login_failed"
	case strings.HasPrefix(message, "Вход "):
		return "login"
	case strings.HasPrefix(message, "Выход "):
		return "logout"
	case strings.HasPrefix(message, "Обновление пары"):
		return "token_refreshed"
	case strings.HasPrefix(message, "Отказ в доступе") || strings.HasPrefix(message, "Неавторизованный"):
		return "access_denied"
	case strings.HasPrefix(message, "Создание аккаунта"):
		return "account_created"
	case strings.HasPrefix(message, "Изменение аккаунта"):
		return "account_updated"
	}
	method, path := auditRequest(message)
	switch {
	case path == "/api/v1/staff" && method == "POST":
		return "staff_created"
	case strings.HasPrefix(path, "/api/v1/staff/") && method == "PUT":
		return "staff_updated"
	case strings.HasPrefix(path, "/api/v1/staff/") && method == "DELETE":
		return "staff_deleted"
	case strings.HasPrefix(path, "/api/v1/records/") && method == "POST":
		return "record_updated"
	case strings.HasPrefix(path, "/api/v1/patients/") && strings.HasSuffix(path, "/record") && method == "POST":
		return "record_created"
	case path == "/api/v1/patients" && method == "POST":
		return "patient_created"
	case strings.HasPrefix(path, "/api/v1/patients/") && method == "PUT":
		return "patient_updated"
	case strings.HasPrefix(path, "/api/v1/patients/") && method == "DELETE":
		return "patient_deleted"
	case path == "/api/v1/accounts" && method == "POST":
		return "account_created"
	case strings.HasPrefix(path, "/api/v1/accounts/"):
		return "account_updated"
	case path == "/api/v1/auth/logout":
		return "logout"
	}
	return "other"
}

func auditDescription(message, kind string) string {
	_, path := auditRequest(message)
	if path == "" {
		return message
	}
	titles := map[string]string{
		"staff_created": "Создание медработника", "staff_updated": "Изменение медработника", "staff_deleted": "Удаление медработника",
		"patient_created": "Создание пациента", "patient_updated": "Изменение пациента", "patient_deleted": "Удаление пациента",
		"record_created": "Создание ЭМК", "record_updated": "Изменение ЭМК",
		"account_created": "Создание аккаунта", "account_updated": "Изменение аккаунта", "logout": "Выход из системы",
	}
	title, ok := titles[kind]
	if !ok {
		return message
	}
	for _, segment := range strings.Split(path, "/") {
		if _, err := uuid.Parse(segment); err == nil {
			title += "; объект " + segment
			break
		}
	}
	if strings.HasPrefix(message, "Запрос ") {
		return "Запрошено: " + title
	}
	fields := strings.Fields(message)
	status, err := strconv.Atoi(fields[len(fields)-1])
	if err != nil {
		return message
	}
	if status >= 400 {
		return fmt.Sprintf("Отклонено: %s; код %d", title, status)
	}
	return "Выполнено: " + title
}

func (s *Security) Audit(ctx context.Context, actor, message string) error {
	kind := auditEventType(message)
	return s.Store.AppendAudit(ctx, actor, auditDescription(message, kind), kind)
}

func (s *Security) AuditEvents(ctx context.Context, filter models.AuditFilter) ([]models.AuditEvent, int, error) {
	if filter.From != nil && filter.To != nil && filter.From.After(*filter.To) {
		return nil, 0, fmt.Errorf("%w: Начало периода позже окончания", models.ErrInvalidInput)
	}
	return s.Store.SearchAuditEvents(ctx, filter)
}
