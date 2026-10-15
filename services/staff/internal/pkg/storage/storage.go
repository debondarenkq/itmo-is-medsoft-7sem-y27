package storage

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/internal/models"
)

type Repository struct{ pool *pgxpool.Pool }

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

const columns = "id, first_name, last_name, position, created_at, updated_at, deleted_at"

func scan(row pgx.Row) (models.Staff, error) {
	var staff models.Staff
	err := row.Scan(&staff.ID, &staff.FirstName, &staff.LastName, &staff.Position, &staff.CreatedAt, &staff.UpdatedAt, &staff.DeletedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return staff, models.ErrNotFound
	}
	return staff, err
}

func (repository *Repository) Get(ctx context.Context, id string) (models.Staff, error) {
	return scan(repository.pool.QueryRow(ctx, "SELECT "+columns+" FROM staff WHERE id=$1", id))
}

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

func (repository *Repository) Create(ctx context.Context, staff models.Staff) (models.Staff, error) {
	return scan(repository.pool.QueryRow(ctx, "INSERT INTO staff(id,first_name,last_name,position) VALUES ($1, $2, $3, $4) RETURNING "+columns, uuid.NewString(), staff.FirstName, staff.LastName, staff.Position))
}

func (repository *Repository) Update(ctx context.Context, id string, staff models.Staff) (models.Staff, error) {
	return scan(repository.pool.QueryRow(ctx, "UPDATE staff SET first_name=$2, last_name=$3, position=$4, updated_at=clock_timestamp() WHERE id=$1 AND deleted_at IS NULL RETURNING "+columns, id, staff.FirstName, staff.LastName, staff.Position))
}

func (repository *Repository) Delete(ctx context.Context, id string) error {
	tag, err := repository.pool.Exec(ctx, "UPDATE staff SET deleted_at=COALESCE(deleted_at,clock_timestamp()), updated_at=clock_timestamp() WHERE id=$1", id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return models.ErrNotFound
	}
	return nil
}
