package storage

import (
	"context"
	"testing"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/catalog/internal/models"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/catalog/internal/pkg/platform"
)

func TestSeedIsIdempotentAndRestoresDeletedEntries(t *testing.T) {
	r := testRepository(t)
	ctx := context.Background()
	first, err := r.Seed(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if err = r.Delete(ctx, first[0].ID); err != nil {
		t.Fatal(err)
	}
	second, err := r.Seed(ctx)
	if err != nil {
		t.Fatal(err)
	}
	for j := range first {
		if first[j].ID != second[j].ID || second[j].DeletedAt != nil {
			t.Fatal("seed duplicated or failed to restore")
		}
	}
	_, total, err := r.List(ctx, platform.Page{Limit: 50})
	if err != nil || total != 5 {
		t.Fatalf("seed count=%d err=%v", total, err)
	}
	if _, err = r.Create(ctx, models.Diagnosis{Code: "J00.00", Name: "Другая запись"}); err == nil {
		t.Fatal("duplicate code accepted")
	}
	if _, err = r.Create(ctx, models.Diagnosis{Code: "Z99.00", Name: first[0].Name}); err == nil {
		t.Fatal("duplicate name accepted")
	}
}
