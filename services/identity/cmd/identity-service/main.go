package main

import (
	"context"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/internal/app/api"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/internal/app/server"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/internal/config"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/internal/pkg/adapters/staff"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/internal/pkg/modules"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/internal/pkg/storage"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/internal/transfer"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/migrations"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"net/url"
	"os"
	"strings"
)

func main() {
	if len(os.Args) == 2 && os.Args[1] == "migrate" {
		if err := transfer.Run(context.Background()); err != nil {
			slog.Error("identity transfer failed", "error", err)
			os.Exit(1)
		}
		return
	}
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
	if err := server.Run("identity-service", migrations.FS, func(_ context.Context, router *mux.Router, pool *pgxpool.Pool) {
		repository := storage.New(pool)
		directory := staff.New(base, internalKey)
		repository.UpdateRole = directory.UpdateRole
		module := modules.NewSecurity(repository, key)
		module.Directory = directory
		(&api.SecurityAPI{Module: module, InternalKey: internalKey}).Register(router)
	}); err != nil {
		slog.Error("service stopped", "error", err)
		os.Exit(1)
	}
}
