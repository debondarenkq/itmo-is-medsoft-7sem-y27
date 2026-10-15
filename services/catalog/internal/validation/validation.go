package validation

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/google/uuid"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/catalog/internal/models"
)

func RequireID(id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("%w: Требуется ненулевой UUID", models.ErrInvalidInput)
	}
	return nil
}

func Text(field, value string, min, max int) error {
	length := utf8.RuneCountInString(value)
	if !utf8.ValidString(value) || length < min || length > max || (min > 0 && strings.TrimSpace(value) == "") {
		return fmt.Errorf("%w: поле %s — от %d до %d символов", models.ErrInvalidInput, field, min, max)
	}
	return nil
}
