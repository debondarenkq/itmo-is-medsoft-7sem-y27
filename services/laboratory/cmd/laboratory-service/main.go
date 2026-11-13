package main

import (
	"log"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/laboratory/internal/app/api"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/laboratory/internal/app/server"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/laboratory/internal/pkg/adapters"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/laboratory/internal/pkg/modules"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/laboratory/internal/pkg/storage"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/laboratory/migrations"
)

func main() {
	directories, err := adapters.NewDirectories()
	if err != nil {
		log.Fatal(err)
	}
	err = server.Run("laboratory", migrations.FS, func(r *mux.Router, pool *pgxpool.Pool) {
		implementation := &api.Implementation{Module: &modules.Module{Store: storage.New(pool), Directories: directories}}
		implementation.Register(r)
	})
	if err != nil {
		log.Fatal(err)
	}
}
