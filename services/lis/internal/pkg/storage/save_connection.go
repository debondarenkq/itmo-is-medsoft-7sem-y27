package storage

import (
	"context"
	"encoding/json"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/lis/internal/models"
)

func (r *Repository) SaveConnection(ctx context.Context, c models.Connection) error {
	data, err := json.Marshal(c.Session)
	if err != nil {
		return err
	}
	_, err = r.Pool.Exec(ctx, "INSERT INTO connection(id,session) VALUES(true,$1) ON CONFLICT(id) DO UPDATE SET session=EXCLUDED.session", data)
	return err
}
