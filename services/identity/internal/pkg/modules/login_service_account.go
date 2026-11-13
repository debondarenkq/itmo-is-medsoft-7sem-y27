package modules

import (
	"context"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/internal/models"
)

func (s *Security) LoginServiceAccount(ctx context.Context, key string) (models.Session, error) {
	if len(key) != serviceAPIKeyLength {
		return models.Session{}, models.ErrUnauthorized
	}
	a, err := s.Store.ServiceAccountByKey(ctx, TokenHash(key))
	if err != nil || !a.IsActive {
		return models.Session{}, models.ErrUnauthorized
	}
	session, err := s.sign(a)
	if err != nil {
		return session, err
	}
	err = s.Store.Issue(ctx, a, "", TokenHash(session.RefreshToken), -1)
	return session, err
}
