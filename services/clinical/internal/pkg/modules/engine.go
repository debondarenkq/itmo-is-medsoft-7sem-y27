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
		case "add_diagnosis":
			if _, err := validation.ParseID(command.DiagnosisID); err != nil {
				return err
			}
			if command.ID != "" || command.Text != "" {
				return fmt.Errorf("%w: Для add_diagnosis допускается только diagnosis_id", models.ErrInvalidInput)
			}
		case "remove_diagnosis", "cancel_prescription", "remove_prescription":
			if _, err := validation.ParseID(command.ID); err != nil {
				return err
			}
			if command.Text != "" || command.DiagnosisID != "" {
				return fmt.Errorf("%w: Для этой команды допускается только id", models.ErrInvalidInput)
			}
		case "add_prescription":
			if err := validation.Text("text", strings.TrimSpace(command.Text), prescriptionTextMinLength, prescriptionTextMaxLength); err != nil {
				return err
			}
			if command.ID != "" || command.DiagnosisID != "" {
				return fmt.Errorf("%w: Для add_prescription допускается только text", models.ErrInvalidInput)
			}
		case "edit_prescription":
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
			return fmt.Errorf("%w: %s", models.ErrInvalidInput, "Неизвестный тип команды: "+command.Type)
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
	emit := func(kind, id string, before, after any) {
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
		case "add_diagnosis":
			if slices.ContainsFunc(nextState.Diagnoses, func(diagnosis models.RecordDiagnosis) bool { return diagnosis.DiagnosisID == command.DiagnosisID }) {
				return state, nil, fmt.Errorf("%w: Диагноз уже содержится в ЭМК", models.ErrConflict)
			}
			diagnosis, ok := diagnoses[command.DiagnosisID]
			if !ok {
				return state, nil, fmt.Errorf("%w: Диагноз не проверен в справочнике", models.ErrInvalidInput)
			}
			recordDiagnosis := models.RecordDiagnosis{ID: uuid.NewString(), DiagnosisID: diagnosis.ID, Code: diagnosis.Code, Name: diagnosis.Name}
			nextState.Diagnoses = append(nextState.Diagnoses, recordDiagnosis)
			emit("diagnosis_added", recordDiagnosis.ID, nil, recordDiagnosis)
		case "remove_diagnosis":
			index := slices.IndexFunc(nextState.Diagnoses, func(diagnosis models.RecordDiagnosis) bool { return diagnosis.ID == command.ID })
			if index < 0 {
				return state, nil, models.ErrNotFound
			}
			emit("diagnosis_removed", command.ID, nextState.Diagnoses[index], nil)
			nextState.Diagnoses = slices.Delete(nextState.Diagnoses, index, index+1)
		case "add_prescription":
			prescription := models.Prescription{ID: uuid.NewString(), Text: command.Text, Status: "active"}
			nextState.Prescriptions = append(nextState.Prescriptions, prescription)
			emit("prescription_added", prescription.ID, nil, prescription)
		case "cancel_prescription", "edit_prescription", "remove_prescription":
			index := slices.IndexFunc(nextState.Prescriptions, func(prescription models.Prescription) bool { return prescription.ID == command.ID })
			if index < 0 {
				return state, nil, models.ErrNotFound
			}
			previousPrescription := nextState.Prescriptions[index]
			if command.Type != "remove_prescription" && previousPrescription.Status != "active" {
				return state, nil, fmt.Errorf("%w: Отменённое назначение нельзя повторно отменить или редактировать", models.ErrConflict)
			}
			if command.Type == "cancel_prescription" {
				nextState.Prescriptions[index].Status = "cancelled"
				emit("prescription_cancelled", command.ID, previousPrescription, nextState.Prescriptions[index])
			} else {
				emit("prescription_removed", command.ID, previousPrescription, nil)
				nextState.Prescriptions = slices.Delete(nextState.Prescriptions, index, index+1)
				if command.Type == "edit_prescription" {
					prescription := models.Prescription{ID: uuid.NewString(), Text: command.Text, Status: "active"}
					nextState.Prescriptions = append(nextState.Prescriptions, prescription)
					emit("prescription_added", prescription.ID, nil, prescription)
				}
			}
		}
	}
	return nextState, events, nil
}

func Replay(events []models.Event) (models.State, error) {
	state := models.EmptyState()
	for _, event := range events {
		switch event.Type {
		case "record_created":
		case "diagnosis_added":
			var diagnosis models.RecordDiagnosis
			if err := json.Unmarshal(event.After, &diagnosis); err != nil {
				return state, err
			}
			state.Diagnoses = append(state.Diagnoses, diagnosis)
		case "diagnosis_removed":
			index := slices.IndexFunc(state.Diagnoses, func(diagnosis models.RecordDiagnosis) bool { return diagnosis.ID == event.EntityID })
			if index < 0 {
				return state, fmt.Errorf("invalid diagnosis event %s", event.ID)
			}
			state.Diagnoses = slices.Delete(state.Diagnoses, index, index+1)
		case "prescription_added":
			var prescription models.Prescription
			if err := json.Unmarshal(event.After, &prescription); err != nil {
				return state, err
			}
			state.Prescriptions = append(state.Prescriptions, prescription)
		case "prescription_cancelled", "prescription_removed":
			index := slices.IndexFunc(state.Prescriptions, func(prescription models.Prescription) bool { return prescription.ID == event.EntityID })
			if index < 0 {
				return state, fmt.Errorf("invalid prescription event %s", event.ID)
			}
			if event.Type == "prescription_removed" {
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
