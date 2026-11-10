package storage

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/internal/models"
	"github.com/google/uuid"
)

func TestAccountSearchAndAtomicEditing(t *testing.T) {
	repository := testRepository(t)
	ctx := context.Background()
	staff, account := securityAccount(t, repository)
	active := true
	items, total, err := repository.SearchAccounts(ctx, models.AccountFilter{Limit: 50, Search: "Б А", IsActive: &active, StaffIDs: []string{staff.ID}})
	if err != nil || total != 1 || len(items) != 1 {
		t.Fatalf("search by worker: %+v total=%d err=%v", items, total, err)
	}
	if err = repository.Issue(ctx, account, "", "before-edit", -1); err != nil {
		t.Fatal(err)
	}
	username := "renamed"
	account, err = repository.UpdateAccount(ctx, account.ID, models.AccountChanges{Username: &username, RevokeSessions: true, IsActive: &active}, "account:"+account.ID)
	if err != nil || account.Username != username {
		t.Fatalf("account edit: %+v %v", account, err)
	}
	if _, err = repository.RefreshAccount(ctx, "before-edit"); !errors.Is(err, models.ErrUnauthorized) {
		t.Fatalf("role edit did not revoke refresh: %v", err)
	}
	other := models.Staff{ID: uuid.NewString()}
	_, err = repository.CreateAccount(ctx, models.Account{ID: uuid.NewString(), Username: "duplicate", StaffID: &other.ID, PasswordHash: "hash"}, "system:test")
	if err != nil {
		t.Fatal(err)
	}
	if err = repository.Issue(ctx, account, "", "must-survive-rollback", -1); err != nil {
		t.Fatal(err)
	}
	username = "duplicate"
	_, err = repository.UpdateAccount(ctx, account.ID, models.AccountChanges{Username: &username, RevokeSessions: true}, "system:test")
	if !errors.Is(err, models.ErrConflict) {
		t.Fatalf("duplicate username accepted: %v", err)
	}
	if _, err = repository.RefreshAccount(ctx, "must-survive-rollback"); err != nil {
		t.Fatalf("failed edit revoked a session: %v", err)
	}
	inactive := false
	items, total, err = repository.SearchAccounts(ctx, models.AccountFilter{Limit: 1, Search: "renamed", IsActive: &inactive})
	if err != nil || total != 0 || len(items) != 0 {
		t.Fatalf("activity filter: %+v %d %v", items, total, err)
	}
}

func TestAuditFiltersAndActorNames(t *testing.T) {
	repository := testRepository(t)
	ctx := context.Background()
	_, account := securityAccount(t, repository)
	before := time.Now().Add(-time.Minute)
	if err := repository.AppendAudit(ctx, "account:"+account.ID, "Вход пользователя", models.AuditEventKindLogin); err != nil {
		t.Fatal(err)
	}
	if err := repository.AppendAudit(ctx, "anonymous:test", "Неуспешная попытка входа", models.AuditEventKindLoginFailed); err != nil {
		t.Fatal(err)
	}
	events, total, err := repository.SearchAuditEvents(ctx, models.AuditFilter{Limit: 1, EventType: "login", From: &before})
	if err != nil || total != 1 || len(events) != 1 || events[0].ActorName != "doctor" || events[0].ActorStaffID == nil {
		t.Fatalf("filtered audit: %+v total=%d err=%v", events, total, err)
	}
	future := time.Now().Add(time.Hour)
	events, total, err = repository.SearchAuditEvents(ctx, models.AuditFilter{Limit: 50, From: &future})
	if err != nil || total != 0 || len(events) != 0 {
		t.Fatalf("period filter: %+v total=%d err=%v", events, total, err)
	}
}
