package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/pkg/adapters/security"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func Run(name string, migrations fs.FS, register func(*mux.Router, *pgxpool.Pool)) error {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)).With("service", name))
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return errors.New("DATABASE_URL is required")
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if err := applyMigrations(ctx, dsn, migrations); err != nil {
		return fmt.Errorf("apply migrations: %w", err)
	}
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return err
	}
	cfg.ConnConfig.RuntimeParams["timezone"] = "UTC"
	cfg.ConnConfig.RuntimeParams["statement_timeout"] = "10000"
	cfg.MaxConns = 10
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return err
	}
	defer pool.Close()
	startup, cancel := context.WithTimeout(ctx, 10*time.Second)
	err = pool.Ping(startup)
	cancel()
	if err != nil {
		return err
	}
	authorization, err := security.New()
	if err != nil {
		return err
	}
	router := mux.NewRouter()
	router.HandleFunc("/live", func(w http.ResponseWriter, _ *http.Request) {
		writeHealth(w, 200, map[string]any{"service": name, "status": "ok"})
	}).Methods("GET")
	router.HandleFunc("/ready", func(w http.ResponseWriter, request *http.Request) {
		check, cancel := context.WithTimeout(request.Context(), time.Second)
		defer cancel()
		if ctx.Err() != nil || pool.Ping(check) != nil {
			writeHealth(w, 503, map[string]string{"status": "not_ready"})
			return
		}
		writeHealth(w, 200, map[string]string{"status": "ready"})
	}).Methods("GET")
	register(router, pool)
	server := &http.Server{
		Addr:              env("HTTP_ADDR", ":8080"),
		Handler:           http.MaxBytesHandler(authorization.Middleware(router), 1<<20),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      20 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	result := make(chan error, 1)
	go func() {
		slog.Info("listening", "address", server.Addr)
		result <- server.ListenAndServe()
	}()
	select {
	case err := <-result:
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	case <-ctx.Done():
		shutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdown); err != nil {
			_ = server.Close()
			return err
		}
	}
	return nil
}

func applyMigrations(ctx context.Context, dsn string, migrations fs.FS) error {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return err
	}
	defer db.Close()
	provider, err := goose.NewProvider(goose.DialectPostgres, db, migrations)
	if err != nil {
		return err
	}
	migrationCtx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()
	_, err = provider.Up(migrationCtx)
	return err
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func writeHealth(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		slog.Warn("encode health response", "error", err)
	}
}
