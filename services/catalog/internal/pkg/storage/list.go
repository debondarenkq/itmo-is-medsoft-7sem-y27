package storage

import (
	"context"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/catalog/internal/models"
)

func (repository *Repository) List(ctx context.Context, page models.Page) ([]models.Diagnosis, int, error) {
	filter := " FROM diagnoses WHERE ($1 OR deleted_at IS NULL) AND ($2='' OR (code || ' ' || name) ILIKE '%' || $2 || '%')"
	var total int
	if err := repository.pool.QueryRow(ctx, "SELECT count(*)"+filter, page.IncludeDeleted, page.Search).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := repository.pool.Query(ctx, "SELECT "+columns+filter+" ORDER BY created_at, id LIMIT $3 OFFSET $4", page.IncludeDeleted, page.Search, page.Limit, page.Offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	items := []models.Diagnosis{}
	for rows.Next() {
		diagnosis, err := scan(rows)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, diagnosis)
	}
	return items, total, rows.Err()
}
