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
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/pkg/platform"
)

type Client struct {
	staff   *staffapi.ClientWithResponses
	catalog *catalogapi.ClientWithResponses
}

func New(staffURL, catalogURL string) (*Client, error) {
	for _, base := range []string{staffURL, catalogURL} {
		u, err := url.Parse(base)
		if err != nil || u.Host == "" || (u.Scheme != "http" && u.Scheme != "https") {
			return nil, fmt.Errorf("invalid service URL %q", base)
		}
	}
	client := &http.Client{Timeout: 3 * time.Second}
	staff, err := staffapi.NewClientWithResponses(strings.TrimRight(staffURL, "/"), staffapi.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}
	catalog, err := catalogapi.NewClientWithResponses(strings.TrimRight(catalogURL, "/"), catalogapi.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}
	return &Client{staff: staff, catalog: catalog}, nil
}

func dependencyError() error {
	return &platform.Error{Status: 503, Code: "DEPENDENCY_UNAVAILABLE", Message: "Справочный сервис недоступен или вернул некорректный ответ"}
}

func (c *Client) GetActor(ctx context.Context, id string) (models.Actor, error) {
	entityID, err := uuid.Parse(id)
	if err != nil {
		return models.Actor{}, platform.Bad("INVALID_ID", "Требуется UUID медработника")
	}
	response, err := c.staff.GetStaffWithResponse(ctx, entityID)
	if err != nil {
		return models.Actor{}, dependencyError()
	}
	if response.StatusCode() == http.StatusNotFound {
		return models.Actor{}, platform.Bad("REFERENCE_NOT_FOUND", "Медработник отсутствует в справочнике")
	}
	if response.JSON200 == nil {
		return models.Actor{}, dependencyError()
	}
	staff := response.JSON200
	if staff.DeletedAt != nil {
		return models.Actor{}, platform.Bad("STAFF_DELETED", "Удалённый медработник не может изменять ЭМК")
	}
	if staff.ID != entityID || staff.FirstName == "" || staff.LastName == "" || staff.Position == "" {
		return models.Actor{}, dependencyError()
	}
	return models.Actor{ID: staff.ID.String(), FirstName: staff.FirstName, LastName: staff.LastName, Position: staff.Position}, nil
}

func (c *Client) GetDiagnosis(ctx context.Context, id string) (models.DiagnosisSnapshot, error) {
	entityID, err := uuid.Parse(id)
	if err != nil {
		return models.DiagnosisSnapshot{}, platform.Bad("INVALID_ID", "Требуется UUID диагноза")
	}
	response, err := c.catalog.GetDiagnosisWithResponse(ctx, entityID)
	if err != nil {
		return models.DiagnosisSnapshot{}, dependencyError()
	}
	if response.StatusCode() == http.StatusNotFound {
		return models.DiagnosisSnapshot{}, platform.Bad("REFERENCE_NOT_FOUND", "Диагноз отсутствует в справочнике")
	}
	if response.JSON200 == nil {
		return models.DiagnosisSnapshot{}, dependencyError()
	}
	diagnosis := response.JSON200
	if diagnosis.DeletedAt != nil {
		return models.DiagnosisSnapshot{}, platform.Bad("DIAGNOSIS_DELETED", "Удалённый диагноз нельзя добавить в ЭМК")
	}
	if diagnosis.ID != entityID || diagnosis.Code == "" || diagnosis.Name == "" {
		return models.DiagnosisSnapshot{}, dependencyError()
	}
	return models.DiagnosisSnapshot{ID: diagnosis.ID.String(), Code: diagnosis.Code, Name: diagnosis.Name}, nil
}
