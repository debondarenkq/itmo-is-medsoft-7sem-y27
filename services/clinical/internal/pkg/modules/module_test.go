package modules

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/models"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/pkg/platform"
)

func TestMissingAuthorRejectedBeforeDependencies(t *testing.T) {
	m := New(Deps{}) // Nil dependencies make accidental calls fail the test.
	_, err := m.Save(context.Background(), uuid.NewString(), models.RecordChanges{})
	var api *platform.Error
	if !errors.As(err, &api) || api.Code != "STAFF_REQUIRED" {
		t.Fatalf("missing author: %v", err)
	}
	_, err = m.CreateRecord(context.Background(), uuid.NewString(), "")
	if !errors.As(err, &api) || api.Code != "STAFF_REQUIRED" {
		t.Fatalf("missing creator: %v", err)
	}
}
