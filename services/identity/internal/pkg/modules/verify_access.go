package modules

import (
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/internal/models"
)

// VerifyAccess validates sessions.
func (s *Security) VerifyAccess(ctx context.Context, token string) (models.Claims, error) {
	claims, err := s.Verify(token)
	if err != nil {
		return claims, err
	}
	a, err := s.Store.AccountByID(ctx, claims.Subject)
	if err != nil {
		return claims, err
	}
	if !a.IsActive || a.SessionVersion != claims.SessionVersion {
		return claims, models.ErrUnauthorized
	}
	if err = s.staffAccount(ctx, &a); err != nil {
		return claims, loginStaffError(err)
	}
	staffID := ""
	if a.StaffID != nil {
		staffID = *a.StaffID
	}
	if staffID != claims.StaffID || a.Role != claims.Role || a.StaffRevision != claims.StaffRevision {
		return claims, models.ErrUnauthorized
	}
	return claims, nil
}

func (security *Security) Verify(token string) (models.Claims, error) {
	var claims models.Claims
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return claims, models.ErrUnauthorized
	}
	enc := base64.RawURLEncoding
	header, err := enc.DecodeString(parts[0])
	if err != nil || string(header) != `{"alg":"EdDSA","typ":"JWT"}` {
		return claims, models.ErrUnauthorized
	}
	signature, err := enc.DecodeString(parts[2])
	if err != nil || !ed25519.Verify(security.key.Public().(ed25519.PublicKey), []byte(parts[0]+"."+parts[1]), signature) {
		return claims, models.ErrUnauthorized
	}
	data, err := enc.DecodeString(parts[1])
	if err != nil || json.Unmarshal(data, &claims) != nil {
		return claims, models.ErrUnauthorized
	}
	if claims.Expires <= time.Now().Unix() || claims.Issuer != tokenIssuer || claims.Audience != "his-api" || claims.Subject == "" {
		return claims, models.ErrUnauthorized
	}
	return claims, nil
}
