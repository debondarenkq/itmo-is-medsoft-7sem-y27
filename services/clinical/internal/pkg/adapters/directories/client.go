package directories

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/models"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/pkg/adapters/catalogapi"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/pkg/adapters/staffapi"
)

type Client struct {
	staff   *staffapi.ClientWithResponses
	catalog *catalogapi.ClientWithResponses
}

func New(gatewayURL string) (*Client, error) {
	serviceURL, err := url.Parse(gatewayURL)
	if err != nil || serviceURL.Host == "" || (serviceURL.Scheme != "http" && serviceURL.Scheme != "https") {
		return nil, fmt.Errorf("invalid GATEWAY_URL %q", gatewayURL)
	}
	gatewayURL = strings.TrimRight(gatewayURL, "/")

	httpClient := &http.Client{Timeout: 3 * time.Second}
	staff, err := staffapi.NewClientWithResponses(gatewayURL, staffapi.WithHTTPClient(httpClient))
	if err != nil {
		return nil, err
	}
	catalog, err := catalogapi.NewClientWithResponses(gatewayURL, catalogapi.WithHTTPClient(httpClient))
	if err != nil {
		return nil, err
	}
	return &Client{staff: staff, catalog: catalog}, nil
}

func (client *Client) GetActor(ctx context.Context, id string) (models.Actor, error) {
	entityID, err := uuid.Parse(id)
	if err != nil {
		return models.Actor{}, fmt.Errorf("%w: Требуется UUID медработника", models.ErrInvalidInput)
	}
	response, err := client.staff.GetStaffWithResponse(ctx, entityID)
	if err != nil {
		return models.Actor{}, models.ErrUnavailable
	}
	if response.StatusCode() == http.StatusNotFound {
		return models.Actor{}, fmt.Errorf("%w: Медработник отсутствует в справочнике", models.ErrInvalidInput)
	}
	if response.JSON200 == nil {
		return models.Actor{}, models.ErrUnavailable
	}
	staff := response.JSON200
	if staff.DeletedAt != nil {
		return models.Actor{}, fmt.Errorf("%w: Удалённый медработник не может изменять ЭМК", models.ErrInvalidInput)
	}
	if staff.ID != entityID || staff.FirstName == "" || staff.LastName == "" || staff.Position == "" {
		return models.Actor{}, models.ErrUnavailable
	}
	return models.Actor{ID: staff.ID.String(), FirstName: staff.FirstName, LastName: staff.LastName, Position: staff.Position}, nil
}

func (client *Client) GetDiagnosis(ctx context.Context, id string) (models.DiagnosisSnapshot, error) {
	entityID, err := uuid.Parse(id)
	if err != nil {
		return models.DiagnosisSnapshot{}, fmt.Errorf("%w: Требуется UUID диагноза", models.ErrInvalidInput)
	}
	response, err := client.catalog.GetDiagnosisWithResponse(ctx, entityID)
	if err != nil {
		return models.DiagnosisSnapshot{}, models.ErrUnavailable
	}
	if response.StatusCode() == http.StatusNotFound {
		return models.DiagnosisSnapshot{}, fmt.Errorf("%w: Диагноз отсутствует в справочнике", models.ErrInvalidInput)
	}
	if response.JSON200 == nil {
		return models.DiagnosisSnapshot{}, models.ErrUnavailable
	}
	diagnosis := response.JSON200
	if diagnosis.DeletedAt != nil {
		return models.DiagnosisSnapshot{}, fmt.Errorf("%w: Удалённый диагноз нельзя добавить в ЭМК", models.ErrInvalidInput)
	}
	if diagnosis.ID != entityID || diagnosis.Code == "" || diagnosis.Name == "" {
		return models.DiagnosisSnapshot{}, models.ErrUnavailable
	}
	return models.DiagnosisSnapshot{ID: diagnosis.ID.String(), Code: diagnosis.Code, Name: diagnosis.Name}, nil
}
