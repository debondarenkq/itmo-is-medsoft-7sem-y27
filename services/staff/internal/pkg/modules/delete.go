package modules

import (
	"context"
)

func (module *Module) Delete(ctx context.Context, id string) error {
	return module.Storage.Delete(ctx, id)
}
