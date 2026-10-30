package storage

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/internal/models"
	"github.com/google/uuid"
)

func securityAccount(t *testing.T, r *Repository) (models.Staff, models.Account) {
	t.Helper()
	ctx := context.Background()
	staff, err := r.Create(ctx, models.Staff{FirstName: "А", LastName: "Б", Position: "Врач", Role: "doctor"})
	if err != nil {
		t.Fatal(err)
	}
	a, err := r.CreateAccount(ctx, models.Account{ID: uuid.NewString(), Username: "doctor", StaffID: &staff.ID, PasswordHash: "test-hash"}, "system:test")
	if err != nil {
		t.Fatal(err)
	}
	return staff, a
}
func TestRefreshRotationAndRevocation(t *testing.T) {
	r := testRepository(t)
	ctx := context.Background()
	staff, a := securityAccount(t, r)
	if err := r.Issue(ctx, a, "", "first", -1); err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	results := make(chan error, 2)
	for _, next := range []string{"second", "third"} {
		wg.Go(func() { results <- r.Issue(ctx, a, "first", next, -1) })
	}
	wg.Wait()
	close(results)
	successes := 0
	for err := range results {
		if err == nil {
			successes++
		} else if !errors.Is(err, models.ErrUnauthorized) {
			t.Fatal(err)
		}
	}
	if successes != 1 {
		t.Fatalf("concurrent rotation successes=%d", successes)
	}
	if _, err := r.RefreshAccount(ctx, "first"); !errors.Is(err, models.ErrUnauthorized) {
		t.Fatalf("old refresh accepted: %v", err)
	}
	staff.Role = "nurse"
	if _, err := r.Update(ctx, staff.ID, staff); err != nil {
		t.Fatal(err)
	}
	var count int
	if err := r.pool.QueryRow(ctx, "SELECT count(*) FROM refresh_tokens WHERE account_id=$1", a.ID).Scan(&count); err != nil || count != 0 {
		t.Fatalf("role revocation: %d %v", count, err)
	}
	updated, err := r.Account(ctx, a.Username)
	if err != nil {
		t.Fatal(err)
	}
	if err := r.Issue(ctx, updated, "", "after-role-change", -1); err != nil {
		t.Fatal(err)
	}
	if err := r.Delete(ctx, staff.ID); err != nil {
		t.Fatal(err)
	}
	updated, err = r.Account(ctx, a.Username)
	if err != nil || updated.IsActive {
		t.Fatalf("deleted staff account active: %+v %v", updated, err)
	}
	if _, err = r.RefreshAccount(ctx, "after-role-change"); !errors.Is(err, models.ErrUnauthorized) {
		t.Fatalf("deleted staff refresh: %v", err)
	}
	if _, err = r.ChangeAccount(ctx, a.ID, new(true), "", "system:test"); !errors.Is(err, models.ErrConflict) {
		t.Fatalf("reactivated deleted staff account: %v", err)
	}
	events, _, err := r.AuditEvents(ctx, 200, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(events) < 5 {
		t.Fatalf("missing audit events: %d", len(events))
	}
}
func TestAccountPasswordAndActivityRevokeRefresh(t *testing.T) {
	r := testRepository(t)
	ctx := context.Background()
	_, a := securityAccount(t, r)
	if err := r.Issue(ctx, a, "", "password-session", -1); err != nil {
		t.Fatal(err)
	}
	a, err := r.ChangeAccount(ctx, a.ID, nil, "new-hash", "system:test")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = r.RefreshAccount(ctx, "password-session"); !errors.Is(err, models.ErrUnauthorized) {
		t.Fatalf("password refresh: %v", err)
	}
	if err = r.Issue(ctx, a, "", "activity-session", -1); err != nil {
		t.Fatal(err)
	}
	if _, err = r.ChangeAccount(ctx, a.ID, new(false), "", "system:test"); err != nil {
		t.Fatal(err)
	}
	if _, err = r.RefreshAccount(ctx, "activity-session"); !errors.Is(err, models.ErrUnauthorized) {
		t.Fatalf("deactivated refresh: %v", err)
	}
}
