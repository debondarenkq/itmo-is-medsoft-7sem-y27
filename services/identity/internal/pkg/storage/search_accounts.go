package storage

import (
	"context"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/internal/models"
)

func (r *Repository) SearchAccounts(ctx context.Context, filter models.AccountFilter) ([]models.Account, int, error) {
	where := accountFrom + "WHERE NOT a.is_service AND ($1='' OR a.username ILIKE '%'||$1||'%' OR a.staff_id=ANY($3::uuid[])) AND ($2::boolean IS NULL OR (a.is_active AND NOT COALESCE(a.staff_id=ANY($4::uuid[]),false))=$2)"
	var total int
	if err := r.pool.QueryRow(ctx, "SELECT count(*)"+where, filter.Search, filter.IsActive, filter.StaffIDs, filter.DeletedStaffIDs).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := r.pool.Query(ctx, "SELECT "+accountColumns+where+" ORDER BY a.created_at,a.id LIMIT $5 OFFSET $6", filter.Search, filter.IsActive, filter.StaffIDs, filter.DeletedStaffIDs, filter.Limit, filter.Offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	items := []models.Account{}
	for rows.Next() {
		account, err := scanAccount(rows)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, account)
	}
	return items, total, rows.Err()
}
