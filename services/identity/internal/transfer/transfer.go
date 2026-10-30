// Package transfer implements an offline, one-time upgrade. Runtime identity has no staff DB access.
package transfer

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/migrations"
	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"os"
)

func Run(ctx context.Context) error {
	return Migrate(ctx, os.Getenv("LEGACY_DATABASE_URL"), os.Getenv("DATABASE_URL"))
}
func Migrate(ctx context.Context, sourceDSN, targetDSN string) error {
	db, err := sql.Open("pgx", targetDSN)
	if err != nil {
		return err
	}
	defer db.Close()
	provider, err := goose.NewProvider(goose.DialectPostgres, db, migrations.FS)
	if err != nil {
		return err
	}
	if _, err = provider.Up(ctx); err != nil {
		return err
	}
	source, err := pgx.Connect(ctx, sourceDSN)
	if err != nil {
		return err
	}
	defer source.Close(ctx)
	var exists bool
	if err = source.QueryRow(ctx, "SELECT to_regclass('accounts') IS NOT NULL").Scan(&exists); err != nil || !exists {
		return err
	}
	from, err := source.Begin(ctx)
	if err != nil {
		return err
	}
	defer from.Rollback(ctx)
	// The upgrade script stops the APIs first. Locks also exclude stray legacy writers during copying.
	if _, err = from.Exec(ctx, "LOCK TABLE staff,accounts,refresh_tokens,audit_events,security_settings IN ACCESS EXCLUSIVE MODE"); err != nil {
		return err
	}
	target, err := pgx.Connect(ctx, targetDSN)
	if err != nil {
		return err
	}
	defer target.Close(ctx)
	to, err := target.Begin(ctx)
	if err != nil {
		return err
	}
	defer to.Rollback(ctx)
	var imported bool
	if err = to.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM legacy_import)").Scan(&imported); err != nil {
		return err
	}
	tables := []struct {
		name, query string
		columns     []string
	}{
		{"staff_profiles", "SELECT id,first_name,last_name,role,deleted_at,0::bigint FROM staff ORDER BY id", []string{"id", "first_name", "last_name", "role", "deleted_at", "revision"}},
		{"accounts", "SELECT id,username,password_hash,created_at,updated_at,is_active,staff_id,is_admin,totp_secret,last_totp_step FROM accounts ORDER BY id", []string{"id", "username", "password_hash", "created_at", "updated_at", "is_active", "staff_id", "is_admin", "totp_secret", "last_totp_step"}},
		{"refresh_tokens", "SELECT token_hash,account_id,expires_at FROM refresh_tokens ORDER BY token_hash", []string{"token_hash", "account_id", "expires_at"}},
		{"audit_events", "SELECT event_id,event_time,actor,message,event_type FROM audit_events ORDER BY event_id", []string{"event_id", "event_time", "actor", "message", "event_type"}},
		{"security_settings", "SELECT id,setup_secret FROM security_settings ORDER BY id", []string{"id", "setup_secret"}},
	}
	digest := sha256.New()
	if !imported {
		var populated bool
		if err = to.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM accounts) OR EXISTS(SELECT 1 FROM audit_events) OR EXISTS(SELECT 1 FROM security_settings)").Scan(&populated); err != nil {
			return err
		}
		if populated {
			return fmt.Errorf("identity database is already populated; refusing to overwrite it")
		}

	}
	for _, table := range tables {
		rows, err := from.Query(ctx, table.query)
		if err != nil {
			return err
		}
		values := [][]any{}
		for rows.Next() {
			v, e := rows.Values()
			if e != nil {
				rows.Close()
				return e
			}
			values = append(values, v)
		}
		err = rows.Err()
		rows.Close()
		if err != nil {
			return err
		}
		encoded, err := json.Marshal(values)
		if err != nil {
			return err
		}
		digest.Write(encoded)
		if imported {
			continue
		}
		count, err := to.CopyFrom(ctx, pgx.Identifier{table.name}, table.columns, pgx.CopyFromRows(values))
		if err != nil {
			return fmt.Errorf("copy %s: %w", table.name, err)
		}
		if count != int64(len(values)) {
			return fmt.Errorf("incomplete copy of %s", table.name)
		}
	}
	fingerprint := fmt.Sprintf("%x", digest.Sum(nil))
	if imported {
		var previous string
		if err = to.QueryRow(ctx, "SELECT source_hash FROM legacy_import WHERE id").Scan(&previous); err != nil {
			return err
		}
		if previous != fingerprint {
			return fmt.Errorf("legacy data changed after import; refusing to discard changes")
		}
	} else if _, err = to.Exec(ctx, "INSERT INTO legacy_import VALUES(true,$1)", fingerprint); err != nil {
		return err
	}

	// Commit the destination first. A crash before the source marker can safely be retried.
	if err = to.Commit(ctx); err != nil {
		return err
	}
	if _, err = from.Exec(ctx, "CREATE TABLE IF NOT EXISTS _identity_transfer(id boolean PRIMARY KEY CHECK(id)); INSERT INTO _identity_transfer VALUES(true) ON CONFLICT DO NOTHING"); err != nil {
		return err
	}
	return from.Commit(ctx)
}
