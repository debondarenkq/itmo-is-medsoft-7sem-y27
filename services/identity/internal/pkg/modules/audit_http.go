package modules

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/internal/models"
	"github.com/google/uuid"
)

// AuditHTTP records requests.
func (s *Security) AuditHTTP(ctx context.Context, input models.HTTPAudit) error {
	status := 0
	if input.Status != nil {
		status = *input.Status
	}
	if !input.Method.Valid() || !input.EventType.Valid() || input.Title == "" || input.Actor == "" {
		return models.ErrInvalidInput
	}
	title := input.Title
	for _, segment := range strings.Split(input.Path, "/") {
		if _, err := uuid.Parse(segment); err == nil {
			title += "; объект " + segment
			break
		}
	}
	message := "Запрошено: " + title
	if status >= http.StatusBadRequest {
		message = fmt.Sprintf("Отклонено: %s; код %d", title, status)
	} else if status != 0 {
		message = "Выполнено: " + title
	}
	return s.Audit(ctx, input.Actor, message, input.EventType)
}
