package storage

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/models"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/pkg/modules"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/pkg/platform"
)

func createPatient(t *testing.T, r *Repository) models.Patient {
	t.Helper()
	v, err := r.CreatePatient(context.Background(), models.PatientData{FirstName: "Иван", LastName: "Иванов", BirthDate: "2000-01-01", AdministrativeSex: "M"})
	if err != nil {
		t.Fatal(err)
	}
	return v
}
func actor() models.Actor {
	return models.Actor{ID: uuid.NewString(), FirstName: "Анна", LastName: "Петрова", Position: "Терапевт"}
}

func TestRecordLifecycleAndHistoricalState(t *testing.T) {
	r := testRepository(t)
	ctx := context.Background()
	patient := createPatient(t, r)
	a := actor()
	record, err := r.CreateRecord(ctx, patient.ID, a)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = r.CreateRecord(ctx, patient.ID, a); err == nil {
		t.Fatal("duplicate record accepted")
	}
	if err = r.DeletePatient(ctx, patient.ID); err == nil {
		t.Fatal("patient with record deleted")
	}
	d := models.DiagnosisSnapshot{ID: uuid.NewString(), Code: "J00.00", Name: "Назофарингит"}
	record, err = r.Save(ctx, record.ID, 1, a, []models.Command{{Type: "add_diagnosis", DiagnosisID: d.ID}, {Type: "add_prescription", Text: "Пить больше воды"}}, map[string]models.DiagnosisSnapshot{d.ID: d})
	if err != nil {
		t.Fatal(err)
	}
	added := record
	originalID := record.State.Prescriptions[0].ID
	record, err = r.Save(ctx, record.ID, 2, a, []models.Command{{Type: "edit_prescription", ID: originalID, Text: "Пить воду после еды"}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	edited := record
	if record.State.Prescriptions[0].ID == originalID {
		t.Fatal("ID must change when editing")
	}
	record, err = r.Save(ctx, record.ID, 3, a, []models.Command{{Type: "cancel_prescription", ID: record.State.Prescriptions[0].ID}, {Type: "remove_diagnosis", ID: record.State.Diagnoses[0].ID}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	events, total, err := r.History(ctx, record.ID, platform.Page{Limit: 200})
	if err != nil {
		t.Fatal(err)
	}
	if total != 7 || len(events) != 7 {
		t.Fatalf("events=%d/%d", len(events), total)
	}
	if !events[3].OccurredAt.Equal(events[4].OccurredAt) || events[3].Version != events[4].Version {
		t.Fatal("edit not atomic in history")
	}
	for _, e := range events {
		if e.Actor != a {
			t.Fatal("author snapshot changed")
		}
	}
	restored, err := modules.Replay(events)
	if err != nil || !reflect.DeepEqual(restored, record.State) {
		t.Fatalf("current replay mismatch: %v", err)
	}
	for _, v := range []models.Record{added, edited} {
		at, err := r.Events(ctx, record.ID, models.StateQuery{At: &v.UpdatedAt})
		if err != nil {
			t.Fatal(err)
		}
		state, err := modules.Replay(at)
		if err != nil || !reflect.DeepEqual(state, v.State) {
			t.Fatalf("historical replay mismatch: %v", err)
		}
	}
	// Failure after a valid command leaves both current state and history untouched.
	_, err = r.Save(ctx, record.ID, 4, a, []models.Command{{Type: "add_prescription", Text: "Новое назначение"}, {Type: "remove_diagnosis", ID: uuid.NewString()}}, nil)
	if err == nil {
		t.Fatal("invalid batch accepted")
	}
	current, err := r.GetRecord(ctx, record.ID)
	if err != nil || current.Version != 4 || !reflect.DeepEqual(current.State, record.State) {
		t.Fatal("failed batch changed record")
	}
	_, total, err = r.History(ctx, record.ID, platform.Page{Limit: 200})
	if err != nil || total != 7 {
		t.Fatal("failed batch changed history")
	}
	// Event snapshots contain complete before/after data, not mutable references.
	var before models.Prescription
	if err = json.Unmarshal(events[3].Before, &before); err != nil || before.ID != originalID || before.Text != "Пить больше воды" {
		t.Fatal("incomplete edit snapshot")
	}
}

func TestConcurrentSavesUseVersion(t *testing.T) {
	r := testRepository(t)
	ctx := context.Background()
	record, err := r.CreateRecord(ctx, createPatient(t, r).ID, actor())
	if err != nil {
		t.Fatal(err)
	}
	start := make(chan struct{})
	results := make(chan error, 2)
	for range 2 {
		go func() {
			<-start
			_, err := r.Save(ctx, record.ID, 1, actor(), []models.Command{{Type: "add_prescription", Text: "Назначение"}}, nil)
			results <- err
		}()
	}
	close(start)
	successes, conflicts := 0, 0
	for range 2 {
		err := <-results
		var api *platform.Error
		if err == nil {
			successes++
		} else if errors.As(err, &api) && api.Code == "VERSION_CONFLICT" {
			conflicts++
		} else {
			t.Fatal(err)
		}
	}
	if successes != 1 || conflicts != 1 {
		t.Fatalf("successes=%d conflicts=%d", successes, conflicts)
	}
}

func TestConcurrentCreateRecordAndDeletePatient(t *testing.T) {
	r := testRepository(t)
	ctx := context.Background()
	for range 20 {
		patient := createPatient(t, r)
		start := make(chan struct{})
		var createErr, deleteErr error
		var wg sync.WaitGroup
		wg.Add(2)
		go func() { defer wg.Done(); <-start; _, createErr = r.CreateRecord(ctx, patient.ID, actor()) }()
		go func() { defer wg.Done(); <-start; deleteErr = r.DeletePatient(ctx, patient.ID) }()
		close(start)
		wg.Wait()
		if createErr == nil {
			var api *platform.Error
			if !errors.As(deleteErr, &api) || api.Code != "PATIENT_HAS_RECORD" {
				t.Fatalf("create success, delete=%v", deleteErr)
			}
		} else if !errors.Is(createErr, pgx.ErrNoRows) || deleteErr != nil {
			t.Fatalf("create=%v delete=%v", createErr, deleteErr)
		}
	}
	var orphans int
	if err := r.pool.QueryRow(ctx, "SELECT count(*) FROM records r LEFT JOIN patients p ON p.id=r.patient_id WHERE p.id IS NULL").Scan(&orphans); err != nil || orphans != 0 {
		t.Fatalf("orphans=%d err=%v", orphans, err)
	}
}
