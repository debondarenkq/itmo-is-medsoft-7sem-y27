package modules

import (
	"context"
	"crypto/ed25519"
	"fmt"
	"net/url"
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
	All(context.Context) ([]models.Staff, error)
	UpdateRole(context.Context, string, string) (models.Staff, error)
}

type SecurityStorage interface {
	AccountByID(context.Context, string) (models.Account, error)
	Account(context.Context, string) (models.Account, error)
	RefreshAccount(context.Context, string) (models.Account, error)
	SetupSecret(context.Context, string) (bool, string, error)
	Setup(context.Context, models.Account, string, string, int64) error
	Issue(context.Context, models.Account, string, string, int64) error
	CreateAccount(context.Context, models.Account, string) (models.Account, error)
	ChangeAccount(context.Context, string, *bool, string, string) (models.Account, error)
	Logout(context.Context, string, string) error
	AppendAudit(context.Context, string, string, models.AuditEventKind) error
	SearchAccounts(context.Context, models.AccountFilter) ([]models.Account, int, error)
	SearchAuditEvents(context.Context, models.AuditFilter) ([]models.AuditEvent, int, error)
	UpdateAccount(context.Context, string, models.AccountChanges, string) (models.Account, error)
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
func (s *Security) Refresh(ctx context.Context, token string) (models.Session, error) {
	if len(token) != 64 {
		return models.Session{}, models.ErrUnauthorized
	}
	a, err := s.Store.RefreshAccount(ctx, TokenHash(token))
	if err != nil {
		return models.Session{}, err
	}
	revision := a.StaffRevision
	if err := s.staffAccount(ctx, &a); err != nil {
		return models.Session{}, loginStaffError(err)
	}
	if a.StaffRevision != revision {
		return models.Session{}, models.ErrUnauthorized
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
	a := models.Account{ID: uuid.NewString(), Username: username, StaffID: &staffID, IsActive: true}
	if err := s.staffAccount(ctx, &a); err != nil {
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
	a.PasswordHash = string(hash)
	stored, err := s.Store.CreateAccount(ctx, a, actor)
	stored.Role, stored.StaffName, stored.StaffRevision = a.Role, a.StaffName, a.StaffRevision
	return stored, password, err
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
	a, err := s.Store.AccountByID(ctx, id)
	if err != nil {
		return a, "", err
	}
	if err = s.staffAccount(ctx, &a); err != nil {
		return a, "", err
	}
	stored, err := s.Store.ChangeAccount(ctx, id, nil, string(hash), actor)
	stored.Role, stored.StaffName, stored.StaffRevision = a.Role, a.StaffName, a.StaffRevision
	return stored, password, err
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
	a, err := s.Store.AccountByID(ctx, id)
	if err != nil {
		return a, err
	}
	if a.IsAdmin {
		return a, models.ErrForbidden
	}
	if err := s.staffAccount(ctx, &a); err != nil {
		return a, err
	}
	if changes.Username != nil {
		other, err := s.Store.Account(ctx, *changes.Username)
		if err == nil && other.ID != id {
			return a, models.ErrConflict
		}
		if err != nil && err != models.ErrUnauthorized {
			return a, err
		}
	}
	if changes.Role != nil && *changes.Role != a.Role {
		profile, err := s.Directory.UpdateRole(ctx, *a.StaffID, *changes.Role)
		if err != nil {
			return a, err
		}
		a.Role = profile.Role
		a.StaffRevision = profile.RoleRevision
	}
	changes.RevokeSessions = changes.Role != nil
	changes.Role = nil
	stored, err := s.Store.UpdateAccount(ctx, id, changes, actor)
	stored.Role, stored.StaffName, stored.StaffRevision = a.Role, a.StaffName, a.StaffRevision
	return stored, err
}

func loginStaffError(err error) error {
	if err == models.ErrConflict || err == models.ErrNotFound {
		return models.ErrUnauthorized
	}
	return err
}

func (s *Security) staffAccount(ctx context.Context, a *models.Account) error {
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

func (s *Security) SearchAccounts(ctx context.Context, filter models.AccountFilter) ([]models.Account, int, error) {
	profiles, err := s.Directory.All(ctx)
	if err != nil {
		return nil, 0, err
	}
	byID := make(map[string]models.Staff, len(profiles))
	filter.StaffIDs = []string{}
	filter.DeletedStaffIDs = []string{}
	for _, p := range profiles {
		byID[p.ID] = p
		if strings.Contains(strings.ToLower(staffName(p)), strings.ToLower(filter.Search)) {
			filter.StaffIDs = append(filter.StaffIDs, p.ID)
		}
		if p.DeletedAt != nil {
			filter.DeletedStaffIDs = append(filter.DeletedStaffIDs, p.ID)
		}
	}
	accounts, total, err := s.Store.SearchAccounts(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	for i := range accounts {
		a := &accounts[i]
		if a.IsAdmin {
			continue
		}
		if a.StaffID != nil {
			p, ok := byID[*a.StaffID]
			if ok {
				a.StaffName = staffName(p)
				a.Role = p.Role
				a.IsActive = a.IsActive && p.DeletedAt == nil
			}
		}
	}
	return accounts, total, nil
}
