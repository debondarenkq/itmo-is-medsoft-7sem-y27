package storage

import (
	"context"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/internal/models"
	"github.com/google/uuid"
	"testing"
	"time"
)

func (r *Repository) Create(ctx context.Context, s models.Staff) (models.Staff, error) {
	s.ID = uuid.NewString()
	s.Revision = 1
	return s, r.ApplyStaff(ctx, s)
}
func (r *Repository) Get(ctx context.Context, id string) (models.Staff, error) {
	var s models.Staff
	err := r.pool.QueryRow(ctx, "SELECT id,first_name,last_name,role,deleted_at,revision,security_revision FROM staff_profiles WHERE id=$1", id).Scan(&s.ID, &s.FirstName, &s.LastName, &s.Role, &s.DeletedAt, &s.Revision, &s.SecurityRevision)
	return s, err
}
func (r *Repository) Update(ctx context.Context, id string, s models.Staff) (models.Staff, error) {
	s.ID = id
	s.Revision++
	s.SecurityRevision = s.Revision
	return s, r.ApplyStaff(ctx, s)
}
func (r *Repository) Delete(ctx context.Context, id string) error {
	s, err := r.Get(ctx, id)
	if err != nil {
		return err
	}
	now := time.Now()
	s.DeletedAt = &now
	s.Revision++
	s.SecurityRevision = s.Revision
	return r.ApplyStaff(ctx, s)
}
func TestProfileReplayDoesNotUndoDeletionOrRevokeTwice(t *testing.T) {
	r := testRepository(t)
	ctx := context.Background()
	profile, a := securityAccount(t, r)
	if err := r.Issue(ctx, a, "", "before-delete", -1); err != nil {
		t.Fatal(err)
	}
	if err := r.Delete(ctx, profile.ID); err != nil {
		t.Fatal(err)
	}
	deleted, err := r.Get(ctx, profile.ID)
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range []models.Staff{profile, deleted, profile} {
		if err := r.ApplyStaff(ctx, s); err != nil {
			t.Fatal(err)
		}
	}
	a, err = r.Account(ctx, a.Username)
	if err != nil || a.IsActive {
		t.Fatalf("stale event reactivated account: %v", err)
	}
	var count int
	if err = r.pool.QueryRow(ctx, "SELECT count(*) FROM audit_events WHERE actor=$1", "system:staff:"+profile.ID).Scan(&count); err != nil || count != 1 {
		t.Fatalf("replayed revocation: count=%d err=%v", count, err)
	}
}

func TestLatestProfileRevokesEvenWhenRoleReturnsToOriginal(t *testing.T) {
	r := testRepository(t)
	ctx := context.Background()
	profile, a := securityAccount(t, r)
	if err := r.Issue(ctx, a, "", "before-role-round-trip", -1); err != nil {
		t.Fatal(err)
	}
	// Deliver the latest snapshot before the intermediate nurse event.
	profile.Revision = 3
	profile.SecurityRevision = 3
	if err := r.ApplyStaff(ctx, profile); err != nil {
		t.Fatal(err)
	}
	if _, err := r.RefreshAccount(ctx, "before-role-round-trip"); err == nil {
		t.Fatal("missed intervening role change")
	}
	profile.Role = "nurse"
	profile.Revision = 2
	profile.SecurityRevision = 2
	if err := r.ApplyStaff(ctx, profile); err != nil {
		t.Fatal(err)
	}
	current, err := r.Get(ctx, profile.ID)
	if err != nil || current.Role != "doctor" {
		t.Fatal("stale event changed role")
	}
}
