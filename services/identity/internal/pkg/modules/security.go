package modules

import (
	"context"
	"crypto/ed25519"
	"fmt"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/internal/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type StaffDirectory interface {
	Get(context.Context, string) (models.Staff, error)
}

type SecurityStorage interface {
	ApplyStaff(context.Context, models.Staff) error
	AccountByID(context.Context, string) (models.Account, error)
	Account(context.Context, string) (models.Account, error)
	RefreshAccount(context.Context, string) (models.Account, error)
	SetupSecret(context.Context, string) (bool, string, error)
	Setup(context.Context, models.Account, string, string, int64) error
	Issue(context.Context, models.Account, string, string, int64) error
	CreateAccount(context.Context, models.Account, string) (models.Account, error)
	ChangeAccount(context.Context, string, *bool, string, string) (models.Account, error)
	Logout(context.Context, string, string) error
	AppendAudit(context.Context, string, string, string) error
	SearchAccounts(context.Context, models.AccountFilter) ([]models.Account, int, error)
	SearchAuditEvents(context.Context, models.AuditFilter) ([]models.AuditEvent, int, error)
	UpdateAccount(context.Context, string, models.AccountChanges, string) (models.Account, error)
	Accounts(context.Context, int, int) ([]models.Account, int, error)
	AuditEvents(context.Context, int, int) ([]models.AuditEvent, int, error)
}
type Security struct {
	Directory StaffDirectory
	Store     SecurityStorage
	key       ed25519.PrivateKey
	dummyHash []byte
}

func NewSecurity(store SecurityStorage, key ed25519.PrivateKey) *Security {
	dummy, _ := bcrypt.GenerateFromPassword([]byte("invalid-login-password"), 12)
	return &Security{Store: store, key: key, dummyHash: dummy}
}
func (s *Security) SetupStatus(ctx context.Context) (bool, string, error) {
	secret, err := NewTOTPSecret()
	if err != nil {
		return false, "", err
	}
	return s.Store.SetupSecret(ctx, secret)
}
func (s *Security) Setup(ctx context.Context, password, confirmation, secret, code string) (models.Session, error) {
	if err := ValidateAdminPassword(password, confirmation); err != nil {
		return models.Session{}, err
	}
	step, ok := TOTPStep(secret, code, time.Now())
	if !ok {
		return models.Session{}, models.ErrUnauthorized
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return models.Session{}, err
	}
	a := models.Account{ID: uuid.NewString(), Username: "admin", PasswordHash: string(hash), IsAdmin: true, IsActive: true, Role: "admin"}
	session, err := s.sign(a)
	if err != nil {
		return session, err
	}
	err = s.Store.Setup(ctx, a, secret, TokenHash(session.RefreshToken), step)
	return session, err
}
func (s *Security) Login(ctx context.Context, username, password, code string) (models.Session, error) {
	a, err := s.Store.Account(ctx, username)
	if err != nil {
		bcrypt.CompareHashAndPassword(s.dummyHash, []byte(password))
		return models.Session{}, models.ErrUnauthorized
	}
	if bcrypt.CompareHashAndPassword([]byte(a.PasswordHash), []byte(password)) != nil || !a.IsActive {
		return models.Session{}, models.ErrUnauthorized
	}
	if !a.IsAdmin {
		checkedPasswordHash := a.PasswordHash
		if err := s.syncStaff(ctx, a.StaffID); err != nil {
			if err == models.ErrConflict || err == models.ErrNotFound {
				return models.Session{}, models.ErrUnauthorized
			}
			return models.Session{}, err
		}
		a, err = s.Store.Account(ctx, username)
		if err != nil || !a.IsActive || a.PasswordHash != checkedPasswordHash {
			return models.Session{}, models.ErrUnauthorized
		}
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
func (s *Security) Refresh(ctx context.Context, token string) (models.Session, error) {
	if len(token) != 64 {
		return models.Session{}, models.ErrUnauthorized
	}
	a, err := s.Store.RefreshAccount(ctx, TokenHash(token))
	if err != nil {
		return models.Session{}, err
	}
	if !a.IsAdmin {
		if err := s.syncStaff(ctx, a.StaffID); err != nil {
			if err == models.ErrConflict || err == models.ErrNotFound {
				return models.Session{}, models.ErrUnauthorized
			}
			return models.Session{}, err
		}
		a, err = s.Store.RefreshAccount(ctx, TokenHash(token))
		if err != nil {
			return models.Session{}, err
		}
	}
	if !a.IsActive {
		return models.Session{}, models.ErrUnauthorized
	}
	session, err := s.sign(a)
	if err != nil {
		return session, err
	}
	err = s.Store.Issue(ctx, a, TokenHash(token), TokenHash(session.RefreshToken), -1)
	return session, err
}
func (s *Security) CreateAccount(ctx context.Context, username, staffID, actor string) (models.Account, string, error) {
	username = strings.TrimSpace(username)
	if n := utf8.RuneCountInString(username); n < 3 || n > 16 || strings.ContainsAny(username, "\r\n\t") || username == "admin" {
		return models.Account{}, "", models.ErrInvalidInput
	}
	if _, err := uuid.Parse(staffID); err != nil {
		return models.Account{}, "", models.ErrInvalidInput
	}
	if err := s.syncStaff(ctx, &staffID); err != nil {
		return models.Account{}, "", err
	}
	password, err := randomText(12)
	if err != nil {
		return models.Account{}, "", err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return models.Account{}, "", err
	}
	a, err := s.Store.CreateAccount(ctx, models.Account{ID: uuid.NewString(), Username: username, StaffID: &staffID, PasswordHash: string(hash)}, actor)
	return a, password, err
}
func (s *Security) ResetPassword(ctx context.Context, id, actor string) (models.Account, string, error) {
	password, err := randomText(12)
	if err != nil {
		return models.Account{}, "", err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return models.Account{}, "", err
	}
	a, err := s.Store.ChangeAccount(ctx, id, nil, string(hash), actor)
	return a, password, err
}

func ValidateAdminPassword(password, confirmation string) error {
	if password != confirmation {
		return fmt.Errorf("%w: Пароли должны совпадать", models.ErrInvalidInput)
	}
	upper, lower, digit := false, false, false
	for _, character := range password {
		upper = upper || (unicode.IsUpper(character))
		lower = lower || (unicode.IsLower(character))
		digit = digit || (unicode.IsNumber(character))
	}
	if utf8.RuneCountInString(password) < 8 || len(password) > 72 || !upper || !lower || !digit {
		return fmt.Errorf("%w: Пароль должен содержать минимум 8 символов, большие и маленькие буквы и цифру; максимум 72 байта", models.ErrInvalidInput)
	}
	return nil
}

func (s *Security) UpdateAccount(ctx context.Context, id string, changes models.AccountChanges, actor string) (models.Account, error) {
	if changes.Username != nil {
		username := strings.TrimSpace(*changes.Username)
		if n := utf8.RuneCountInString(username); n < 3 || n > 16 || strings.ContainsAny(username, "\r\n\t") || username == "admin" {
			return models.Account{}, models.ErrInvalidInput
		}
		changes.Username = &username
	}
	if changes.Role != nil {
		switch *changes.Role {
		case "doctor", "nurse", "laboratory-technician", "other":
		default:
			return models.Account{}, models.ErrInvalidInput
		}
	}
	a, err := s.Store.AccountByID(ctx, id)
	if err != nil {
		return a, err
	}
	if a.IsAdmin {
		return a, models.ErrForbidden
	}
	if err := s.syncStaff(ctx, a.StaffID); err != nil {
		return a, err
	}
	return s.Store.UpdateAccount(ctx, id, changes, actor)
}

func (s *Security) syncStaff(ctx context.Context, id *string) error {
	if id == nil {
		return models.ErrUnauthorized
	}
	if s.Directory == nil {
		return models.ErrUnavailable
	}
	profile, err := s.Directory.Get(ctx, *id)
	if err != nil {
		return err
	}
	if err = s.Store.ApplyStaff(ctx, profile); err != nil {
		return err
	}
	if profile.DeletedAt != nil {
		return models.ErrConflict
	}
	return nil
}
