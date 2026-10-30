package main

import (
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/internal/config"
	"log/slog"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/internal/app/api"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/internal/app/server"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/internal/pkg/modules"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/internal/pkg/storage"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/staff/migrations"
)

func main() {
	key, err := config.SigningKey()
	if err != nil {
		slog.Error("signing key", "error", err)
		os.Exit(1)
	}
	internalKey, err := os.ReadFile(os.Getenv("SERVICE_KEY_FILE"))
	if err != nil {
		slog.Error("service key", "error", err)
		os.Exit(1)
	}
	if err := server.Run("staff-service", migrations.FS, func(router *mux.Router, pool *pgxpool.Pool) {
		repository := storage.New(pool)
		module := modules.New(modules.Deps{Storage: repository})
		api.New(api.Deps{Module: module}).Register(router)
		(&api.SecurityAPI{Module: modules.NewSecurity(repository, key), InternalKey: strings.TrimSpace(string(internalKey))}).Register(router)
	}); err != nil {
		slog.Error("service stopped", "error", err)
		os.Exit(1)
	}
}
