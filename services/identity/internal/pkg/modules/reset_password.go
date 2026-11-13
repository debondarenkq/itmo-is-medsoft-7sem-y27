package modules

import (
	"context"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/internal/models"
	"golang.org/x/crypto/bcrypt"
)

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
	if a.IsService {
		return a, "", models.ErrForbidden
	}
	if err = s.staffAccount(ctx, &a); err != nil {
		return a, "", err
	}
	stored, err := s.Store.ChangeAccount(ctx, id, nil, string(hash), actor)
	stored.Role, stored.StaffName, stored.StaffRevision = a.Role, a.StaffName, a.StaffRevision
	return stored, password, err
}
