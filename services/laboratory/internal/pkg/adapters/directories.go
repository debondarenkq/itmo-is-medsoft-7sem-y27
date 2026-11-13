package adapters

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"

	api "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/laboratory/api"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/laboratory/internal/models"
	clinical "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/laboratory/internal/pkg/adapters/clinicalapi"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/laboratory/internal/pkg/adapters/security"
)

const directoryTimeout = 5 * time.Second

type Directories struct {
	clinical *clinical.ClientWithResponses
	staff    *security.ClientWithResponses
}

func NewDirectories() (*Directories, error) {
	clinicalURL, staffURL := os.Getenv("CLINICAL_URL"), os.Getenv("STAFF_URL")
	for _, base := range []string{clinicalURL, staffURL} {
		parsed, err := url.Parse(base)
		if err != nil || parsed.Host == "" || (parsed.Scheme != "http" && parsed.Scheme != "https") {
			return nil, fmt.Errorf("CLINICAL_URL and STAFF_URL are required")
		}
	}
	data, err := os.ReadFile(os.Getenv("SERVICE_KEY_FILE"))
	if err != nil {
		return nil, err
	}
	key := strings.TrimSpace(string(data))
	edit := func(_ context.Context, r *http.Request) error { r.Header.Set("X-Service-Key", key); return nil }
	c, err := clinical.NewClientWithResponses(clinicalURL, clinical.WithHTTPClient(&http.Client{Timeout: directoryTimeout}), clinical.WithRequestEditorFn(edit))
	if err != nil {
		return nil, err
	}
	staff, err := security.NewClientWithResponses(staffURL, security.WithHTTPClient(&http.Client{Timeout: directoryTimeout}), security.WithRequestEditorFn(edit))
	return &Directories{c, staff}, err
}
func (d *Directories) Patient(ctx context.Context, id uuid.UUID, requireRecord bool) (api.FHIRResource, error) {
	response, err := d.clinical.GetPatientIdentityWithResponse(ctx, id)
	if err != nil {
		return api.FHIRResource{}, models.ErrUnavailable
	}
	if response.JSON404 != nil {
		return api.FHIRResource{}, models.ErrNotFound
	}
	if response.JSON200 == nil {
		return api.FHIRResource{}, models.ErrUnavailable
	}
	p := response.JSON200
	if requireRecord && !p.HasRecord {
		return api.FHIRResource{}, fmt.Errorf("%w: Сначала создайте ЭМК", models.ErrInvalidInput)
	}
	given := []string{p.FirstName}
	if p.MiddleName != nil {
		given = append(given, *p.MiddleName)
	}
	names := []api.FHIRHumanName{{Family: p.LastName, Given: given}}
	gender := api.Unknown
	if p.AdministrativeSex == clinical.M {
		gender = api.Male
	}
	if p.AdministrativeSex == clinical.F {
		gender = api.Female
	}
	return api.FHIRResource{ID: id, ResourceType: api.Patient, Name: &names, BirthDate: &p.BirthDate, Gender: &gender}, nil
}
func (d *Directories) Practitioner(ctx context.Context, id uuid.UUID) (api.FHIRResource, error) {
	response, err := d.staff.GetStaffIdentityWithResponse(ctx, id)
	if err != nil {
		return api.FHIRResource{}, models.ErrUnavailable
	}
	if response.JSON404 != nil {
		return api.FHIRResource{}, models.ErrNotFound
	}
	if response.JSON200 == nil {
		return api.FHIRResource{}, models.ErrUnavailable
	}
	s := response.JSON200
	names := []api.FHIRHumanName{{Family: s.LastName, Given: []string{s.FirstName}}}
	active := s.DeletedAt == nil
	return api.FHIRResource{ID: id, ResourceType: api.Practitioner, Name: &names, Active: &active}, nil
}
