package storage

import (
	"context"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/laboratory/internal/models"
)

func (r *Repository) SaveTest(ctx context.Context, t models.Test, create bool) (models.Test, error) {
	var err error
	if create {
		_, err = r.pool.Exec(ctx, "INSERT INTO lab_tests(id,code,name) VALUES($1,$2,$3)", t.ID, t.Code, t.Name)
	} else {
		err = r.pool.QueryRow(ctx, "UPDATE lab_tests SET code=$2,name=$3 WHERE id=$1 RETURNING active", t.ID, t.Code, t.Name).Scan(&t.Active)
	}
	return t, databaseError(err)
}
