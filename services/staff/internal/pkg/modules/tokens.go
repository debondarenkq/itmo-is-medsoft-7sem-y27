package modules

import (
	"crypto/ed25519"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base32"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/internal/models"
	"github.com/google/uuid"
)

func randomText(length int) (string, error) {
	const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, length)
	for i := range result {
		for {
			var b [1]byte
			if _, err := rand.Read(b[:]); err != nil {
				return "", err
			}
			if b[0] < 248 {
				result[i] = alphabet[int(b[0])%len(alphabet)]
				break
			}
		}
	}
	return string(result), nil
}

func TokenHash(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])
}
func NewTOTPSecret() (string, error) {
	b := make([]byte, 20)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(b), nil
}

// RFC 6238: SHA-1, six decimal digits, 30-second steps, one step of clock skew.
func TOTPStep(secret, code string, now time.Time) (int64, bool) {
	if len(code) != 6 {
		return 0, false
	}
	if _, err := strconv.Atoi(code); err != nil {
		return 0, false
	}
	key, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(secret)
	if err != nil {
		return 0, false
	}
	for delta := int64(-1); delta <= 1; delta++ {
		step := now.Unix()/30 + delta
		var counter [8]byte
		binary.BigEndian.PutUint64(counter[:], uint64(step))
		mac := hmac.New(sha1.New, key)
		mac.Write(counter[:])
		digest := mac.Sum(nil)
		offset := digest[len(digest)-1] & 15
		number := (binary.BigEndian.Uint32(digest[offset:offset+4]) & 0x7fffffff) % 1000000
		expected := fmt.Sprintf("%06d", number)
		if subtle.ConstantTimeCompare([]byte(expected), []byte(code)) == 1 {
			return step, true
		}
	}
	return 0, false
}

func (security *Security) sign(account models.Account) (models.Session, error) {
	now := time.Now().UTC()
	expires := now.Add(5 * time.Minute)
	claims := models.Claims{Subject: account.ID, Role: account.Role, Issuer: "his-staff", Audience: "his-api", Expires: expires.Unix(), Issued: now.Unix(), ID: uuid.NewString()}
	if account.StaffID != nil {
		claims.StaffID = *account.StaffID
	}
	data, err := json.Marshal(claims)
	if err != nil {
		return models.Session{}, err
	}
	enc := base64.RawURLEncoding
	payload := enc.EncodeToString([]byte(`{"alg":"EdDSA","typ":"JWT"}`)) + "." + enc.EncodeToString(data)
	token := payload + "." + enc.EncodeToString(ed25519.Sign(security.key, []byte(payload)))
	refresh, err := randomText(64)
	if err != nil {
		return models.Session{}, err
	}
	return models.Session{AccessToken: token, RefreshToken: refresh, ExpiresAt: expires, AccountID: account.ID, StaffID: account.StaffID, Role: account.Role}, nil
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
	if claims.Expires <= time.Now().Unix() || claims.Issuer != "his-staff" || claims.Audience != "his-api" || claims.Subject == "" {
		return claims, models.ErrUnauthorized
	}
	return claims, nil
}

func Allowed(claims models.Claims, method, path string) bool {
	if claims.Role == "admin" {
		return strings.HasPrefix(path, "/api/v1/staff") || strings.HasPrefix(path, "/api/v1/diagnoses") || strings.HasPrefix(path, "/api/v1/accounts") || path == "/api/v1/audit" || path == "/api/v1/auth/logout"
	}
	if claims.StaffID == "" {
		return false
	}
	if path == "/api/v1/auth/logout" {
		return true
	}
	if method == "GET" {
		return strings.HasPrefix(path, "/api/v1/patients") || strings.HasPrefix(path, "/api/v1/records") || strings.HasPrefix(path, "/api/v1/staff/") || strings.HasPrefix(path, "/api/v1/diagnoses")
	}
	if strings.HasPrefix(path, "/api/v1/records/") {
		return claims.Role == "doctor"
	}
	if strings.HasPrefix(path, "/api/v1/patients") {
		return claims.Role == "doctor" || claims.Role == "nurse"
	}
	return false
}
