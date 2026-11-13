package modules

import (
	"context"
	"strings"
	"unicode/utf8"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/internal/models"
)

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
	if a.IsAdmin || a.IsService {
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
