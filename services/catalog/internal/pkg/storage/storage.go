package storage

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/catalog/internal/models"
)

type Repository struct{ pool *pgxpool.Pool }

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

const columns = "id, code, name, created_at, updated_at, deleted_at"

func scan(row pgx.Row) (models.Diagnosis, error) {
	var diagnosis models.Diagnosis
	err := row.Scan(&diagnosis.ID, &diagnosis.Code, &diagnosis.Name, &diagnosis.CreatedAt, &diagnosis.UpdatedAt, &diagnosis.DeletedAt)
	return diagnosis, mapDiagnosisError(err)
}

func (repository *Repository) Get(ctx context.Context, id string) (models.Diagnosis, error) {
	return scan(repository.pool.QueryRow(ctx, "SELECT "+columns+" FROM diagnoses WHERE id=$1", id))
}

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

func (repository *Repository) Create(ctx context.Context, diagnosis models.Diagnosis) (models.Diagnosis, error) {
	return scan(repository.pool.QueryRow(ctx, "INSERT INTO diagnoses(id,code,name) VALUES ($1, $2, $3) RETURNING "+columns, uuid.NewString(), diagnosis.Code, diagnosis.Name))
}

func (repository *Repository) Update(ctx context.Context, id string, diagnosis models.Diagnosis) (models.Diagnosis, error) {
	return scan(repository.pool.QueryRow(ctx, "UPDATE diagnoses SET code=$2, name=$3, updated_at=clock_timestamp() WHERE id=$1 AND deleted_at IS NULL RETURNING "+columns, id, diagnosis.Code, diagnosis.Name))
}

func (repository *Repository) Delete(ctx context.Context, id string) error {
	tag, err := repository.pool.Exec(ctx, "UPDATE diagnoses SET deleted_at=COALESCE(deleted_at,clock_timestamp()), updated_at=clock_timestamp() WHERE id=$1", id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return models.ErrNotFound
	}
	return nil
}
