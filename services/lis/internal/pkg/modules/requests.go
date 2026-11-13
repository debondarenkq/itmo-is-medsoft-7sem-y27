package modules

import (
	"context"
	"time"

	api "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/lis/api"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/lis/internal/models"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/lis/internal/pkg/adapters"
	identity "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/lis/internal/pkg/adapters/identityapi"
	his "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/lis/internal/pkg/adapters/laboratoryapi"
)

const refreshMargin = 15 * time.Second

func (m *Module) Requests(ctx context.Context) (api.FHIRBundle, error) {
	client, _, err := m.client(ctx)
	if err != nil {
		return api.FHIRBundle{}, err
	}
	response, err := client.SearchFHIRServiceRequestWithResponse(ctx, nil)
	if err != nil {
		return api.FHIRBundle{}, err
	}
	if response.ApplicationFhirJSON200 == nil {
		return api.FHIRBundle{}, adapters.ResponseError(response.StatusCode(), response.Body)
	}
	return adapters.Convert[api.FHIRBundle](response.ApplicationFhirJSON200)
}

func (m *Module) client(ctx context.Context) (*his.ClientWithResponses, models.Connection, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	c, err := m.Store.Connection(ctx)
	if err != nil {
		return nil, c, err
	}
	if time.Until(c.Session.ExpiresAt) <= refreshMargin {
		client, err := adapters.Identity(m.Gateway, m.HTTP)
		if err != nil {
			return nil, c, err
		}
		response, err := client.RefreshSessionWithResponse(ctx, identity.TokenInput{RefreshToken: c.Session.RefreshToken})
		if err != nil {
			return nil, c, err
		}
		if response.JSON200 == nil {
			return nil, c, adapters.ResponseError(response.StatusCode(), response.Body)
		}
		c.Session = *response.JSON200
		if err = m.Store.SaveConnection(ctx, c); err != nil {
			return nil, c, err
		}
	}
	client, err := adapters.Laboratory(m.Gateway, c.Session.AccessToken, m.HTTP)
	return client, c, err
}
