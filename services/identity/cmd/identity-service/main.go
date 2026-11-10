package main

import (
	"context"
	"crypto/ed25519"
	"log/slog"
	"net/url"
	"os"
	"strings"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/internal/app/api"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/internal/app/server"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/internal/config"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/internal/pkg/adapters/staff"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/internal/pkg/modules"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/internal/pkg/storage"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/migrations"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	key, err := config.SigningKey()
	if err != nil {
		slog.Error("signing key", "error", err)
		os.Exit(1)
	}
	data, err := os.ReadFile(os.Getenv("SERVICE_KEY_FILE"))
	if err != nil {
		slog.Error("service key", "error", err)
		os.Exit(1)
	}
	internalKey := strings.TrimSpace(string(data))
	if len(internalKey) < 32 {
		slog.Error("invalid service key")
		os.Exit(1)
	}
	base := os.Getenv("STAFF_URL")
	u, err := url.Parse(base)
	if err != nil || u.Host == "" || (u.Scheme != "http" && u.Scheme != "https") {
		slog.Error("STAFF_URL is required")
		os.Exit(1)
	}
	directory, err := staff.New(base, internalKey)
	if err != nil {
		slog.Error("staff client", "error", err)
		os.Exit(1)
	}
	if err := server.Run("identity-service", migrations.FS, func(_ context.Context, router *mux.Router, pool *pgxpool.Pool) {
		api.New(api.Deps{Module: newSecurity(pool, key, directory), InternalKey: internalKey}).Register(router)
	}); err != nil {
		slog.Error("service stopped", "error", err)
		os.Exit(1)
	}
}

func newSecurity(pool *pgxpool.Pool, key ed25519.PrivateKey, directory *staff.Directory) *modules.Security {
	repository := storage.New(pool)
	module := modules.NewSecurity(repository, key)
	module.Directory = directory
	return module
}
