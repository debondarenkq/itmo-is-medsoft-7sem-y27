package platform

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/api"
)

type Error struct {
	Status  int               `json:"-"`
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitempty"`
}

func (e *Error) Error() string {
	return e.Message
}

func Bad(code, message string) *Error {
	return &Error{Status: 400, Code: code, Message: message}
}

func Conflict(code, message string) *Error {
	return &Error{Status: 409, Code: code, Message: message}
}

func NotFound() *Error {
	return &Error{Status: 404, Code: "NOT_FOUND", Message: "Объект не найден"}
}

type Handler func(http.ResponseWriter, *http.Request) error

func Handle(fn Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
		defer cancel()
		r = r.WithContext(ctx)
		if err := fn(w, r); err != nil {
			WriteError(w, r, err)
		}
	}
}

func WriteError(w http.ResponseWriter, r *http.Request, err error) {
	var api *Error
	if !errors.As(err, &api) {
		if errors.Is(err, pgx.ErrNoRows) {
			api = NotFound()
		} else {
			var db *pgconn.PgError
			if errors.As(err, &db) && db.Code == "23505" {
				api = Conflict("DUPLICATE", "Запись с такими уникальными данными уже существует")
			} else if errors.As(err, &db) && db.Code == "23503" {
				api = Conflict("REFERENCED", "Объект связан с другими записями")
			} else {
				slog.Error("request failed", "method", r.Method, "path", r.URL.Path, "error", err)
				api = &Error{Status: 500, Code: "INTERNAL", Message: "Внутренняя ошибка сервиса"}
			}
		}
	}
	detail := contract.ErrorDetail{Code: api.Code, Message: api.Message}
	if len(api.Fields) > 0 {
		detail.Fields = &api.Fields
	}
	JSON(w, api.Status, contract.Error{Error: detail})
}

func JSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Warn("encode response", "error", err)
	}
}

func ParseID(s string) (string, error) {
	id, err := uuid.Parse(s)
	if err != nil || id == uuid.Nil {
		return "", Bad("INVALID_ID", "Требуется ненулевой UUID")
	}
	return id.String(), nil
}

func Text(field, value string, min, max int) error {
	n := utf8.RuneCountInString(value)
	if !utf8.ValidString(value) || n < min || n > max || (min > 0 && strings.TrimSpace(value) == "") {
		return &Error{Status: 400, Code: "VALIDATION", Message: "Неверное значение поля", Fields: map[string]string{field: fmt.Sprintf("От %d до %d символов", min, max)}}
	}
	return nil
}

type Page struct {
	Limit, Offset  int
	Query          string
	IncludeDeleted bool
}
