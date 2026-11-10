package staff

import (
	"context"
	"fmt"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/internal/models"
	"github.com/google/uuid"
	"net/http"
	"time"
)

const staffRequestTimeout = 2 * time.Second

type Directory struct{ client *ClientWithResponses }

func New(base, key string) (*Directory, error) {
	client, err := NewClientWithResponses(base, WithHTTPClient(&http.Client{Timeout: staffRequestTimeout}), WithRequestEditorFn(func(_ context.Context, r *http.Request) error { r.Header.Set("X-Service-Key", key); return nil }))
	if err != nil {
		return nil, err
	}
	return &Directory{client: client}, nil
}
func profile(value *Profile, expected uuid.UUID) (models.Staff, error) {
	if value == nil || value.ID != expected || value.RoleRevision < 0 || !value.Role.Valid() {
		return models.Staff{}, models.ErrUnavailable
	}
	s := models.Staff{ID: value.ID.String(), FirstName: value.FirstName, LastName: value.LastName, Role: string(value.Role), DeletedAt: value.DeletedAt, RoleRevision: value.RoleRevision}
	if value.Position != nil {
		s.Position = *value.Position
	}
	if value.CreatedAt != nil {
		s.CreatedAt = *value.CreatedAt
	}
	if value.UpdatedAt != nil {
		s.UpdatedAt = *value.UpdatedAt
	}
	return s, nil
}
func (d *Directory) Get(ctx context.Context, id string) (models.Staff, error) {
	identifier, err := uuid.Parse(id)
	if err != nil {
		return models.Staff{}, models.ErrInvalidInput
	}
	response, err := d.client.GetStaffIdentityWithResponse(ctx, identifier)
	if err != nil {
		return models.Staff{}, models.ErrUnavailable
	}
	if response.JSON404 != nil && response.JSON404.Error.Code == NOTFOUND {
		return models.Staff{}, models.ErrNotFound
	}
	if response.JSON200 == nil {
		return models.Staff{}, fmt.Errorf("%w: HTTP %d", models.ErrUnavailable, response.StatusCode())
	}
	return profile(response.JSON200, identifier)
}
func (d *Directory) UpdateRole(ctx context.Context, id, role string) (models.Staff, error) {
	identifier, err := uuid.Parse(id)
	if err != nil {
		return models.Staff{}, models.ErrInvalidInput
	}
	response, err := d.client.UpdateStaffIdentityRoleWithResponse(ctx, identifier, UpdateStaffIdentityRoleJSONRequestBody{Role: Role(role)})
	if err != nil {
		return models.Staff{}, models.ErrUnavailable
	}
	if response.JSON404 != nil && response.JSON404.Error.Code == NOTFOUND {
		return models.Staff{}, models.ErrNotFound
	}
	if response.JSON400 != nil && response.JSON400.Error.Code == INVALIDREQUEST {
		return models.Staff{}, models.ErrInvalidInput
	}
	if response.JSON409 != nil && response.JSON409.Error.Code == CONFLICT {
		return models.Staff{}, models.ErrConflict
	}
	if response.JSON200 == nil {
		return models.Staff{}, fmt.Errorf("%w: HTTP %d", models.ErrUnavailable, response.StatusCode())
	}
	return profile(response.JSON200, identifier)
}

func (d *Directory) All(ctx context.Context) ([]models.Staff, error) {
	response, err := d.client.ListStaffIdentitiesWithResponse(ctx)
	if err != nil {
		return nil, models.ErrUnavailable
	}
	if response.JSON200 == nil {
		return nil, fmt.Errorf("%w: HTTP %d", models.ErrUnavailable, response.StatusCode())
	}
	items := make([]models.Staff, 0, len(*response.JSON200))
	for _, value := range *response.JSON200 {
		p := Profile(value)
		item, err := profile(&p, value.ID)
		if err != nil || value.ID == uuid.Nil {
			return nil, models.ErrUnavailable
		}
		items = append(items, item)
	}
	return items, nil
}
