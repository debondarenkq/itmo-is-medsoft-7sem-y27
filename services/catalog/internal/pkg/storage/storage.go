package storage

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/catalog/internal/models"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/catalog/internal/pkg/platform"
)

type Repository struct{ pool *pgxpool.Pool }

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

const columns = "id, code, name, created_at, updated_at, deleted_at"

func scan(row pgx.Row) (models.Diagnosis, error) {
	var v models.Diagnosis
	err := row.Scan(&v.ID, &v.Code, &v.Name, &v.CreatedAt, &v.UpdatedAt, &v.DeletedAt)
	return v, err
}

func (r *Repository) Get(ctx context.Context, id string) (models.Diagnosis, error) {
	return scan(r.pool.QueryRow(ctx, "SELECT "+columns+" FROM diagnoses WHERE id=$1", id))
}

func (r *Repository) List(ctx context.Context, p platform.Page) ([]models.Diagnosis, int, error) {
	filter := " FROM diagnoses WHERE ($1 OR deleted_at IS NULL) AND ($2='' OR (code || ' ' || name) ILIKE '%' || $2 || '%')"
	var total int
	if err := r.pool.QueryRow(ctx, "SELECT count(*)"+filter, p.IncludeDeleted, p.Query).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := r.pool.Query(ctx, "SELECT "+columns+filter+" ORDER BY created_at, id LIMIT $3 OFFSET $4", p.IncludeDeleted, p.Query, p.Limit, p.Offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	items := []models.Diagnosis{}
	for rows.Next() {
		v, err := scan(rows)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, v)
	}
	return items, total, rows.Err()
}

func (r *Repository) Create(ctx context.Context, v models.Diagnosis) (models.Diagnosis, error) {
	return scan(r.pool.QueryRow(ctx, "INSERT INTO diagnoses(id,code,name) VALUES ($1, $2, $3) RETURNING "+columns, uuid.NewString(), v.Code, v.Name))
}

func (r *Repository) Update(ctx context.Context, id string, v models.Diagnosis) (models.Diagnosis, error) {
	return scan(r.pool.QueryRow(ctx, "UPDATE diagnoses SET code=$2, name=$3, updated_at=clock_timestamp() WHERE id=$1 AND deleted_at IS NULL RETURNING "+columns, id, v.Code, v.Name))
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, "UPDATE diagnoses SET deleted_at=COALESCE(deleted_at,clock_timestamp()), updated_at=clock_timestamp() WHERE id=$1", id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
