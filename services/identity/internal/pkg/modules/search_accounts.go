package modules

import (
	"context"
	"strings"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/internal/models"
)

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
