package staff

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/internal/models"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	Base, Key string
	HTTP      *http.Client
}

func New(base, key string) *Client {
	return &Client{Base: strings.TrimRight(base, "/"), Key: key, HTTP: &http.Client{Timeout: 2 * time.Second}}
}
func (c *Client) call(ctx context.Context, id string, role *string) (models.Staff, error) {
	var s models.Staff
	method, path, body := "GET", "/internal/staff/"+id, []byte(nil)
	if role != nil {
		method = "PUT"
		path += "/role"
		body, _ = json.Marshal(map[string]string{"role": *role})
	}
	req, err := http.NewRequestWithContext(ctx, method, c.Base+path, bytes.NewReader(body))
	if err != nil {
		return s, models.ErrUnavailable
	}
	req.Header.Set("X-Service-Key", c.Key)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return s, models.ErrUnavailable
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case 404:
		return s, models.ErrNotFound
	case 409:
		return s, models.ErrConflict
	case 200:
	default:
		return s, fmt.Errorf("%w: HTTP %d", models.ErrUnavailable, resp.StatusCode)
	}
	if json.NewDecoder(http.MaxBytesReader(nil, resp.Body, 65536)).Decode(&s) != nil || s.ID != id || s.Revision < 1 {
		return s, models.ErrUnavailable
	}
	return s, nil
}
func (c *Client) Get(ctx context.Context, id string) (models.Staff, error) {
	return c.call(ctx, id, nil)
}
func (c *Client) UpdateRole(ctx context.Context, id, role string) (models.Staff, error) {
	return c.call(ctx, id, &role)
}
