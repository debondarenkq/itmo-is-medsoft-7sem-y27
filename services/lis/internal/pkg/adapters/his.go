package adapters

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	identity "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/lis/internal/pkg/adapters/identityapi"
	laboratory "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/lis/internal/pkg/adapters/laboratoryapi"
)

const hisTimeout = 15 * time.Second

func GatewayAddress() (string, error) {
	address := strings.TrimRight(strings.TrimSpace(os.Getenv("GATEWAY_URL")), "/")
	parsed, err := url.Parse(address)
	if err != nil || parsed.Host == "" || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.User != nil || parsed.Path != "" || parsed.RawQuery != "" || parsed.Fragment != "" {
		return "", fmt.Errorf("GATEWAY_URL must identify the HIS gateway")
	}
	return address, nil
}

func HTTPClient() (*http.Client, error) {
	roots, err := x509.SystemCertPool()
	if err != nil {
		return nil, err
	}
	if path := os.Getenv("HIS_CA_FILE"); path != "" {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		if !roots.AppendCertsFromPEM(data) {
			return nil, fmt.Errorf("invalid HIS CA")
		}
	}
	transport := &http.Transport{TLSClientConfig: &tls.Config{MinVersion: tls.VersionTLS13, RootCAs: roots}}
	return &http.Client{Timeout: hisTimeout, Transport: transport, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}, nil
}
func Identity(base string, client *http.Client) (*identity.ClientWithResponses, error) {
	return identity.NewClientWithResponses(base, identity.WithHTTPClient(client))
}
func Laboratory(base, token string, client *http.Client) (*laboratory.ClientWithResponses, error) {
	return laboratory.NewClientWithResponses(base, laboratory.WithHTTPClient(client), laboratory.WithRequestEditorFn(func(_ context.Context, r *http.Request) error {
		r.Header.Set("Authorization", "Bearer "+token)
		return nil
	}))
}
func ResponseError(status int, body []byte) error {
	var outcome struct {
		Issue []struct {
			Diagnostics string `json:"diagnostics"`
		} `json:"issue"`
		Error struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	if json.Unmarshal(body, &outcome) == nil {
		if len(outcome.Issue) > 0 {
			return fmt.Errorf("HIS HTTP %d: %s", status, outcome.Issue[0].Diagnostics)
		}
		if outcome.Error.Message != "" {
			return fmt.Errorf("HIS HTTP %d: %s", status, outcome.Error.Message)
		}
	}
	return fmt.Errorf("HIS HTTP %d", status)
}
func Convert[T any](value any) (T, error) {
	var result T
	data, err := json.Marshal(value)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(data, &result)
	return result, err
}
