package platform

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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
	JSON(w, api.Status, map[string]any{"error": api})
}

func JSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Warn("encode response", "error", err)
	}
}

func Decode(w http.ResponseWriter, r *http.Request, v any) error {
	media, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || media != "application/json" {
		return &Error{Status: 415, Code: "CONTENT_TYPE", Message: "Требуется Content-Type: application/json"}
	}
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(v); err != nil {
		return Bad("INVALID_JSON", "Некорректный JSON: "+err.Error())
	}
	if err := d.Decode(new(any)); err != io.EOF {
		return Bad("INVALID_JSON", "Ожидается один JSON-объект")
	}
	return nil
}

func ID(r *http.Request) (string, error) {
	return ParseID(mux.Vars(r)["id"])
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

func Pagination(r *http.Request) (Page, error) {
	p := Page{Limit: 50, Query: strings.TrimSpace(r.URL.Query().Get("q"))}
	var err error
	if v := r.URL.Query().Get("limit"); v != "" {
		p.Limit, err = strconv.Atoi(v)
		if err != nil || p.Limit < 1 || p.Limit > 200 {
			return p, Bad("PAGINATION", "limit должен быть от 1 до 200")
		}
	}
	if v := r.URL.Query().Get("offset"); v != "" {
		p.Offset, err = strconv.Atoi(v)
		if err != nil || p.Offset < 0 {
			return p, Bad("PAGINATION", "offset должен быть неотрицательным")
		}
	}
	if v := r.URL.Query().Get("include_deleted"); v != "" {
		p.IncludeDeleted, err = strconv.ParseBool(v)
		if err != nil {
			return p, Bad("PAGINATION", "include_deleted должен быть true или false")
		}
	}
	return p, nil
}

func List(w http.ResponseWriter, items any, total int, p Page) {
	JSON(w, 200, map[string]any{"items": items, "total": total, "limit": p.Limit, "offset": p.Offset})
}
