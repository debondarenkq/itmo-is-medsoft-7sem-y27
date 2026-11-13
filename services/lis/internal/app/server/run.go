package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

const (
	databaseMaxConnections   = 10
	databaseStatementTimeout = 10 * time.Second
	databaseStartupTimeout   = 10 * time.Second
	databaseReadyTimeout     = time.Second
	migrationTimeout         = time.Minute
	shutdownTimeout          = 10 * time.Second
	maxRequestBytes          = 1 << 20
	readHeaderTimeout        = 5 * time.Second
	readTimeout              = 15 * time.Second
	writeTimeout             = 20 * time.Second
	idleTimeout              = time.Minute
)

func Run(name string, migrations fs.FS, register func(*mux.Router, *pgxpool.Pool)) error {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)).With("service", name))
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return errors.New("DATABASE_URL is required")
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return err
	}
	if path := os.Getenv("DATABASE_PASSWORD_FILE"); path != "" {
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		cfg.ConnConfig.Password = strings.TrimSpace(string(data))
	}
	if err := applyMigrations(ctx, cfg.ConnConfig, migrations); err != nil {
		return fmt.Errorf("apply migrations: %w", err)
	}
	cfg.ConnConfig.RuntimeParams["timezone"] = "UTC"
	cfg.ConnConfig.RuntimeParams["statement_timeout"] = strconv.FormatInt(databaseStatementTimeout.Milliseconds(), 10)
	cfg.MaxConns = databaseMaxConnections
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return err
	}
	defer pool.Close()
	startup, cancel := context.WithTimeout(ctx, databaseStartupTimeout)
	err = pool.Ping(startup)
	cancel()
	if err != nil {
		return err
	}

	router := mux.NewRouter()
	router.HandleFunc("/live", func(w http.ResponseWriter, _ *http.Request) {
		writeHealth(w, 200, map[string]any{"service": name, "status": "ok"})
	}).Methods("GET")
	router.HandleFunc("/ready", func(w http.ResponseWriter, request *http.Request) {
		check, cancel := context.WithTimeout(request.Context(), databaseReadyTimeout)
		defer cancel()
		if ctx.Err() != nil || pool.Ping(check) != nil {
			writeHealth(w, 503, map[string]string{"status": "not_ready"})
			return
		}
		writeHealth(w, 200, map[string]string{"status": "ready"})
	}).Methods("GET")
	register(router, pool)
	server := &http.Server{
		Addr: env("HTTP_ADDR", ":8080"),
		Handler: http.MaxBytesHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", "no-store")
			router.ServeHTTP(w, r)
		}), maxRequestBytes),
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
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
		shutdown, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		if err := server.Shutdown(shutdown); err != nil {
			_ = server.Close()
			return err
		}
	}
	return nil
}

func applyMigrations(ctx context.Context, config *pgx.ConnConfig, migrations fs.FS) error {
	db := stdlib.OpenDB(*config)
	defer db.Close()
	provider, err := goose.NewProvider(goose.DialectPostgres, db, migrations)
	if err != nil {
		return err
	}
	migrationCtx, cancel := context.WithTimeout(ctx, migrationTimeout)
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
