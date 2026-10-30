package identity

import (
	"bytes"
	"context"
	"fmt"
	"io"
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
func (c *Client) Publish(ctx context.Context, payload []byte) error {
	req, err := http.NewRequestWithContext(ctx, "POST", c.Base+"/internal/staff", bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Service-Key", c.Key)
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)
	if resp.StatusCode != 200 {
		return fmt.Errorf("identity synchronization HTTP %d", resp.StatusCode)
	}
	return nil
}
