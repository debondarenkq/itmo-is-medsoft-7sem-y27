package storage

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	recordPatientUniqueConstraint = "records_patient_id_key"
	recordPatientForeignKey       = "records_patient_id_fkey"
)

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

func insertEvent(ctx context.Context, tx pgx.Tx, event models.Event) error {
	actor, err := json.Marshal(event.Actor)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, "INSERT INTO record_events(id,record_id,sequence,version,occurred_at,actor,event_type,entity_id,before_data,after_data) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)", event.ID, event.RecordID, event.Sequence, event.Version, event.OccurredAt, actor, event.Type, event.EntityID, event.Before, event.After)
	return err
}
