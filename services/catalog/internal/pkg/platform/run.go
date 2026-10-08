package platform

import (
	"context"
	"database/sql"
	"errors"
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

func Env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func Run(name string, migrations fs.FS, register func(*mux.Router, *pgxpool.Pool)) {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)).With("service", name))
	if err := run(name, migrations, register); err != nil {
		slog.Error("service stopped", "error", err)
		os.Exit(1)
	}
}

func run(name string, migrations fs.FS, register func(*mux.Router, *pgxpool.Pool)) error {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return errors.New("DATABASE_URL is required")
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if len(os.Args) > 1 {
		if os.Args[1] != "migrate" {
			return errors.New("usage: service [migrate]")
		}
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
	r := mux.NewRouter()
	r.HandleFunc("/live", func(w http.ResponseWriter, _ *http.Request) {
		JSON(w, 200, map[string]any{"service": name, "status": "ok"})
	}).Methods("GET")
	r.HandleFunc("/ready", func(w http.ResponseWriter, req *http.Request) {
		check, cancel := context.WithTimeout(req.Context(), time.Second)
		defer cancel()
		if ctx.Err() != nil || pool.Ping(check) != nil {
			JSON(w, 503, map[string]string{"status": "not_ready"})
			return
		}
		JSON(w, 200, map[string]string{"status": "ready"})
	}).Methods("GET")
	register(r, pool)
	r.NotFoundHandler = Handle(func(http.ResponseWriter, *http.Request) error { return NotFound() })
	r.MethodNotAllowedHandler = Handle(func(http.ResponseWriter, *http.Request) error {
		return &Error{Status: 405, Code: "METHOD_NOT_ALLOWED", Message: "Метод не поддерживается"}
	})
	server := &http.Server{Addr: Env("HTTP_ADDR", ":8080"), Handler: r, ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 15 * time.Second, WriteTimeout: 20 * time.Second, IdleTimeout: 60 * time.Second}
	result := make(chan error, 1)
	go func() { slog.Info("listening", "address", server.Addr); result <- server.ListenAndServe() }()
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
