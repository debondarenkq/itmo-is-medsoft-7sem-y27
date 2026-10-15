package api

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/api"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/models"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/validation"
)

func (implementation *Implementation) GetRecordHistory(ctx context.Context, request contract.GetRecordHistoryRequestObject) (contract.GetRecordHistoryResponseObject, error) {
	if err := validation.RequireID(request.ID); err != nil {
		return nil, err
	}
	page := pageFromParams(request.Params.Limit, request.Params.Offset, nil)
	events, total, err := implementation.Module.History(ctx, request.ID.String(), page)
	if err != nil {
		return nil, err
	}
	items := make([]contract.Event, len(events))
	for index, event := range events {
		items[index], err = toEvent(event)
		if err != nil {
			return nil, err
		}
	}
	return contract.GetRecordHistory200JSONResponse{Items: items, Total: total, Limit: page.Limit, Offset: page.Offset}, nil
}

func toEvent(value models.Event) (contract.Event, error) {
	id, err := uuid.Parse(value.ID)
	if err != nil {
		return contract.Event{}, fmt.Errorf("map event id: %w", err)
	}
	recordID, err := uuid.Parse(value.RecordID)
	if err != nil {
		return contract.Event{}, fmt.Errorf("map event record id: %w", err)
	}
	entityID, err := uuid.Parse(value.EntityID)
	if err != nil {
		return contract.Event{}, fmt.Errorf("map event entity id: %w", err)
	}
	actor, err := toActor(value.Actor)
	if err != nil {
		return contract.Event{}, err
	}
	before, err := snapshot(value.Before)
	if err != nil {
		return contract.Event{}, fmt.Errorf("map event before: %w", err)
	}
	after, err := snapshot(value.After)
	if err != nil {
		return contract.Event{}, fmt.Errorf("map event after: %w", err)
	}
	return contract.Event{
		ID:         id,
		RecordID:   recordID,
		Sequence:   value.Sequence,
		Version:    value.Version,
		OccurredAt: value.OccurredAt,
		Actor:      actor,
		Type:       contract.EventType(value.Type),
		EntityID:   entityID,
		Before:     before,
		After:      after,
	}, nil
}

func toActor(value models.Actor) (contract.Actor, error) {
	id, err := uuid.Parse(value.ID)
	if err != nil {
		return contract.Actor{}, fmt.Errorf("map actor id: %w", err)
	}
	return contract.Actor{
		ID:        id,
		FirstName: value.FirstName,
		LastName:  value.LastName,
		Position:  value.Position,
	}, nil
}

func snapshot(value json.RawMessage) (*map[string]any, error) {
	var result *map[string]any
	if err := json.Unmarshal(value, &result); err != nil {
		return nil, fmt.Errorf("decode event snapshot: %w", err)
	}
	return result, nil
}
