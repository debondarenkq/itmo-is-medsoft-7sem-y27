package modules

import (
	"context"
	"strings"
	"unicode/utf8"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/internal/models"
	"github.com/google/uuid"
)

const serviceAPIKeyLength = 64

func (s *Security) SaveServiceAccount(ctx context.Context, id, name string, active, rotate bool, actor string) (models.Account, string, error) {
	name = strings.TrimSpace(name)
	const minNameLength, maxNameLength = 3, 16
	if utf8.RuneCountInString(name) < minNameLength || utf8.RuneCountInString(name) > maxNameLength || strings.ContainsAny(name, "\r\n\t") || name == "admin" {
		return models.Account{}, "", models.ErrInvalidInput
	}
	create := id == ""
	if create {
		id = uuid.NewString()
	}
	a := models.Account{ID: id, Username: name, IsActive: active, IsService: true}
	key := ""
	if create || rotate {
		var err error
		key, err = randomText(serviceAPIKeyLength)
		if err != nil {
			return a, "", err
		}
		a.PasswordHash = TokenHash(key)
	}
	stored, err := s.Store.SaveServiceAccount(ctx, a, create, actor)
	return stored, key, err
}
