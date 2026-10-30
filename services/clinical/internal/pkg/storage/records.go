package storage

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/models"
)

const recordColumns = "id,patient_id,version,created_at,updated_at,state,number"
const eventColumns = "id,record_id,sequence,version,occurred_at,actor,event_type,entity_id,before_data,after_data"

func scanRecord(row pgx.Row) (models.Record, error) {
	var record models.Record
	var state []byte
	if err := row.Scan(&record.ID, &record.PatientID, &record.Version, &record.CreatedAt, &record.UpdatedAt, &state, &record.Number); err != nil {
		return record, mapNotFoundError(err)
	}
	return record, json.Unmarshal(state, &record.State)
}

func scanEvent(row pgx.Row) (models.Event, error) {
	var event models.Event
	var actor []byte
	if err := row.Scan(&event.ID, &event.RecordID, &event.Sequence, &event.Version, &event.OccurredAt, &actor, &event.Type, &event.EntityID, &event.Before, &event.After); err != nil {
		return event, mapNotFoundError(err)
	}
	return event, json.Unmarshal(actor, &event.Actor)
}

func (repository *Repository) GetRecord(ctx context.Context, id string) (models.Record, error) {
	return scanRecord(repository.pool.QueryRow(ctx, "SELECT "+recordColumns+" FROM records WHERE id=$1", id))
}

func (repository *Repository) GetPatientRecord(ctx context.Context, id string) (models.Record, error) {
	return scanRecord(repository.pool.QueryRow(ctx, "SELECT "+recordColumns+" FROM records WHERE patient_id=$1", id))
}

func insertEvent(ctx context.Context, tx pgx.Tx, event models.Event) error {
	actor, err := json.Marshal(event.Actor)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, "INSERT INTO record_events(id,record_id,sequence,version,occurred_at,actor,event_type,entity_id,before_data,after_data) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)", event.ID, event.RecordID, event.Sequence, event.Version, event.OccurredAt, actor, event.Type, event.EntityID, event.Before, event.After)
	return err
}

func (repository *Repository) CreateRecord(ctx context.Context, patientID string, actor models.Actor) (models.Record, error) {
	tx, err := repository.pool.Begin(ctx)
	if err != nil {
		return models.Record{}, err
	}
	defer tx.Rollback(ctx)
	var patient string
	if err = tx.QueryRow(ctx, "SELECT id FROM patients WHERE id=$1 FOR UPDATE", patientID).Scan(&patient); err != nil {
		return models.Record{}, mapNotFoundError(err)
	}
	state, _ := json.Marshal(models.EmptyState())
	record, err := scanRecord(tx.QueryRow(ctx, "INSERT INTO records(id,patient_id,state) VALUES($1,$2,$3) RETURNING "+recordColumns, uuid.NewString(), patientID, state))
	var databaseError *pgconn.PgError
	if errors.As(err, &databaseError) && databaseError.Code == pgerrcode.UniqueViolation && databaseError.ConstraintName == recordPatientUniqueConstraint {
		return record, models.ErrAlreadyExists
	}
	if err != nil {
		return record, err
	}
	after, _ := json.Marshal(map[string]string{"id": record.ID, "patient_id": record.PatientID})
	event := models.Event{ID: uuid.NewString(), RecordID: record.ID, Sequence: 1, Version: 1, OccurredAt: record.CreatedAt, Actor: actor, Type: models.EventRecordCreated, EntityID: record.ID, Before: json.RawMessage("null"), After: after}
	if err = insertEvent(ctx, tx, event); err != nil {
		return record, err
	}
	return record, tx.Commit(ctx)
}

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

func (repository *Repository) History(ctx context.Context, id string, page models.Page) ([]models.Event, int, error) {
	var total int
	if err := repository.pool.QueryRow(ctx, "SELECT count(*) FROM record_events WHERE record_id=$1", id).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := repository.pool.Query(ctx, "SELECT "+eventColumns+" FROM record_events WHERE record_id=$1 ORDER BY sequence LIMIT $2 OFFSET $3", id, page.Limit, page.Offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	events := []models.Event{}
	for rows.Next() {
		event, err := scanEvent(rows)
		if err != nil {
			return nil, 0, err
		}
		events = append(events, event)
	}
	return events, total, rows.Err()
}

func (repository *Repository) Events(ctx context.Context, id string, stateQuery models.StateQuery) ([]models.Event, error) {
	rows, err := repository.pool.Query(ctx, "SELECT "+eventColumns+" FROM record_events WHERE record_id=$1 AND ($2::timestamptz IS NULL OR occurred_at<=$2) AND ($3::bigint IS NULL OR version<=$3) ORDER BY sequence", id, stateQuery.At, stateQuery.Version)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	events := []models.Event{}
	for rows.Next() {
		event, err := scanEvent(rows)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, rows.Err()
}
