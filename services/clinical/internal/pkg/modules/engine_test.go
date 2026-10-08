package modules

import (
	"reflect"
	"strings"
	"testing"

	"github.com/google/uuid"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/models"
)

func TestHistoryRestoresEveryChange(t *testing.T) {
	id := uuid.NewString()
	diagnoses := map[string]models.DiagnosisSnapshot{id: {ID: id, Code: "J00.00", Name: "Назофарингит"}}
	state := models.EmptyState()
	all := []models.Event{}
	apply := func(commands ...models.Command) []models.Event {
		t.Helper()
		next, events, err := Apply(state, commands, diagnoses)
		if err != nil {
			t.Fatal(err)
		}
		all = append(all, events...)
		restored, err := Replay(all)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(restored, next) {
			t.Fatalf("replay differs: %#v != %#v", restored, next)
		}
		state = next
		return events
	}
	apply(models.Command{Type: "add_diagnosis", DiagnosisID: id}, models.Command{Type: "add_prescription", Text: "Пить больше воды"})
	originalID := state.Prescriptions[0].ID
	events := apply(models.Command{Type: "edit_prescription", ID: originalID, Text: "Пить воду после еды"})
	if len(events) != 2 || events[0].Type != "prescription_removed" || events[1].Type != "prescription_added" {
		t.Fatalf("edit events: %#v", events)
	}
	if state.Prescriptions[0].ID == originalID {
		t.Fatal("editing reused prescription ID")
	}
	apply(models.Command{Type: "cancel_prescription", ID: state.Prescriptions[0].ID})
	if state.Prescriptions[0].Status != "cancelled" {
		t.Fatal("cancel status lost")
	}
	apply(models.Command{Type: "remove_diagnosis", ID: state.Diagnoses[0].ID}, models.Command{Type: "remove_prescription", ID: state.Prescriptions[0].ID})
	if len(state.Diagnoses) != 0 || len(state.Prescriptions) != 0 {
		t.Fatal("removed entities remain")
	}
}

func TestFailedBatchDoesNotMutateState(t *testing.T) {
	state := models.State{Diagnoses: []models.RecordDiagnosis{}, Prescriptions: []models.Prescription{{ID: uuid.NewString(), Text: "Назначение", Status: "active"}}}
	before := state.Prescriptions[0]
	_, events, err := Apply(state, []models.Command{{Type: "cancel_prescription", ID: before.ID}, {Type: "remove_diagnosis", ID: uuid.NewString()}}, nil)
	if err == nil || events != nil {
		t.Fatal("expected failed batch")
	}
	if state.Prescriptions[0] != before {
		t.Fatal("input state mutated")
	}
}

func TestValidationCountsUnicodeCharacters(t *testing.T) {
	if err := ValidateCommands([]models.Command{{Type: "add_prescription", Text: strings.Repeat("я", 128)}}); err != nil {
		t.Fatal(err)
	}
	if err := ValidateCommands([]models.Command{{Type: "add_prescription", Text: strings.Repeat("я", 129)}}); err == nil {
		t.Fatal("overlong text accepted")
	}
	if err := ValidateCommands([]models.Command{{Type: "add_prescription", Text: "   "}}); err == nil {
		t.Fatal("blank text accepted")
	}
}

func TestPatientValidation(t *testing.T) {
	valid := models.PatientInput{FirstName: "Иван", LastName: "Иванов", BirthDate: "2000-02-29", AdministrativeSex: "UNKNOWN"}
	if _, err := ValidatePatient(valid); err != nil {
		t.Fatal(err)
	}
	for _, date := range []string{"2001-02-29", "3000-01-01", "", "01.01.2000"} {
		v := valid
		v.BirthDate = date
		if _, err := ValidatePatient(v); err == nil {
			t.Errorf("accepted birth date %q", date)
		}
	}
}
