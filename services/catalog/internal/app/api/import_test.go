package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/catalog/internal/models"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/catalog/internal/pkg/modules"
)

type importRepository struct {
	modules.Storage
	items []models.Diagnosis
}

func (r *importRepository) UpsertBatch(_ context.Context, diagnoses []models.Diagnosis) ([]models.Diagnosis, error) {
	r.items = diagnoses
	for index := range diagnoses {
		diagnoses[index].ID = uuid.NewString()
	}
	return diagnoses, nil
}

func TestImportFileContents(t *testing.T) {
	tests := []struct {
		name        string
		contentType string
		body        string
		status      int
	}{
		{name: "JSON file", contentType: "application/json", body: `{"diagnoses":[{"code":"TEST01","name":"Первый диагноз"}]}`, status: http.StatusOK},
		{name: "empty file", contentType: "application/json", body: `{"diagnoses":[]}`, status: http.StatusBadRequest},
		{name: "malformed JSON", contentType: "application/json", body: `{"diagnoses":`, status: http.StatusBadRequest},
		{name: "unknown field", contentType: "application/json", body: `{"diagnoses":[{"code":"TEST01","name":"Первый диагноз","id":"client-id"}]}`, status: http.StatusBadRequest},
		{name: "unsupported file format", contentType: "text/csv", body: "TEST01,Первый диагноз", status: http.StatusUnsupportedMediaType},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repository := &importRepository{}
			module := modules.New(modules.Deps{Storage: repository})
			router := mux.NewRouter()
			New(Deps{Module: module}).Register(router)
			request := httptest.NewRequest(http.MethodPost, "/api/v1/diagnoses/import", strings.NewReader(test.body))
			request.Header.Set("Content-Type", test.contentType)
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)
			if response.Code != test.status {
				t.Fatalf("status = %d, response = %s", response.Code, response.Body.String())
			}
			if test.status == http.StatusOK {
				if len(repository.items) != 1 || repository.items[0].Code != "TEST01" {
					t.Fatalf("unexpected data saved: %#v", repository.items)
				}
			} else if len(repository.items) != 0 {
				t.Fatal("invalid file changed the catalog")
			}
		})
	}
}
