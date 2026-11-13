package modules

import (
	"context"
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

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/internal/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const tokenIssuer = "his-identity"

func (s *Security) Login(ctx context.Context, username, password, code string) (models.Session, error) {
	a, err := s.Store.Account(ctx, username)
	if err != nil {
		bcrypt.CompareHashAndPassword(s.dummyHash, []byte(password))
		return models.Session{}, models.ErrUnauthorized
	}
	if bcrypt.CompareHashAndPassword([]byte(a.PasswordHash), []byte(password)) != nil || !a.IsActive {
		return models.Session{}, models.ErrUnauthorized
	}
	if err := s.staffAccount(ctx, &a); err != nil {
		return models.Session{}, loginStaffError(err)
	}
	step := int64(-1)
	if a.IsAdmin {
		var ok bool
		step, ok = TOTPStep(a.TOTPSecret, code, time.Now())
		if !ok {
			return models.Session{}, models.ErrUnauthorized
		}
	}
	session, err := s.sign(a)
	if err != nil {
		return session, err
	}
	err = s.Store.Issue(ctx, a, "", TokenHash(session.RefreshToken), step)
	return session, err
}

func loginStaffError(err error) error {
	if err == models.ErrConflict || err == models.ErrNotFound {
		return models.ErrUnauthorized
	}
	return err
}

func (s *Security) staffAccount(ctx context.Context, a *models.Account) error {
	if a.IsService {
		a.Role = models.ServiceRole
		return nil
	}
	if a.IsAdmin {
		a.Role = "admin"
		return nil
	}
	if a.StaffID == nil {
		return models.ErrUnauthorized
	}
	if s.Directory == nil {
		return models.ErrUnavailable
	}
	profile, err := s.Directory.Get(ctx, *a.StaffID)
	if err != nil {
		return err
	}
	if profile.DeletedAt != nil {
		return models.ErrConflict
	}
	a.Role = profile.Role
	a.StaffName = staffName(profile)
	a.StaffRevision = profile.RoleRevision
	return nil
}

func staffName(s models.Staff) string { return strings.TrimSpace(s.LastName + " " + s.FirstName) }

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

// TOTPStep validates codes.
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
	claims := models.Claims{SessionVersion: account.SessionVersion, StaffRevision: account.StaffRevision, Subject: account.ID, Role: account.Role, Issuer: tokenIssuer, Audience: "his-api", Expires: expires.Unix(), Issued: now.Unix(), ID: uuid.NewString()}
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
