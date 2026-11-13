package modules

import (
	"context"
	"strings"
	"unicode/utf8"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/internal/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

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
