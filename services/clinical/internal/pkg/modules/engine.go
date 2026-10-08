package modules

import (
	"encoding/json"
	"fmt"
	"slices"

	"github.com/google/uuid"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/models"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/pkg/platform"
)

func ValidateCommands(commands []models.Command) error {
	if len(commands) == 0 || len(commands) > 100 {
		return platform.Bad("COMMANDS", "Требуется от 1 до 100 команд")
	}
	for _, c := range commands {
		switch c.Type {
		case "add_diagnosis":
			if _, err := platform.ParseID(c.DiagnosisID); err != nil {
				return err
			}
			if c.ID != "" || c.Text != "" {
				return platform.Bad("COMMAND", "Для add_diagnosis допускается только diagnosis_id")
			}
		case "remove_diagnosis", "cancel_prescription", "remove_prescription":
			if _, err := platform.ParseID(c.ID); err != nil {
				return err
			}
			if c.Text != "" || c.DiagnosisID != "" {
				return platform.Bad("COMMAND", "Для этой команды допускается только id")
			}
		case "add_prescription":
			if err := platform.Text("text", c.Text, 3, 128); err != nil {
				return err
			}
			if c.ID != "" || c.DiagnosisID != "" {
				return platform.Bad("COMMAND", "Для add_prescription допускается только text")
			}
		case "edit_prescription":
			if _, err := platform.ParseID(c.ID); err != nil {
				return err
			}
			if err := platform.Text("text", c.Text, 3, 128); err != nil {
				return err
			}
			if c.DiagnosisID != "" {
				return platform.Bad("COMMAND", "Для edit_prescription допускаются id и text")
			}
		default:
			return platform.Bad("COMMAND", "Неизвестный тип команды: "+c.Type)
		}
	}
	return nil
}

// Apply is pure with respect to input state. A failed batch never changes it.
func Apply(state models.State, commands []models.Command, diagnoses map[string]models.DiagnosisSnapshot) (models.State, []models.Event, error) {
	if err := ValidateCommands(commands); err != nil {
		return state, nil, err
	}
	next := models.EmptyState()
	next.Diagnoses = append(next.Diagnoses, state.Diagnoses...)
	next.Prescriptions = append(next.Prescriptions, state.Prescriptions...)
	events := []models.Event{}
	emit := func(kind, id string, before, after any) {
		b, _ := json.Marshal(before)
		a, _ := json.Marshal(after)
		events = append(events, models.Event{ID: uuid.NewString(), Type: kind, EntityID: id, Before: b, After: a})
	}
	for _, c := range commands {
		// UUID spellings are normalized even for upper-case or alternative valid input.
		if c.ID != "" {
			c.ID, _ = platform.ParseID(c.ID)
		}
		if c.DiagnosisID != "" {
			c.DiagnosisID, _ = platform.ParseID(c.DiagnosisID)
		}
		switch c.Type {
		case "add_diagnosis":
			if slices.ContainsFunc(next.Diagnoses, func(d models.RecordDiagnosis) bool { return d.DiagnosisID == c.DiagnosisID }) {
				return state, nil, platform.Conflict("DIAGNOSIS_EXISTS", "Диагноз уже содержится в ЭМК")
			}
			d, ok := diagnoses[c.DiagnosisID]
			if !ok {
				return state, nil, platform.Bad("DIAGNOSIS", "Диагноз не проверен в справочнике")
			}
			v := models.RecordDiagnosis{ID: uuid.NewString(), DiagnosisID: d.ID, Code: d.Code, Name: d.Name}
			next.Diagnoses = append(next.Diagnoses, v)
			emit("diagnosis_added", v.ID, nil, v)
		case "remove_diagnosis":
			idx := slices.IndexFunc(next.Diagnoses, func(v models.RecordDiagnosis) bool { return v.ID == c.ID })
			if idx < 0 {
				return state, nil, platform.NotFound()
			}
			emit("diagnosis_removed", c.ID, next.Diagnoses[idx], nil)
			next.Diagnoses = slices.Delete(next.Diagnoses, idx, idx+1)
		case "add_prescription":
			v := models.Prescription{ID: uuid.NewString(), Text: c.Text, Status: "active"}
			next.Prescriptions = append(next.Prescriptions, v)
			emit("prescription_added", v.ID, nil, v)
		case "cancel_prescription", "edit_prescription", "remove_prescription":
			idx := slices.IndexFunc(next.Prescriptions, func(v models.Prescription) bool { return v.ID == c.ID })
			if idx < 0 {
				return state, nil, platform.NotFound()
			}
			old := next.Prescriptions[idx]
			if c.Type != "remove_prescription" && old.Status != "active" {
				return state, nil, platform.Conflict("PRESCRIPTION_CANCELLED", "Отменённое назначение нельзя повторно отменить или редактировать")
			}
			if c.Type == "cancel_prescription" {
				next.Prescriptions[idx].Status = "cancelled"
				emit("prescription_cancelled", c.ID, old, next.Prescriptions[idx])
			} else {
				emit("prescription_removed", c.ID, old, nil)
				next.Prescriptions = slices.Delete(next.Prescriptions, idx, idx+1)
				if c.Type == "edit_prescription" {
					v := models.Prescription{ID: uuid.NewString(), Text: c.Text, Status: "active"}
					next.Prescriptions = append(next.Prescriptions, v)
					emit("prescription_added", v.ID, nil, v)
				}
			}
		}
	}
	return next, events, nil
}

func Replay(events []models.Event) (models.State, error) {
	state := models.EmptyState()
	for _, e := range events {
		switch e.Type {
		case "record_created":
		case "diagnosis_added":
			var v models.RecordDiagnosis
			if err := json.Unmarshal(e.After, &v); err != nil {
				return state, err
			}
			state.Diagnoses = append(state.Diagnoses, v)
		case "diagnosis_removed":
			idx := slices.IndexFunc(state.Diagnoses, func(v models.RecordDiagnosis) bool { return v.ID == e.EntityID })
			if idx < 0 {
				return state, fmt.Errorf("invalid diagnosis event %s", e.ID)
			}
			state.Diagnoses = slices.Delete(state.Diagnoses, idx, idx+1)
		case "prescription_added":
			var v models.Prescription
			if err := json.Unmarshal(e.After, &v); err != nil {
				return state, err
			}
			state.Prescriptions = append(state.Prescriptions, v)
		case "prescription_cancelled", "prescription_removed":
			idx := slices.IndexFunc(state.Prescriptions, func(v models.Prescription) bool { return v.ID == e.EntityID })
			if idx < 0 {
				return state, fmt.Errorf("invalid prescription event %s", e.ID)
			}
			if e.Type == "prescription_removed" {
				state.Prescriptions = slices.Delete(state.Prescriptions, idx, idx+1)
			} else {
				if err := json.Unmarshal(e.After, &state.Prescriptions[idx]); err != nil {
					return state, err
				}
			}
		default:
			return state, fmt.Errorf("unknown event %s", e.Type)
		}
	}
	return state, nil
}
