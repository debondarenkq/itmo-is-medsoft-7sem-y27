package modules

import (
	"context"
	"crypto/rand"
	"encoding/base32"
	"net/url"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/internal/models"
)

func (s *Security) SetupStatus(ctx context.Context) (models.AdminSetup, error) {
	candidate, err := NewTOTPSecret()
	if err != nil {
		return models.AdminSetup{}, err
	}
	configured, secret, err := s.Store.SetupSecret(ctx, candidate)
	if err != nil {
		return models.AdminSetup{}, err
	}
	if configured {
		return models.AdminSetup{Configured: true}, nil
	}
	query := url.Values{"secret": {secret}, "issuer": {"HIS"}, "algorithm": {"SHA1"}, "digits": {"6"}, "period": {"30"}}
	uri := (&url.URL{Scheme: "otpauth", Host: "totp", Path: "/HIS:admin", RawQuery: query.Encode()}).String()
	return models.AdminSetup{Secret: &secret, ProvisioningURI: &uri}, nil
}

func NewTOTPSecret() (string, error) {
	b := make([]byte, 20)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(b), nil
}
