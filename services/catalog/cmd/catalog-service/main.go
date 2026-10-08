package main

import (
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/catalog/internal/app/api"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/catalog/internal/pkg/modules"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/catalog/internal/pkg/platform"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/catalog/internal/pkg/storage"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/catalog/migrations"
)

func main() {
	platform.Run("catalog-service", migrations.FS, func(r *mux.Router, pool *pgxpool.Pool) {
		repository := storage.New(pool)
		module := modules.New(modules.Deps{Storage: repository})
		api.New(api.Deps{Module: module}).Register(r)
	})
}
