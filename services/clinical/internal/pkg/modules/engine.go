package modules

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/google/uuid"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/models"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/validation"
)

const (
	prescriptionTextMinLength = 3
	prescriptionTextMaxLength = 128
)

func ValidateCommands(commands []models.Command) error {
	if len(commands) == 0 || len(commands) > 100 {
		return fmt.Errorf("%w: Требуется от 1 до 100 команд", models.ErrInvalidInput)
	}
	for _, command := range commands {
		switch command.Type {
		case models.CommandAddDiagnosis, models.CommandEditDiagnosis:
			if _, err := validation.ParseID(command.DiagnosisID); err != nil {
				return err
			}
			if command.Type == models.CommandEditDiagnosis {
				if _, err := validation.ParseID(command.ID); err != nil {
					return err
				}
			} else if command.ID != "" {
				return fmt.Errorf("%w: Для add_diagnosis допускается только diagnosis_id", models.ErrInvalidInput)
			}
			if command.Text != "" {
				return fmt.Errorf("%w: Для диагноза нельзя передавать текст назначения", models.ErrInvalidInput)
			}
		case models.CommandRemoveDiagnosis, models.CommandCancelPrescription, models.CommandRemovePrescription:
			if _, err := validation.ParseID(command.ID); err != nil {
				return err
			}
			if command.Text != "" || command.DiagnosisID != "" {
				return fmt.Errorf("%w: Для этой команды допускается только id", models.ErrInvalidInput)
			}
		case models.CommandAddPrescription:
			if err := validation.Text("text", strings.TrimSpace(command.Text), prescriptionTextMinLength, prescriptionTextMaxLength); err != nil {
				return err
			}
			if command.ID != "" || command.DiagnosisID != "" {
				return fmt.Errorf("%w: Для add_prescription допускается только text", models.ErrInvalidInput)
			}
		case models.CommandEditPrescription:
			if _, err := validation.ParseID(command.ID); err != nil {
				return err
			}
			if err := validation.Text("text", strings.TrimSpace(command.Text), prescriptionTextMinLength, prescriptionTextMaxLength); err != nil {
				return err
			}
			if command.DiagnosisID != "" {
				return fmt.Errorf("%w: Для edit_prescription допускаются id и text", models.ErrInvalidInput)
			}
		default:
			return fmt.Errorf("%w: Неизвестный тип команды: %s", models.ErrInvalidInput, command.Type)
		}
	}
	return nil
}

