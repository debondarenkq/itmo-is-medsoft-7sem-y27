package modules

import (
	"context"
	"errors"
	"testing"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/catalog/internal/models"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/catalog/internal/pkg/platform"
)

type importStorage struct {
	Storage
	called bool
	items  []models.Diagnosis
	err    error
}

func (s *importStorage) UpsertBatch(_ context.Context, diagnoses []models.Diagnosis) ([]models.Diagnosis, error) {
	s.called = true
	s.items = diagnoses
	return diagnoses, s.err
}

func TestImportValidatesWholeFileBeforeStorage(t *testing.T) {
	tests := []struct {
		name string
		data []models.Diagnosis
		code string
	}{
		{name: "empty", code: "IMPORT_SIZE"},
		{name: "missing code", data: []models.Diagnosis{{Name: "Диагноз"}}, code: "VALIDATION"},
		{name: "seven character code", data: []models.Diagnosis{{Code: "E11.900", Name: "Диагноз"}}, code: "VALIDATION"},
		{name: "invalid second row", data: []models.Diagnosis{{Code: "TEST01", Name: "Первый диагноз"}, {Code: "TEST02", Name: "Я"}}, code: "VALIDATION"},
		{name: "duplicate trimmed code", data: []models.Diagnosis{{Code: "TEST01", Name: "Первый диагноз"}, {Code: " TEST01 ", Name: "Второй диагноз"}}, code: "IMPORT_DUPLICATE_CODE"},
		{name: "duplicate name", data: []models.Diagnosis{{Code: "TEST01", Name: "Первый диагноз"}, {Code: "TEST02", Name: "Первый диагноз"}}, code: "IMPORT_DUPLICATE_NAME"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			storage := &importStorage{}
			module := New(Deps{Storage: storage})
			_, err := module.Import(context.Background(), test.data)
			var api *platform.Error
			if !errors.As(err, &api) || api.Code != test.code {
				t.Fatalf("error = %v, want %s", err, test.code)
			}
			if storage.called {
				t.Fatal("invalid file reached storage")
			}
		})
	}
}

func TestImportPassesOnlySelectedDataAndPropagatesStorageFailure(t *testing.T) {
	expected := errors.New("batch failed")
	storage := &importStorage{err: expected}
	module := New(Deps{Storage: storage})
	_, err := module.Import(context.Background(), []models.Diagnosis{{Code: " TEST01 ", Name: " Первый диагноз "}})
	if !errors.Is(err, expected) {
		t.Fatalf("error = %v", err)
	}
	if !storage.called || len(storage.items) != 1 || storage.items[0].Code != "TEST01" || storage.items[0].Name != "Первый диагноз" {
		t.Fatalf("imported unexpected data: %#v", storage.items)
	}
}
