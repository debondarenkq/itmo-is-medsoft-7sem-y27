package main

import (
	"log/slog"
	"os"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/app/api"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/config"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/pkg/adapters/directories"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/pkg/modules"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/pkg/platform"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/pkg/storage"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/migrations"
)

func main() {
	cfg := config.New()
	client, err := directories.New(cfg.StaffURL, cfg.CatalogURL)
	if err != nil {
		slog.Error("configuration", "error", err)
		os.Exit(1)
	}
	platform.Run("clinical-service", migrations.FS, func(r *mux.Router, pool *pgxpool.Pool) {
		module := modules.New(modules.Deps{Storage: storage.New(pool), Staff: client, Catalog: client})
		api.New(api.Deps{Module: module}).Register(r)
	})
}
