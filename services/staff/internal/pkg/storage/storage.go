package storage

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/internal/models"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/internal/pkg/platform"
)

type Repository struct{ pool *pgxpool.Pool }

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

const columns = "id, first_name, last_name, position, created_at, updated_at, deleted_at"

func scan(row pgx.Row) (models.Staff, error) {
	var v models.Staff
	err := row.Scan(&v.ID, &v.FirstName, &v.LastName, &v.Position, &v.CreatedAt, &v.UpdatedAt, &v.DeletedAt)
	return v, err
}

func (r *Repository) Get(ctx context.Context, id string) (models.Staff, error) {
	return scan(r.pool.QueryRow(ctx, "SELECT "+columns+" FROM staff WHERE id=$1", id))
}

func (r *Repository) List(ctx context.Context, p platform.Page) ([]models.Staff, int, error) {
	filter := " FROM staff WHERE ($1 OR deleted_at IS NULL) AND ($2='' OR (first_name || ' ' || last_name || ' ' || position) ILIKE '%' || $2 || '%')"
	var total int
	if err := r.pool.QueryRow(ctx, "SELECT count(*)"+filter, p.IncludeDeleted, p.Query).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := r.pool.Query(ctx, "SELECT "+columns+filter+" ORDER BY created_at, id LIMIT $3 OFFSET $4", p.IncludeDeleted, p.Query, p.Limit, p.Offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	items := []models.Staff{}
	for rows.Next() {
		v, err := scan(rows)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, v)
	}
	return items, total, rows.Err()
}

func (r *Repository) Create(ctx context.Context, v models.Staff) (models.Staff, error) {
	return scan(r.pool.QueryRow(ctx, "INSERT INTO staff(id,first_name,last_name,position) VALUES ($1, $2, $3, $4) RETURNING "+columns, uuid.NewString(), v.FirstName, v.LastName, v.Position))
}

func (r *Repository) Update(ctx context.Context, id string, v models.Staff) (models.Staff, error) {
	return scan(r.pool.QueryRow(ctx, "UPDATE staff SET first_name=$2, last_name=$3, position=$4, updated_at=clock_timestamp() WHERE id=$1 AND deleted_at IS NULL RETURNING "+columns, id, v.FirstName, v.LastName, v.Position))
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, "UPDATE staff SET deleted_at=COALESCE(deleted_at,clock_timestamp()), updated_at=clock_timestamp() WHERE id=$1", id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
