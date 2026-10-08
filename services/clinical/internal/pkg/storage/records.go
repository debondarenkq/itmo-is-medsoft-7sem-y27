package storage

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/models"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/pkg/modules"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/pkg/platform"
)

const recordColumns = "id,patient_id,version,created_at,updated_at,state"
const eventColumns = "id,record_id,sequence,version,occurred_at,actor,event_type,entity_id,before_data,after_data"

func scanRecord(row pgx.Row) (models.Record, error) {
	var v models.Record
	var state []byte
	if err := row.Scan(&v.ID, &v.PatientID, &v.Version, &v.CreatedAt, &v.UpdatedAt, &state); err != nil {
		return v, err
	}
	return v, json.Unmarshal(state, &v.State)
}

func scanEvent(row pgx.Row) (models.Event, error) {
	var v models.Event
	var actor []byte
	if err := row.Scan(&v.ID, &v.RecordID, &v.Sequence, &v.Version, &v.OccurredAt, &actor, &v.Type, &v.EntityID, &v.Before, &v.After); err != nil {
		return v, err
	}
	return v, json.Unmarshal(actor, &v.Actor)
}

func (r *Repository) GetRecord(ctx context.Context, id string) (models.Record, error) {
	return scanRecord(r.pool.QueryRow(ctx, "SELECT "+recordColumns+" FROM records WHERE id=$1", id))
}

func (r *Repository) GetPatientRecord(ctx context.Context, id string) (models.Record, error) {
	return scanRecord(r.pool.QueryRow(ctx, "SELECT "+recordColumns+" FROM records WHERE patient_id=$1", id))
}

func insertEvent(ctx context.Context, tx pgx.Tx, e models.Event) error {
	actor, err := json.Marshal(e.Actor)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, "INSERT INTO record_events(id,record_id,sequence,version,occurred_at,actor,event_type,entity_id,before_data,after_data) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)", e.ID, e.RecordID, e.Sequence, e.Version, e.OccurredAt, actor, e.Type, e.EntityID, e.Before, e.After)
	return err
}

func (r *Repository) CreateRecord(ctx context.Context, patientID string, actor models.Actor) (models.Record, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return models.Record{}, err
	}
	defer tx.Rollback(ctx)
	var patient string
	if err = tx.QueryRow(ctx, "SELECT id FROM patients WHERE id=$1 FOR UPDATE", patientID).Scan(&patient); err != nil {
		return models.Record{}, err
	}
	state, _ := json.Marshal(models.EmptyState())
	record, err := scanRecord(tx.QueryRow(ctx, "INSERT INTO records(id,patient_id,state) VALUES($1,$2,$3) RETURNING "+recordColumns, uuid.NewString(), patientID, state))
	var db *pgconn.PgError
	if errors.As(err, &db) && db.Code == "23505" {
		return record, platform.Conflict("RECORD_EXISTS", "У пациента уже есть ЭМК")
	}
	if err != nil {
		return record, err
	}
	after, _ := json.Marshal(map[string]string{"id": record.ID, "patient_id": record.PatientID})
	e := models.Event{ID: uuid.NewString(), RecordID: record.ID, Sequence: 1, Version: 1, OccurredAt: record.CreatedAt, Actor: actor, Type: "record_created", EntityID: record.ID, Before: json.RawMessage("null"), After: after}
	if err = insertEvent(ctx, tx, e); err != nil {
		return record, err
	}
	return record, tx.Commit(ctx)
}

func (r *Repository) Save(ctx context.Context, id string, expected int64, actor models.Actor, commands []models.Command, diagnoses map[string]models.DiagnosisSnapshot) (models.Record, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return models.Record{}, err
	}
	defer tx.Rollback(ctx)
	record, err := scanRecord(tx.QueryRow(ctx, "SELECT "+recordColumns+" FROM records WHERE id=$1 FOR UPDATE", id))
	if err != nil {
		return record, err
	}
	if record.Version != expected {
		return record, platform.Conflict("VERSION_CONFLICT", "ЭМК изменена другим пользователем; обновите карту")
	}
	state, events, err := modules.Apply(record.State, commands, diagnoses)
	if err != nil {
		return record, err
	}
	var sequence int64
	if err = tx.QueryRow(ctx, "SELECT COALESCE(max(sequence),0) FROM record_events WHERE record_id=$1", id).Scan(&sequence); err != nil {
		return record, err
	}
	var at time.Time
	if err = tx.QueryRow(ctx, "SELECT GREATEST(clock_timestamp(),$1::timestamptz + interval '1 microsecond')", record.UpdatedAt).Scan(&at); err != nil {
		return record, err
	}
	record.Version++
	for j := range events {
		e := events[j]
		e.RecordID = id
		e.Sequence = sequence + int64(j) + 1
		e.Version = record.Version
		e.OccurredAt = at
		e.Actor = actor
		if err = insertEvent(ctx, tx, e); err != nil {
			return record, err
		}
	}
	data, err := json.Marshal(state)
	if err != nil {
		return record, err
	}
	if _, err = tx.Exec(ctx, "UPDATE records SET state=$2,version=$3,updated_at=$4 WHERE id=$1", id, data, record.Version, at); err != nil {
		return record, err
	}
	record.State = state
	record.UpdatedAt = at
	return record, tx.Commit(ctx)
}

func (r *Repository) History(ctx context.Context, id string, p platform.Page) ([]models.Event, int, error) {
	var total int
	if err := r.pool.QueryRow(ctx, "SELECT count(*) FROM record_events WHERE record_id=$1", id).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := r.pool.Query(ctx, "SELECT "+eventColumns+" FROM record_events WHERE record_id=$1 ORDER BY sequence LIMIT $2 OFFSET $3", id, p.Limit, p.Offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	events := []models.Event{}
	for rows.Next() {
		v, err := scanEvent(rows)
		if err != nil {
			return nil, 0, err
		}
		events = append(events, v)
	}
	return events, total, rows.Err()
}

func (r *Repository) Events(ctx context.Context, id string, q models.StateQuery) ([]models.Event, error) {
	rows, err := r.pool.Query(ctx, "SELECT "+eventColumns+" FROM record_events WHERE record_id=$1 AND ($2::timestamptz IS NULL OR occurred_at<=$2) AND ($3::bigint IS NULL OR version<=$3) ORDER BY sequence", id, q.At, q.Version)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	events := []models.Event{}
	for rows.Next() {
		v, err := scanEvent(rows)
		if err != nil {
			return nil, err
		}
		events = append(events, v)
	}
	return events, rows.Err()
}
