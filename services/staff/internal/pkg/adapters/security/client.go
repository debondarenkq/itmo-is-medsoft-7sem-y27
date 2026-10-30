package security

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type contextKey int

const tokenKey contextKey = 0
const staffKey contextKey = 1

func Authorization(ctx context.Context) string {
	value, _ := ctx.Value(tokenKey).(string)
	return value
}
func StaffID(ctx context.Context) string { value, _ := ctx.Value(staffKey).(string); return value }

type Client struct {
	base, key string
	http      *http.Client
}

func New() (*Client, error) {
	base := os.Getenv("AUTH_URL")
	u, err := url.Parse(base)
	if err != nil || u.Host == "" || (u.Scheme != "http" && u.Scheme != "https") {
		return nil, fmt.Errorf("AUTH_URL is required")
	}
	key, err := os.ReadFile(os.Getenv("SERVICE_KEY_FILE"))
	if err != nil {
		return nil, err
	}
	if len(strings.TrimSpace(string(key))) < 32 {
		return nil, fmt.Errorf("invalid service key")
	}
	return &Client{base: strings.TrimRight(base, "/"), key: strings.TrimSpace(string(key)), http: &http.Client{Timeout: 2 * time.Second}}, nil
}
func (client *Client) call(ctx context.Context, path, token string, input any) ([]byte, int, error) {
	body, err := json.Marshal(input)
	if err != nil {
		return nil, 0, err
	}
	request, err := http.NewRequestWithContext(ctx, "POST", client.base+path, bytes.NewReader(body))
	if err != nil {
		return nil, 0, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Service-Key", client.key)
	request.Header.Set("Authorization", token)
	response, err := client.http.Do(request)
	if err != nil {
		return nil, 0, err
	}
	defer response.Body.Close()
	data, err := io.ReadAll(io.LimitReader(response.Body, 65536))
	return data, response.StatusCode, err
}
func reject(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]any{"error": map[string]string{"code": "ACCESS_DENIED", "message": message}})
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}
func (client *Client) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/api/") {
			next.ServeHTTP(w, r)
			return
		}
		w.Header().Set("Cache-Control", "no-store")
		token := r.Header.Get("Authorization")
		data, status, err := client.call(r.Context(), "/internal/authorize", token, map[string]string{"method": r.Method, "path": r.URL.Path})
		if err != nil || status >= 500 {
			reject(w, 503, "Сервис авторизации недоступен")
			return
		}
		if status != 200 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(status)
			w.Write(data)
			return
		}
		var claims struct {
			Subject string `json:"sub"`
			StaffID string `json:"staff_id"`
		}
		if json.Unmarshal(data, &claims) != nil || claims.Subject == "" {
			reject(w, 503, "Некорректный ответ авторизации")
			return
		}
		ctx := context.WithValue(r.Context(), tokenKey, token)
		ctx = context.WithValue(ctx, staffKey, claims.StaffID)
		r = r.WithContext(ctx)
		w.Header().Set("Cache-Control", "no-store")
		if r.Method == "GET" {
			next.ServeHTTP(w, r)
			return
		}
		actor := "account:" + claims.Subject
		_, status, err = client.call(ctx, "/internal/audit", token, map[string]string{"actor": actor, "message": "Запрос " + r.Method + " " + r.URL.Path})
		if err != nil || status != 200 {
			reject(w, 503, "Сервис аудита недоступен")
			return
		}
		recorder := &statusWriter{ResponseWriter: w, status: 200}
		next.ServeHTTP(recorder, r)
		auditContext, cancel := context.WithTimeout(context.WithoutCancel(ctx), 2*time.Second)
		defer cancel()
		_, status, err = client.call(auditContext, "/internal/audit", token, map[string]string{"actor": actor, "message": fmt.Sprintf("Результат %s %s: HTTP %d", r.Method, r.URL.Path, recorder.status)})
		if err != nil || status != 200 {
			slog.Error("audit outcome delivery failed", "path", r.URL.Path, "error", err, "status", status)
		}
	})
}
