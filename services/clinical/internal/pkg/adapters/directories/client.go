package directories

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/models"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/pkg/platform"
)

type Client struct {
	staffURL, catalogURL string
	http                 *http.Client
}

func New(staffURL, catalogURL string) (*Client, error) {
	for _, base := range []string{staffURL, catalogURL} {
		u, err := url.Parse(base)
		if err != nil || u.Host == "" || (u.Scheme != "http" && u.Scheme != "https") {
			return nil, fmt.Errorf("invalid service URL %q", base)
		}
	}
	return &Client{staffURL: strings.TrimRight(staffURL, "/"), catalogURL: strings.TrimRight(catalogURL, "/"), http: &http.Client{Timeout: 3 * time.Second}}, nil
}

func (c *Client) get(ctx context.Context, endpoint string, v any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}
	res, err := c.http.Do(req)
	if err != nil {
		return &platform.Error{Status: 503, Code: "DEPENDENCY_UNAVAILABLE", Message: "Справочный сервис недоступен"}
	}
	defer res.Body.Close()
	if res.StatusCode == 404 {
		return platform.Bad("REFERENCE_NOT_FOUND", "Медработник или диагноз отсутствует в справочнике")
	}
	if res.StatusCode != 200 {
		return &platform.Error{Status: 503, Code: "DEPENDENCY_UNAVAILABLE", Message: "Справочный сервис вернул ошибку"}
	}
	if err = json.NewDecoder(io.LimitReader(res.Body, 1<<20)).Decode(v); err != nil {
		return &platform.Error{Status: 503, Code: "DEPENDENCY_RESPONSE", Message: "Некорректный ответ справочного сервиса"}
	}
	return nil
}

func (c *Client) GetActor(ctx context.Context, id string) (models.Actor, error) {
	var v struct {
		models.Actor
		DeletedAt *time.Time `json:"deleted_at"`
	}
	if err := c.get(ctx, c.staffURL+"/api/v1/staff/"+url.PathEscape(id), &v); err != nil {
		return models.Actor{}, err
	}
	if v.DeletedAt != nil {
		return models.Actor{}, platform.Bad("STAFF_DELETED", "Удалённый медработник не может изменять ЭМК")
	}
	if v.ID != id || v.FirstName == "" || v.LastName == "" || v.Position == "" {
		return models.Actor{}, &platform.Error{Status: 503, Code: "DEPENDENCY_RESPONSE", Message: "Неверные данные медработника"}
	}
	return v.Actor, nil
}

func (c *Client) GetDiagnosis(ctx context.Context, id string) (models.DiagnosisSnapshot, error) {
	var v struct {
		models.DiagnosisSnapshot
		DeletedAt *time.Time `json:"deleted_at"`
	}
	if err := c.get(ctx, c.catalogURL+"/api/v1/diagnoses/"+url.PathEscape(id), &v); err != nil {
		return models.DiagnosisSnapshot{}, err
	}
	if v.DeletedAt != nil {
		return models.DiagnosisSnapshot{}, platform.Bad("DIAGNOSIS_DELETED", "Удалённый диагноз нельзя добавить в ЭМК")
	}
	if v.ID != id || v.Code == "" || v.Name == "" {
		return models.DiagnosisSnapshot{}, &platform.Error{Status: 503, Code: "DEPENDENCY_RESPONSE", Message: "Неверные данные диагноза"}
	}
	return v.DiagnosisSnapshot, nil
}
