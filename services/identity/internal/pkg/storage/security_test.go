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

func TestServiceAccountRevocation(t *testing.T) {
	for _, test := range []struct {
		name   string
		active bool
		hash   string
	}{{"deactivate", false, ""}, {"rotate key", true, "new-api-hash"}} {
		t.Run(test.name, func(t *testing.T) {
			r := testRepository(t)
			ctx := context.Background()
			a, err := r.SaveServiceAccount(ctx, models.Account{ID: uuid.NewString(), Username: "lab", PasswordHash: "api-hash", IsService: true, IsActive: true}, true, "system:test")
			if err != nil {
				t.Fatal(err)
			}
			if a.Role != "lis" || !a.IsService {
				t.Fatalf("service role: %+v", a)
			}
			if err = r.Issue(ctx, a, "", "service-refresh", -1); err != nil {
				t.Fatal(err)
			}
			before := a
			a.IsActive = test.active
			a.PasswordHash = test.hash
			updated, err := r.SaveServiceAccount(ctx, a, false, "system:test")
			if err != nil {
				t.Fatal(err)
			}
			if _, err = r.RefreshAccount(ctx, "service-refresh"); !errors.Is(err, models.ErrUnauthorized) {
				t.Fatalf("refresh survived: %v", err)
			}
			if err = r.Issue(ctx, before, "", "late-refresh", -1); !errors.Is(err, models.ErrUnauthorized) {
				t.Fatalf("stale login survived: %v", err)
			}
			if updated.SessionVersion != before.SessionVersion+1 {
				t.Fatal("access tokens were not revoked")
			}
			accounts, total, err := r.SearchAccounts(ctx, models.AccountFilter{Limit: 50, StaffIDs: []string{}, DeletedStaffIDs: []string{}})
			if err != nil || total != 0 || len(accounts) != 0 {
				t.Fatalf("mixed account lists: %d %v", total, err)
			}
			if _, err = r.Account(ctx, "lab"); !errors.Is(err, models.ErrUnauthorized) {
				t.Fatalf("service account accepted as user: %v", err)
			}
		})
	}
}
