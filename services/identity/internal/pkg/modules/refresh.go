package modules

import (
	"context"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/internal/models"
)

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
