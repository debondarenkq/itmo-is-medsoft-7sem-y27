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
	items, total, err := repository.SearchAccounts(ctx, models.AccountFilter{Limit: 50, Search: "Б А", IsActive: &active})
	if err != nil || total != 1 || len(items) != 1 || items[0].StaffName != "Б А" || items[0].Role != "doctor" {
		t.Fatalf("search by worker: %+v total=%d err=%v", items, total, err)
	}
	if err = repository.Issue(ctx, account, "", "before-edit", -1); err != nil {
		t.Fatal(err)
	}
	username, role := "renamed", "nurse"
	account, err = repository.UpdateAccount(ctx, account.ID, models.AccountChanges{Username: &username, Role: &role, IsActive: &active}, "account:"+account.ID)
	if err != nil || account.Username != username || account.Role != role {
		t.Fatalf("account edit: %+v %v", account, err)
	}
	if _, err = repository.RefreshAccount(ctx, "before-edit"); !errors.Is(err, models.ErrUnauthorized) {
		t.Fatalf("role edit did not revoke refresh: %v", err)
	}
	other, err := repository.Create(ctx, models.Staff{FirstName: "В", LastName: "Г", Position: "Врач", Role: "doctor"})
	if err != nil {
		t.Fatal(err)
	}
	_, err = repository.CreateAccount(ctx, models.Account{ID: uuid.NewString(), Username: "duplicate", StaffID: &other.ID, PasswordHash: "hash"}, "system:test")
	if err != nil {
		t.Fatal(err)
	}
	if err = repository.Issue(ctx, account, "", "must-survive-rollback", -1); err != nil {
		t.Fatal(err)
	}
	username, role = "duplicate", "doctor"
	_, err = repository.UpdateAccount(ctx, account.ID, models.AccountChanges{Username: &username, Role: &role}, "system:test")
	if !errors.Is(err, models.ErrConflict) {
		t.Fatalf("duplicate username accepted: %v", err)
	}
	worker, err := repository.Get(ctx, staff.ID)
	if err != nil || worker.Role != "nurse" {
		t.Fatalf("role was not rolled back: %+v %v", worker, err)
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
	if err := repository.AppendAudit(ctx, "account:"+account.ID, "Вход пользователя", "login"); err != nil {
		t.Fatal(err)
	}
	if err := repository.AppendAudit(ctx, "anonymous:test", "Неуспешная попытка входа", "login_failed"); err != nil {
		t.Fatal(err)
	}
	events, total, err := repository.SearchAuditEvents(ctx, models.AuditFilter{Limit: 1, EventType: "login", From: &before})
	if err != nil || total != 1 || len(events) != 1 || events[0].ActorName != "Б А" {
		t.Fatalf("filtered audit: %+v total=%d err=%v", events, total, err)
	}
	future := time.Now().Add(time.Hour)
	events, total, err = repository.SearchAuditEvents(ctx, models.AuditFilter{Limit: 50, From: &future})
	if err != nil || total != 0 || len(events) != 0 {
		t.Fatalf("period filter: %+v total=%d err=%v", events, total, err)
	}
}
