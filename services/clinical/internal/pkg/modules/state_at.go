package modules

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/models"
)

func (module *Module) StateAt(ctx context.Context, id string, stateQuery models.StateQuery) (models.Record, error) {
	record, err := module.Storage.GetRecord(ctx, id)
	if err != nil {
		return record, err
	}
	if stateQuery.Version != nil && (*stateQuery.Version < 1 || *stateQuery.Version > record.Version) {
		return record, models.ErrNotFound
	}
	if stateQuery.At != nil && stateQuery.At.Before(record.CreatedAt) {
		return record, models.ErrNotFound
	}
	events, err := module.Storage.Events(ctx, id, stateQuery)
	if err != nil {
		return record, err
	}
	if len(events) == 0 {
		return record, models.ErrNotFound
	}
	record.State, err = Replay(events)
	last := events[len(events)-1]
	record.Version = last.Version
	record.UpdatedAt = last.OccurredAt
	return record, err
}

func Replay(events []models.Event) (models.State, error) {
	state := models.EmptyState()
	for _, event := range events {
		switch event.Type {
		case models.EventRecordCreated:
		case models.EventDiagnosisAdded:
			var diagnosis models.RecordDiagnosis
			if err := json.Unmarshal(event.After, &diagnosis); err != nil {
				return state, err
			}
			state.Diagnoses = append(state.Diagnoses, diagnosis)
		case models.EventDiagnosisRemoved, models.EventDiagnosisUpdated:
			index := slices.IndexFunc(state.Diagnoses, func(diagnosis models.RecordDiagnosis) bool { return diagnosis.ID == event.EntityID })
			if index < 0 {
				return state, fmt.Errorf("invalid diagnosis event %s", event.ID)
			}
			if event.Type == models.EventDiagnosisRemoved {
				state.Diagnoses = slices.Delete(state.Diagnoses, index, index+1)
			} else if err := json.Unmarshal(event.After, &state.Diagnoses[index]); err != nil {
				return state, err
			}
		case models.EventPrescriptionAdded:
			var prescription models.Prescription
			if err := json.Unmarshal(event.After, &prescription); err != nil {
				return state, err
			}
			state.Prescriptions = append(state.Prescriptions, prescription)
		case models.EventPrescriptionCancelled, models.EventPrescriptionRemoved, models.EventPrescriptionUpdated:
			index := slices.IndexFunc(state.Prescriptions, func(prescription models.Prescription) bool { return prescription.ID == event.EntityID })
			if index < 0 {
				return state, fmt.Errorf("invalid prescription event %s", event.ID)
			}
			if event.Type == models.EventPrescriptionRemoved {
				state.Prescriptions = slices.Delete(state.Prescriptions, index, index+1)
			} else {
				if err := json.Unmarshal(event.After, &state.Prescriptions[index]); err != nil {
					return state, err
				}
			}
		default:
			return state, fmt.Errorf("unknown event %s", event.Type)
		}
	}
	return state, nil
}
