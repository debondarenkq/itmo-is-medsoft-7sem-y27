package storage

import (
	"context"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool       *pgxpool.Pool
	UpdateRole func(context.Context, string, string) (models.Staff, error)
}

func New(pool *pgxpool.Pool) *Repository { return &Repository{pool: pool} }

func (r *Repository) ApplyStaff(ctx context.Context, s models.Staff) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO staff_profiles(id,first_name,last_name,role,deleted_at,revision,security_revision)
 VALUES($1,$2,$3,$4,$5,$6,$7) ON CONFLICT(id) DO UPDATE SET
 first_name=EXCLUDED.first_name,last_name=EXCLUDED.last_name,role=EXCLUDED.role,
 deleted_at=EXCLUDED.deleted_at,revision=EXCLUDED.revision,security_revision=EXCLUDED.security_revision
 WHERE staff_profiles.revision<EXCLUDED.revision`, s.ID, s.FirstName, s.LastName, s.Role, s.DeletedAt, s.Revision, s.SecurityRevision)
	return err
}
