package storage

import (
	"context"
	"encoding/json"
	"time"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/models"
)

func (repository *Repository) Save(ctx context.Context, id string, actor models.Actor, update func(models.Record) (models.State, []models.Event, error)) (models.Record, error) {
	tx, err := repository.pool.Begin(ctx)
	if err != nil {
		return models.Record{}, err
	}
	defer tx.Rollback(ctx)
	record, err := scanRecord(tx.QueryRow(ctx, "SELECT "+recordColumns+" FROM records WHERE id=$1 FOR UPDATE", id))
	if err != nil {
		return record, err
	}
	state, events, err := update(record)
	if err != nil {
		return record, err
	}
	var sequence int64
	if err = tx.QueryRow(ctx, "SELECT COALESCE(max(sequence),0) FROM record_events WHERE record_id=$1", id).Scan(&sequence); err != nil {
		return record, err
	}
	var occurredAt time.Time
	if err = tx.QueryRow(ctx, "SELECT GREATEST(clock_timestamp(),$1::timestamptz + interval '1 microsecond')", record.UpdatedAt).Scan(&occurredAt); err != nil {
		return record, err
	}
	record.Version++
	for index := range events {
		event := events[index]
		event.RecordID = id
		event.Sequence = sequence + int64(index) + 1
		event.Version = record.Version
		event.OccurredAt = occurredAt
		event.Actor = actor
		if err = insertEvent(ctx, tx, event); err != nil {
			return record, err
		}
	}
	data, err := json.Marshal(state)
	if err != nil {
		return record, err
	}
	if _, err = tx.Exec(ctx, "UPDATE records SET state=$2,version=$3,updated_at=$4 WHERE id=$1", id, data, record.Version, occurredAt); err != nil {
		return record, err
	}
	record.State = state
	record.UpdatedAt = occurredAt
	return record, tx.Commit(ctx)
}
