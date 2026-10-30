package api

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/internal/models"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/internal/pkg/modules"
)

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (api *SecurityAPI) protect(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if !strings.HasPrefix(path, "/api/") {
			next.ServeHTTP(w, r)
			return
		}
		w.Header().Set("Cache-Control", "no-store")
		if path == "/api/v1/auth/login" || path == "/api/v1/auth/setup" || path == "/api/v1/auth/refresh" {
			w.Header().Set("Cache-Control", "no-store")
			next.ServeHTTP(w, r)
			return
		}
		claims, err := api.claims(r)
		if err != nil {
			api.failed(w, r, err, "Неавторизованный запрос "+r.Method+" "+path)
			return
		}
		if !modules.Allowed(claims, r.Method, path) {
			if err := api.Module.Audit(r.Context(), "account:"+claims.Subject, "Отказ в доступе "+r.Method+" "+path); err != nil {
				securityError(w, err)
				return
			}
			securityError(w, models.ErrForbidden)
			return
		}
		w.Header().Set("Cache-Control", "no-store")
		if r.Method == "GET" {
			next.ServeHTTP(w, r)
			return
		}
		// Record intent before a write, so an unavailable audit store prevents unaudited mutations.
		actor := "account:" + claims.Subject
		if err := api.Module.Audit(r.Context(), actor, "Запрос "+r.Method+" "+path); err != nil {
			securityError(w, err)
			return
		}
		recorder := &statusWriter{ResponseWriter: w, status: 200}
		next.ServeHTTP(recorder, r)
		auditContext, cancel := context.WithTimeout(context.WithoutCancel(r.Context()), 2*time.Second)
		defer cancel()
		if err := api.Module.Audit(auditContext, actor, fmt.Sprintf("Результат %s %s: HTTP %d", r.Method, path, recorder.status)); err != nil {
			slog.Error("audit outcome delivery failed", "path", path, "error", err)
		}
	})
}
