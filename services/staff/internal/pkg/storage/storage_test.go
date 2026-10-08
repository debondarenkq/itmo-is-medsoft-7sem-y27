package storage

import (
	"context"
	"testing"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/internal/models"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/internal/pkg/platform"
)

func TestStaffSoftDelete(t *testing.T) {
	r := testRepository(t)
	ctx := context.Background()
	v, err := r.Create(ctx, models.Staff{FirstName: "Анна", LastName: "Петрова", Position: "Терапевт"})
	if err != nil {
		t.Fatal(err)
	}
	if err = r.Delete(ctx, v.ID); err != nil {
		t.Fatal(err)
	}
	deleted, err := r.Get(ctx, v.ID)
	if err != nil || deleted.DeletedAt == nil {
		t.Fatalf("soft delete: %v", err)
	}
	active, total, err := r.List(ctx, platform.Page{Limit: 50})
	if err != nil || total != 0 || len(active) != 0 {
		t.Fatal("deleted staff in active list")
	}
	all, total, err := r.List(ctx, platform.Page{Limit: 50, IncludeDeleted: true})
	if err != nil || total != 1 || len(all) != 1 {
		t.Fatal("staff lost from database")
	}
}
