package modules

import (
	"context"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/models"
)

func (module *Module) History(ctx context.Context, id string, page models.Page) ([]models.Event, int, error) {
	if _, err := module.Storage.GetRecord(ctx, id); err != nil {
		return nil, 0, err
	}
	return module.Storage.History(ctx, id, page)
}
