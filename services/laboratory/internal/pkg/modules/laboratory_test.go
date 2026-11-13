package modules

import (
	"context"
	"errors"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"

	api "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/laboratory/api"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/laboratory/internal/models"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/laboratory/internal/pkg/storage"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/laboratory/migrations"
)

func testModule(t *testing.T) *Module {
	t.Helper()
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("set TEST_DATABASE_URL to run PostgreSQL integration tests")
	}
	ctx := context.Background()
	admin, err := pgx.Connect(ctx, dsn)
	if err != nil {
		t.Fatal(err)
	}
	schema := "test_" + uuid.New().String()
	quoted := pgx.Identifier{schema}.Sanitize()
	if _, err = admin.Exec(ctx, "CREATE SCHEMA "+quoted); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if _, err := admin.Exec(ctx, "DROP SCHEMA "+quoted+" CASCADE"); err != nil {
			t.Error(err)
		}
		_ = admin.Close(ctx)
	})
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		t.Fatal(err)
	}
	cfg.ConnConfig.RuntimeParams["search_path"] = quoted
	cfg.ConnConfig.RuntimeParams["timezone"] = "UTC"
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(pool.Close)
	db := stdlib.OpenDB(*cfg.ConnConfig)
	defer db.Close()
	provider, err := goose.NewProvider(goose.DialectPostgres, db, migrations.FS)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = provider.Up(ctx); err != nil {
		t.Fatal(err)
	}
	if _, err = provider.Up(ctx); err != nil {
		t.Fatalf("repeat migration: %v", err)
	}
	return &Module{Store: storage.New(pool)}
}

func fixture(t *testing.T, m *Module) (api.FHIRServiceRequest, api.FHIRObservation) {
	t.Helper()
	now := time.Now().UTC()
	request := api.FHIRServiceRequest{ID: uuid.New(), ResourceType: api.ServiceRequest, Status: api.Active, Intent: api.Order, AuthoredOn: now.Add(-time.Minute), Subject: api.FHIRReference{Reference: "Patient/" + uuid.NewString()}, Requester: api.FHIRReference{Reference: "Practitioner/" + uuid.NewString()}, Code: api.FHIRCodeableReference{Concept: api.FHIRCodeableConcept{Coding: []api.FHIRCoding{{System: api.HTTPLoincOrg, Code: "718-7", Display: "Гемоглобин"}}}}}
	patient, _ := ReferenceID(request.Subject.Reference, "Patient")
	if err := m.Store.CreateRequests(context.Background(), patient, []api.FHIRServiceRequest{request}); err != nil {
		t.Fatal(err)
	}
	return request, api.FHIRObservation{ID: uuid.New(), ResourceType: api.Observation, Status: api.Final, Code: request.Code.Concept, Subject: request.Subject, BasedOn: []api.FHIRReference{{Reference: "ServiceRequest/" + request.ID.String()}}, EffectiveDateTime: now, Issued: now, ValueString: "142"}
}
func TestObservationValidation(t *testing.T) {
	m := testModule(t)
	ctx := context.Background()
	for _, test := range []struct {
		name   string
		change func(*api.FHIRObservation)
		valid  bool
	}{
		{"valid", func(*api.FHIRObservation) {}, true},
		{"wrong patient", func(o *api.FHIRObservation) { o.Subject.Reference = "Patient/" + uuid.NewString() }, false},
		{"wrong code", func(o *api.FHIRObservation) {
			o.Code = api.FHIRCodeableConcept{Coding: []api.FHIRCoding{{System: api.HTTPLoincOrg, Code: "2345-7", Display: "Глюкоза"}}}
		}, false},
		{"missing reference", func(o *api.FHIRObservation) { o.BasedOn = nil }, false},
		{"multiple references", func(o *api.FHIRObservation) { o.BasedOn = append(o.BasedOn, o.BasedOn[0]) }, false},
		{"empty value", func(o *api.FHIRObservation) { o.ValueString = "  " }, false},
		{"before order", func(o *api.FHIRObservation) { o.EffectiveDateTime = o.EffectiveDateTime.Add(-time.Hour) }, false},
		{"future date", func(o *api.FHIRObservation) { o.EffectiveDateTime = o.EffectiveDateTime.Add(time.Hour) }, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			request, result := fixture(t, m)
			test.change(&result)
			_, err := m.SaveObservation(ctx, result)
			if (err == nil) != test.valid {
				t.Fatalf("valid=%v: %v", test.valid, err)
			}
			requests, e := m.Store.Requests(ctx, "", request.ID.String())
			if e != nil {
				t.Fatal(e)
			}
			if (requests[0].Status == api.Completed) != test.valid {
				t.Fatal("request status changed after failed validation")
			}
		})
	}
}
func TestConcurrentObservation(t *testing.T) {
	m := testModule(t)
	request, result := fixture(t, m)
	ctx := context.Background()
	var wg sync.WaitGroup
	results := make(chan error, 8)
	for range 8 {
		wg.Go(func() { copy := result; copy.ID = uuid.New(); _, err := m.SaveObservation(ctx, copy); results <- err })
	}
	wg.Wait()
	close(results)
	successes := 0
	for err := range results {
		if err == nil {
			successes++
		} else if !errors.Is(err, models.ErrConflict) {
			t.Fatal(err)
		}
	}
	if successes != 1 {
		t.Fatalf("results accepted: %d", successes)
	}
	patient, _ := ReferenceID(request.Subject.Reference, "Patient")
	observations, err := m.Store.Observations(ctx, patient.String(), "")
	if err != nil || len(observations) != 1 {
		t.Fatalf("observations=%d: %v", len(observations), err)
	}
}
