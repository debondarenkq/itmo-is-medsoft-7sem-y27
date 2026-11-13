package modules

import (
	"context"
	"fmt"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/internal/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

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
