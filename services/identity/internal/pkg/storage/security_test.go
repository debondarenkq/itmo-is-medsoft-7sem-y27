package storage

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/internal/models"
	"github.com/google/uuid"
)

func securityAccount(t *testing.T, r *Repository) (models.Staff, models.Account) {
	t.Helper()
	ctx := context.Background()
	staff := models.Staff{ID: uuid.NewString(), FirstName: "А", LastName: "Б", Position: "Врач", Role: "doctor"}
	a, err := r.CreateAccount(ctx, models.Account{ID: uuid.NewString(), Username: "doctor", StaffID: &staff.ID, PasswordHash: "test-hash"}, "system:test")
	if err != nil {
		t.Fatal(err)
	}
	return staff, a
}
func TestRefreshRotationAndRevocation(t *testing.T) {
	r := testRepository(t)
	ctx := context.Background()
	_, a := securityAccount(t, r)
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
	updated, err := r.UpdateAccount(ctx, a.ID, models.AccountChanges{RevokeSessions: true}, "system:test")
	if err != nil {
		t.Fatal(err)
	}
	if updated.SessionVersion != a.SessionVersion+1 {
		t.Fatal("account session version did not change")
	}
	if err = r.Issue(ctx, a, "", "stale-account", -1); !errors.Is(err, models.ErrUnauthorized) {
		t.Fatalf("stale account issued a token: %v", err)
	}
	var count int
	if err := r.pool.QueryRow(ctx, "SELECT count(*) FROM refresh_tokens WHERE account_id=$1", a.ID).Scan(&count); err != nil || count != 0 {
		t.Fatalf("revocation: %d %v", count, err)
	}
	events, _, err := r.AuditEvents(ctx, 200, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(events) < 4 {
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
