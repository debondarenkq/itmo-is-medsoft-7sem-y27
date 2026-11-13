package storage

import (
	"context"
	"encoding/json"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/models"
	"github.com/jackc/pgx/v5"
)

const eventColumns = "id,record_id,sequence,version,occurred_at,actor,event_type,entity_id,before_data,after_data"

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

func scanEvent(row pgx.Row) (models.Event, error) {
	var event models.Event
	var actor []byte
	if err := row.Scan(&event.ID, &event.RecordID, &event.Sequence, &event.Version, &event.OccurredAt, &actor, &event.Type, &event.EntityID, &event.Before, &event.After); err != nil {
		return event, mapNotFoundError(err)
	}
	return event, json.Unmarshal(actor, &event.Actor)
}
