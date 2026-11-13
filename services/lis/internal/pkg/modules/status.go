package modules

import (
	"context"
	"errors"

	api "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/lis/api"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/lis/internal/models"
)

func (m *Module) Status(ctx context.Context) (api.LISConnectionStatus, error) {
	_, err := m.Store.Connection(ctx)
	if errors.Is(err, models.ErrUnconfigured) {
		return api.LISConnectionStatus{Message: "Введите API-ключ сервисного аккаунта LIS"}, nil
	}
	if err != nil {
		return api.LISConnectionStatus{}, err
	}
	return api.LISConnectionStatus{Configured: true, Message: "Авторизация настроена. Можно получать назначения из HIS."}, nil
}