// Apply is pure with respect to input state. A failed batch never changes it.
func Apply(state models.State, commands []models.Command, diagnoses map[string]models.DiagnosisSnapshot) (models.State, []models.Event, error) {
	if err := ValidateCommands(commands); err != nil {
		return state, nil, err
	}
	nextState := models.EmptyState()
	nextState.Diagnoses = append(nextState.Diagnoses, state.Diagnoses...)
	nextState.Prescriptions = append(nextState.Prescriptions, state.Prescriptions...)
	events := []models.Event{}
	emit := func(kind models.EventType, id string, before, after any) {
		beforeData, _ := json.Marshal(before)
		afterData, _ := json.Marshal(after)
		events = append(events, models.Event{ID: uuid.NewString(), Type: kind, EntityID: id, Before: beforeData, After: afterData})
	}
	for _, command := range commands {
		command.Text = strings.TrimSpace(command.Text)
		if command.ID != "" {
			command.ID, _ = validation.ParseID(command.ID)
		}
		if command.DiagnosisID != "" {
			command.DiagnosisID, _ = validation.ParseID(command.DiagnosisID)
		}
		switch command.Type {
		case models.CommandAddDiagnosis:
			if slices.ContainsFunc(nextState.Diagnoses, func(diagnosis models.RecordDiagnosis) bool { return diagnosis.DiagnosisID == command.DiagnosisID }) {
				return state, nil, fmt.Errorf("%w: Диагноз уже содержится в ЭМК", models.ErrConflict)
			}
			diagnosis, ok := diagnoses[command.DiagnosisID]
			if !ok {
				return state, nil, fmt.Errorf("%w: Диагноз не проверен в справочнике", models.ErrInvalidInput)
			}
			recordDiagnosis := models.RecordDiagnosis{ID: uuid.NewString(), DiagnosisID: diagnosis.ID, Code: diagnosis.Code, Name: diagnosis.Name}
			nextState.Diagnoses = append(nextState.Diagnoses, recordDiagnosis)
			emit(models.EventDiagnosisAdded, recordDiagnosis.ID, nil, recordDiagnosis)
		case models.CommandEditDiagnosis:
			index := slices.IndexFunc(nextState.Diagnoses, func(diagnosis models.RecordDiagnosis) bool { return diagnosis.ID == command.ID })
			if index < 0 {
				return state, nil, models.ErrNotFound
			}
			if slices.ContainsFunc(nextState.Diagnoses, func(diagnosis models.RecordDiagnosis) bool {
				return diagnosis.ID != command.ID && diagnosis.DiagnosisID == command.DiagnosisID
			}) {
				return state, nil, fmt.Errorf("%w: Диагноз уже содержится в ЭМК", models.ErrConflict)
			}
			diagnosis, ok := diagnoses[command.DiagnosisID]
			if !ok {
				return state, nil, fmt.Errorf("%w: Диагноз не проверен в справочнике", models.ErrInvalidInput)
			}
			previousDiagnosis := nextState.Diagnoses[index]
			nextState.Diagnoses[index] = models.RecordDiagnosis{ID: command.ID, DiagnosisID: diagnosis.ID, Code: diagnosis.Code, Name: diagnosis.Name}
			emit(models.EventDiagnosisUpdated, command.ID, previousDiagnosis, nextState.Diagnoses[index])
		case models.CommandRemoveDiagnosis:
			index := slices.IndexFunc(nextState.Diagnoses, func(diagnosis models.RecordDiagnosis) bool { return diagnosis.ID == command.ID })
			if index < 0 {
				return state, nil, models.ErrNotFound
			}
			emit(models.EventDiagnosisRemoved, command.ID, nextState.Diagnoses[index], nil)
			nextState.Diagnoses = slices.Delete(nextState.Diagnoses, index, index+1)
		case models.CommandAddPrescription:
			prescription := models.Prescription{ID: uuid.NewString(), Text: command.Text, Status: models.PrescriptionActive}
			nextState.Prescriptions = append(nextState.Prescriptions, prescription)
			emit(models.EventPrescriptionAdded, prescription.ID, nil, prescription)
		case models.CommandCancelPrescription, models.CommandEditPrescription, models.CommandRemovePrescription:
			index := slices.IndexFunc(nextState.Prescriptions, func(prescription models.Prescription) bool { return prescription.ID == command.ID })
			if index < 0 {
				return state, nil, models.ErrNotFound
			}
			previousPrescription := nextState.Prescriptions[index]
			if command.Type != models.CommandRemovePrescription && previousPrescription.Status != models.PrescriptionActive {
				return state, nil, fmt.Errorf("%w: Отменённое назначение нельзя повторно отменить или редактировать", models.ErrConflict)
			}
			if command.Type == models.CommandCancelPrescription {
				nextState.Prescriptions[index].Status = models.PrescriptionCancelled
				emit(models.EventPrescriptionCancelled, command.ID, previousPrescription, nextState.Prescriptions[index])
			} else if command.Type == models.CommandEditPrescription {
				nextState.Prescriptions[index].Text = command.Text
				emit(models.EventPrescriptionUpdated, command.ID, previousPrescription, nextState.Prescriptions[index])
			} else {
				emit(models.EventPrescriptionRemoved, command.ID, previousPrescription, nil)
				nextState.Prescriptions = slices.Delete(nextState.Prescriptions, index, index+1)
			}
		}
	}
	return nextState, events, nil
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
