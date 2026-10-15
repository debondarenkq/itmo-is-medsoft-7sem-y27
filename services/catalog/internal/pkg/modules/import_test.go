package modules

import (
	"context"
	"errors"
	"testing"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/catalog/internal/models"
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
	}{
		{name: "empty"},
		{name: "missing code", data: []models.Diagnosis{{Name: "Диагноз"}}},
		{name: "seven character code", data: []models.Diagnosis{{Code: "E11.900", Name: "Диагноз"}}},
		{name: "invalid second row", data: []models.Diagnosis{{Code: "TEST01", Name: "Первый диагноз"}, {Code: "TEST02", Name: "Я"}}},
		{name: "duplicate trimmed code", data: []models.Diagnosis{{Code: "TEST01", Name: "Первый диагноз"}, {Code: " TEST01 ", Name: "Второй диагноз"}}},
		{name: "duplicate name", data: []models.Diagnosis{{Code: "TEST01", Name: "Первый диагноз"}, {Code: "TEST02", Name: "Первый диагноз"}}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			storage := &importStorage{}
			module := New(Deps{Storage: storage})
			_, err := module.Import(context.Background(), test.data)
			if !errors.Is(err, models.ErrInvalidInput) {
				t.Fatalf("expected invalid input, got %v", err)
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
