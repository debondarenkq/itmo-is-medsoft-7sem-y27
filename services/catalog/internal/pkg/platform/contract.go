package platform

import (
	"context"
	"mime"
	"net/http"
	"strings"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/google/uuid"
	nethttpmiddleware "github.com/oapi-codegen/nethttp-middleware"
)

func RequireID(id uuid.UUID) error {
	if id == uuid.Nil {
		return Bad("INVALID_ID", "Требуется ненулевой UUID")
	}
	return nil
}

func PageFromParams(limit, offset *int, query *string, includeDeleted *bool) Page {
	page := Page{Limit: 50}
	if limit != nil {
		page.Limit = *limit
	}
	if offset != nil {
		page.Offset = *offset
	}
	if query != nil {
		page.Query = strings.TrimSpace(*query)
	}
	if includeDeleted != nil {
		page.IncludeDeleted = *includeDeleted
	}
	return page
}

func WriteRequestError(w http.ResponseWriter, r *http.Request, err error) {
	WriteError(w, r, Bad("INVALID_REQUEST", err.Error()))
}

// Validate and bind transport data before invoking a typed business handler.
func ContractMiddleware(spec *openapi3.T) func(http.Handler) http.Handler {
	validator := nethttpmiddleware.OapiRequestValidatorWithOptions(spec, &nethttpmiddleware.Options{
		DoNotValidateServers: true,
		ErrorHandlerWithOpts: func(_ context.Context, err error, w http.ResponseWriter, r *http.Request, options nethttpmiddleware.ErrorHandlerOpts) {
			WriteError(w, r, &Error{Status: options.StatusCode, Code: "INVALID_REQUEST", Message: err.Error()})
		},
	})
	return func(next http.Handler) http.Handler {
		validated := validator(next)
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !strings.HasPrefix(r.URL.Path, "/api/v1/") {
				next.ServeHTTP(w, r)
				return
			}
			ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
			defer cancel()
			r = r.WithContext(ctx)
			r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
			if r.Method == http.MethodPost || r.Method == http.MethodPut {
				media, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
				if err != nil || media != "application/json" {
					WriteError(w, r, &Error{Status: 415, Code: "CONTENT_TYPE", Message: "Требуется Content-Type: application/json"})
					return
				}
			}
			validated.ServeHTTP(w, r)
		})
	}
}
