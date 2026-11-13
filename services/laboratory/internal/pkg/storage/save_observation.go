package storage

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/laboratory/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func (r *Repository) SaveObservation(ctx context.Context, requestID, patientID uuid.UUID, result models.Observation, validate func(models.Request) error) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	var data []byte
	if err = tx.QueryRow(ctx, "SELECT resource FROM lab_requests WHERE id=$1 FOR UPDATE", requestID).Scan(&data); err != nil {
		return databaseError(err)
	}
	var request models.Request
	if err = json.Unmarshal(data, &request); err != nil {
		return err
	}
	if err = validate(request); err != nil {
		return err
	}
	data, err = json.Marshal(result)
	if err != nil {
		return err
	}
	if _, err = tx.Exec(ctx, "INSERT INTO lab_observations(id,request_id,patient_id,resource) VALUES($1,$2,$3,$4)", result.ID, requestID, patientID, data); err != nil {
		return databaseError(err)
	}
	_, err = tx.Exec(ctx, `UPDATE lab_requests SET resource=jsonb_set(resource,'{status}','"completed"') WHERE id=$1`, requestID)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func databaseError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return models.ErrNotFound
	}
	var pe *pgconn.PgError
	if errors.As(err, &pe) {
		switch pe.Code {
		case "23505":
			return models.ErrConflict
		case "23514":
			return models.ErrInvalidInput
		}
	}
	return err
}
