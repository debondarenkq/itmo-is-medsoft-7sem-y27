package storage

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"

	api "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/lis/api"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/lis/internal/models"
	identity "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/lis/internal/pkg/adapters/identityapi"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/lis/migrations"
)

func testRepository(t *testing.T) *Repository {
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
	return &Repository{Pool: pool}
}

func TestHistorySnapshot(t *testing.T) {
	r := testRepository(t)
	ctx := context.Background()
	h := models.History{ID: uuid.New(), Request: api.FHIRServiceRequest{ID: uuid.New()}, Result: api.FHIRObservation{ID: uuid.New(), ValueString: "142"}, Patient: api.FHIRResource{ID: uuid.New()}, Practitioner: api.FHIRResource{ID: uuid.New()}}
	if err := r.Record(ctx, h); err != nil {
		t.Fatal(err)
	}
	for _, test := range []struct {
		name   string
		change func(*models.History)
	}{{"same result", func(*models.History) {}}, {"new result same request", func(h *models.History) { h.ID = uuid.New(); h.Result.ID = uuid.New(); h.Result.ValueString = "999" }}} {
		t.Run(test.name, func(t *testing.T) {
			copy := h
			test.change(&copy)
			if err := r.Record(ctx, copy); !errors.Is(err, models.ErrConflict) {
				t.Fatalf("duplicate accepted: %v", err)
			}
		})
	}
	stored, err := r.ByRequest(ctx, h.Request.ID)
	if err != nil || stored == nil || stored.Result.ValueString != "142" || stored.Delivered {
		t.Fatalf("snapshot changed: %+v %v", stored, err)
	}
	if err = r.Delivered(ctx, h.ID); err != nil {
		t.Fatal(err)
	}
	history, err := r.History(ctx, "142")
	if err != nil || len(history) != 1 || !history[0].Delivered {
		t.Fatalf("history: %+v %v", history, err)
	}
}

func TestConnection(t *testing.T) {
	for _, test := range []struct {
		name       string
		configured bool
	}{{"authorized", true}, {"not authorized", false}} {
		t.Run(test.name, func(t *testing.T) {
			r := testRepository(t)
			ctx := context.Background()
			session := identity.Session{AccessToken: "test-access", RefreshToken: "test-refresh"}
			if test.configured {
				if err := r.SaveConnection(ctx, models.Connection{Session: session}); err != nil {
					t.Fatal(err)
				}
			}
			stored, err := r.Connection(ctx)
			if test.configured {
				if err != nil || stored.Session.AccessToken != session.AccessToken || stored.Session.RefreshToken != session.RefreshToken {
					t.Fatalf("session mismatch: %v", err)
				}
			} else if !errors.Is(err, models.ErrUnconfigured) {
				t.Fatalf("unexpected session: %v", err)
			}
			session.RefreshToken = "rotated-refresh"
			if err = r.SaveConnection(ctx, models.Connection{Session: session}); err != nil {
				t.Fatal(err)
			}
			stored, err = r.Connection(ctx)
			if err != nil || stored.Session.RefreshToken != session.RefreshToken {
				t.Fatalf("session not updated: %v", err)
			}
		})
	}
}
