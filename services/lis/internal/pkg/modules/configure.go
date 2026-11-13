package modules

import (
	"context"

	api "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/lis/api"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/lis/internal/models"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/lis/internal/pkg/adapters"
	identity "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/lis/internal/pkg/adapters/identityapi"
)

func (m *Module) Configure(ctx context.Context, input api.LISConnectionInput) (api.LISConnectionStatus, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	client, err := adapters.Identity(m.Gateway, m.HTTP)
	if err != nil {
		return api.LISConnectionStatus{}, err
	}
	response, err := client.LoginServiceAccountWithResponse(ctx, identity.ServiceLoginInput{APIKey: input.APIKey})
	if err != nil {
		return api.LISConnectionStatus{}, err
	}
	if response.JSON200 == nil {
		return api.LISConnectionStatus{}, adapters.ResponseError(response.StatusCode(), response.Body)
	}
	c := models.Connection{Session: *response.JSON200}
	status, err := m.check(ctx, c)
	if err != nil {
		return status, err
	}
	return status, m.Store.SaveConnection(ctx, c)
}

func (m *Module) check(ctx context.Context, c models.Connection) (api.LISConnectionStatus, error) {
	client, err := adapters.Laboratory(m.Gateway, c.Session.AccessToken, m.HTTP)
	if err != nil {
		return api.LISConnectionStatus{}, err
	}
	response, err := client.TestLabConnectionWithResponse(ctx)
	if err != nil {
		return api.LISConnectionStatus{}, err
	}
	if response.JSON200 == nil {
		return api.LISConnectionStatus{}, adapters.ResponseError(response.StatusCode(), response.Body)
	}
	return api.LISConnectionStatus{Configured: true, Message: "Доступ к HIS подтверждён. Роль LIS."}, nil
}
