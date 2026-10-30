package storage

import (
	"context"
	"log/slog"
	"time"
)

type IdentityPublisher interface {
	Publish(context.Context, []byte) error
}

// Each delivery is acknowledged only after identity commits it. Replays are idempotent.
func (r *Repository) DeliverIdentity(ctx context.Context, publisher IdentityPublisher) error {
	rows, err := r.pool.Query(ctx, "SELECT revision,payload FROM identity_outbox ORDER BY revision LIMIT 100")
	if err != nil {
		return err
	}
	type event struct {
		revision int64
		payload  []byte
	}
	events := []event{}
	for rows.Next() {
		var e event
		if err = rows.Scan(&e.revision, &e.payload); err != nil {
			rows.Close()
			return err
		}
		events = append(events, e)
	}
	err = rows.Err()
	rows.Close()
	if err != nil {
		return err
	}
	for _, e := range events {
		if err = publisher.Publish(ctx, e.payload); err != nil {
			return err
		}
		if _, err = r.pool.Exec(ctx, "DELETE FROM identity_outbox WHERE revision=$1", e.revision); err != nil {
			return err
		}
	}
	return nil
}
func (r *Repository) RunIdentityDelivery(ctx context.Context, publisher IdentityPublisher) {
	ticker := time.NewTicker(250 * time.Millisecond)
	defer ticker.Stop()
	for {
		if err := r.DeliverIdentity(ctx, publisher); err != nil && ctx.Err() == nil {
			slog.Warn("identity delivery will retry", "error", err)
		}
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}
