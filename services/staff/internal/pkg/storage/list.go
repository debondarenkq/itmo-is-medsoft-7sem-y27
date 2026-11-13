package storage

import (
	"context"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/internal/models"
)

func (repository *Repository) List(ctx context.Context, page models.Page) ([]models.Staff, int, error) {
	filter := " FROM staff WHERE ($1 OR deleted_at IS NULL) AND ($2='' OR (first_name || ' ' || last_name || ' ' || position) ILIKE '%' || $2 || '%')"
	var total int
	if err := repository.pool.QueryRow(ctx, "SELECT count(*)"+filter, page.IncludeDeleted, page.Search).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := repository.pool.Query(ctx, "SELECT "+columns+filter+" ORDER BY created_at, id LIMIT $3 OFFSET $4", page.IncludeDeleted, page.Search, page.Limit, page.Offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	items := []models.Staff{}
	for rows.Next() {
		staff, err := scan(rows)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, staff)
	}
	return items, total, rows.Err()
}
