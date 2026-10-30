package main

import (
	"context"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/internal/app/api"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/internal/app/server"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/internal/pkg/adapters/identity"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/internal/pkg/modules"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/internal/pkg/storage"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/migrations"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"os"
	"strings"
)

func main() {
	key, err := os.ReadFile(os.Getenv("SERVICE_KEY_FILE"))
	if err != nil {
		slog.Error("service key", "error", err)
		os.Exit(1)
	}
	internalKey := strings.TrimSpace(string(key))
	if len(internalKey) < 32 {
		slog.Error("invalid service key")
		os.Exit(1)
	}
	if err := server.Run("staff-service", migrations.FS, func(ctx context.Context, router *mux.Router, pool *pgxpool.Pool) {
		repository := storage.New(pool)
		implementation := api.New(api.Deps{Module: modules.New(modules.Deps{Storage: repository})})
		implementation.Register(router)
		implementation.RegisterInternal(router, internalKey)
		go repository.RunIdentityDelivery(ctx, identity.New(os.Getenv("AUTH_URL"), internalKey))
	}); err != nil {
		slog.Error("service stopped", "error", err)
		os.Exit(1)
	}
}
