package main

import (
	"log"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/lis/internal/app/api"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/lis/internal/app/server"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/lis/internal/pkg/adapters"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/lis/internal/pkg/modules"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/lis/internal/pkg/storage"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/lis/migrations"
)

func main() {
	gateway, err := adapters.GatewayAddress()
	if err != nil {
		log.Fatal(err)
	}
	client, err := adapters.HTTPClient()
	if err != nil {
		log.Fatal(err)
	}
	err = server.Run("lis", migrations.FS, func(r *mux.Router, pool *pgxpool.Pool) {
		implementation := &api.Implementation{Module: &modules.Module{Gateway: gateway, Store: &storage.Repository{Pool: pool}, HTTP: client}}
		implementation.Register(r)
	})
	if err != nil {
		log.Fatal(err)
	}
}
