package directories

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/models"
)

func TestDeletedStaffCannotAuthorChanges(t *testing.T) {
	id := uuid.NewString()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path != "/api/v1/staff/"+id {
			t.Errorf("path: %s", r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"id": id, "first_name": "Анна", "last_name": "Петрова", "position": "Терапевт", "deleted_at": time.Now().UTC()})
	}))
	defer server.Close()
	c, err := New(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.GetActor(context.Background(), id)
	if !errors.Is(err, models.ErrInvalidInput) {
		t.Fatalf("deleted actor accepted: %v", err)
	}
}

func TestUnavailableDirectoryReturnsDependencyError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(500) }))
	defer server.Close()
	c, err := New(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.GetDiagnosis(context.Background(), uuid.NewString())
	if !errors.Is(err, models.ErrUnavailable) {
		t.Fatalf("dependency failure: %v", err)
	}
}
