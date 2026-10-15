package storage

import (
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/models"
)

const (
	recordPatientUniqueConstraint = "records_patient_id_key"
	recordPatientForeignKey       = "records_patient_id_fkey"
)

func mapNotFoundError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return models.ErrNotFound
	}
	return err
}
