package modules

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/models"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/validation"
	"github.com/google/uuid"
)

const (
	prescriptionTextMinLength = 3
	prescriptionTextMaxLength = 128
)

func (module *Module) Save(ctx context.Context, id string, changes models.RecordChanges) (models.Record, error) {
	staffID, err := validation.ParseID(changes.StaffID)
	if err != nil {
		return models.Record{}, fmt.Errorf("%w: Укажите медработника, изменяющего ЭМК", models.ErrInvalidInput)
	}
	if changes.ExpectedVersion == nil || *changes.ExpectedVersion < 1 {
		return models.Record{}, fmt.Errorf("%w: Требуется expected_version >= 1", models.ErrInvalidInput)
	}
	if err = ValidateCommands(changes.Commands); err != nil {
		return models.Record{}, err
	}
	if _, err = module.Storage.GetRecord(ctx, id); err != nil {
		return models.Record{}, err
	}
	actor, err := module.Staff.GetActor(ctx, staffID)
	if err != nil {
		return models.Record{}, err
	}
	diagnoses := map[string]models.DiagnosisSnapshot{}
	for _, command := range changes.Commands {
		if command.Type == models.CommandAddDiagnosis || command.Type == models.CommandEditDiagnosis {
			diagnosisID, _ := validation.ParseID(command.DiagnosisID)
			if _, ok := diagnoses[diagnosisID]; ok {
				continue
			}
			diagnosis, err := module.Catalog.GetDiagnosis(ctx, diagnosisID)
			if err != nil {
				return models.Record{}, err
			}
			diagnoses[diagnosisID] = diagnosis
		}
	}
	return module.Storage.Save(ctx, id, actor, func(record models.Record) (models.State, []models.Event, error) {
		if record.Version != *changes.ExpectedVersion {
			return models.State{}, nil, models.ErrVersionConflict
		}
		return Apply(record.State, changes.Commands, diagnoses)
	})
}

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
