package storage

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/internal/models"
	"testing"
)

type recordingPublisher struct {
	fail     bool
	profiles []models.Staff
}

func (p *recordingPublisher) Publish(_ context.Context, payload []byte) error {
	if p.fail {
		return errors.New("identity offline")
	}
	var s models.Staff
	if err := json.Unmarshal(payload, &s); err != nil {
		return err
	}
	p.profiles = append(p.profiles, s)
	return nil
}
func TestIdentityOutboxRetriesAndRetainsSecurityRevision(t *testing.T) {
	r := testRepository(t)
	ctx := context.Background()
	s, err := r.Create(ctx, models.Staff{FirstName: "А", LastName: "Б", Position: "Врач", Role: "doctor"})
	if err != nil {
		t.Fatal(err)
	}
	s.Role = "nurse"
	s, err = r.Update(ctx, s.ID, s)
	if err != nil {
		t.Fatal(err)
	}
	s.Role = "doctor"
	s, err = r.Update(ctx, s.ID, s)
	if err != nil {
		t.Fatal(err)
	}
	if s.SecurityRevision == 0 {
		t.Fatal("missing security revision after role round trip")
	}
	publisher := &recordingPublisher{fail: true}
	if err = r.DeliverIdentity(ctx, publisher); err == nil {
		t.Fatal("offline identity acknowledged")
	}
	var count int
	if err = r.pool.QueryRow(ctx, "SELECT count(*) FROM identity_outbox").Scan(&count); err != nil || count != 3 {
		t.Fatalf("lost pending events: %d %v", count, err)
	}
	publisher.fail = false
	if err = r.DeliverIdentity(ctx, publisher); err != nil {
		t.Fatal(err)
	}
	if len(publisher.profiles) != 3 || publisher.profiles[2].SecurityRevision != s.SecurityRevision {
		t.Fatal("events were not delivered in order")
	}
	if err = r.DeliverIdentity(ctx, publisher); err != nil {
		t.Fatal(err)
	}
	if len(publisher.profiles) != 3 {
		t.Fatal("acknowledged events delivered again")
	}
}
