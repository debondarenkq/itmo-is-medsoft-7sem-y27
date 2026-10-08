package storage

import (
	"context"
	"testing"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/catalog/internal/models"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/catalog/internal/pkg/platform"
)

func TestUpsertBatchIsIdempotentAndRestoresDeletedEntries(t *testing.T) {
	r := testRepository(t)
	ctx := context.Background()
	diagnoses := []models.Diagnosis{
		{Code: "TEST01", Name: "Первый диагноз"},
		{Code: "TEST02", Name: "Второй диагноз"},
	}
	first, err := r.UpsertBatch(ctx, diagnoses)
	if err != nil {
		t.Fatal(err)
	}
	if err = r.Delete(ctx, first[0].ID); err != nil {
		t.Fatal(err)
	}
	second, err := r.UpsertBatch(ctx, diagnoses)
	if err != nil {
		t.Fatal(err)
	}
	for j := range first {
		if first[j].ID != second[j].ID || second[j].DeletedAt != nil {
			t.Fatal("batch duplicated or failed to restore")
		}
	}
	_, total, err := r.List(ctx, platform.Page{Limit: 50})
	if err != nil || total != len(diagnoses) {
		t.Fatalf("batch count=%d err=%v", total, err)
	}
	if _, err = r.Create(ctx, models.Diagnosis{Code: diagnoses[0].Code, Name: "Другая запись"}); err == nil {
		t.Fatal("duplicate code accepted")
	}
	if _, err = r.Create(ctx, models.Diagnosis{Code: "TEST03", Name: first[0].Name}); err == nil {
		t.Fatal("duplicate name accepted")
	}
}

func TestUpsertBatchRollsBackOnConflict(t *testing.T) {
	r := testRepository(t)
	ctx := context.Background()
	initial, err := r.UpsertBatch(ctx, []models.Diagnosis{
		{Code: "TEST01", Name: "Первый диагноз"},
		{Code: "TEST02", Name: "Второй диагноз"},
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = r.UpsertBatch(ctx, []models.Diagnosis{
		{Code: "TEST01", Name: "Изменённый диагноз"},
		{Code: "TEST03", Name: initial[1].Name},
	})
	if err == nil {
		t.Fatal("conflicting batch accepted")
	}
	first, err := r.Get(ctx, initial[0].ID)
	if err != nil || first.Name != initial[0].Name {
		t.Fatalf("failed batch changed existing diagnosis: %v", err)
	}
	_, total, err := r.List(ctx, platform.Page{Limit: 50})
	if err != nil || total != 2 {
		t.Fatalf("failed batch changed row count: %d, %v", total, err)
	}
}
