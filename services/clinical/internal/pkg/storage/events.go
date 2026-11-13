package storage

import (
	"context"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/models"
)

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
