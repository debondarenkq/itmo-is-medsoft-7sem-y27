package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/models"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/migrations"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"testing"
)

func TestPatientSearchByRecordAndRecordFilter(t *testing.T) {
	repository := testRepository(t)
	ctx := context.Background()
	patient := createPatient(t, repository)
	createPatient(t, repository)
	record, err := repository.CreateRecord(ctx, patient.ID, actor())
	if err != nil {
		t.Fatal(err)
	}
	items, total, err := repository.ListPatients(ctx, models.Page{Limit: 50, Search: record.ID[:8]})
	if err != nil || total != 1 || len(items) != 1 || items[0].RecordID == nil || *items[0].RecordID != record.ID {
		t.Fatalf("record search: %+v total=%d err=%v", items, total, err)
	}
	if record.Number < 1 {
		t.Fatal("record number was not allocated")
	}
	items, total, err = repository.ListPatients(ctx, models.Page{Limit: 50, Search: fmt.Sprintf("ЭМК-%06d", record.Number)})
	if err != nil || total != 1 || len(items) != 1 || items[0].RecordNumber == nil || *items[0].RecordNumber != record.Number {
		t.Fatalf("human record number search: %+v %d %v", items, total, err)
	}
	withoutRecord := false
	items, total, err = repository.ListPatients(ctx, models.Page{Limit: 50, HasRecord: &withoutRecord})
	if err != nil || total != 1 || len(items) != 1 || items[0].HasRecord || items[0].RecordID != nil {
		t.Fatalf("without record: %+v total=%d err=%v", items, total, err)
	}
	withRecord := true
	items, total, err = repository.ListPatients(ctx, models.Page{Limit: 50, HasRecord: &withRecord})
	if err != nil || total != 1 || len(items) != 1 || items[0].ID != patient.ID {
		t.Fatalf("with record: %+v total=%d err=%v", items, total, err)
	}
}

func TestRecordNumbersPreserveExistingCards(t *testing.T) {
	repository := testRepository(t)
	ctx := context.Background()
	patient := createPatient(t, repository)
	database := stdlib.OpenDB(*repository.pool.Config().ConnConfig)
	defer database.Close()
	provider, err := goose.NewProvider(goose.DialectPostgres, database, migrations.FS)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = provider.DownTo(ctx, 2); err != nil {
		t.Fatal(err)
	}
	id := uuid.NewString()
	state, err := json.Marshal(models.EmptyState())
	if err != nil {
		t.Fatal(err)
	}
	if _, err = repository.pool.Exec(ctx, "INSERT INTO records(id,patient_id,state) VALUES($1,$2,$3)", id, patient.ID, state); err != nil {
		t.Fatal(err)
	}
	if _, err = provider.Up(ctx); err != nil {
		t.Fatal(err)
	}
	record, err := repository.GetRecord(ctx, id)
	if err != nil || record.ID != id || record.PatientID != patient.ID || record.Number < 1 || record.Version != 1 {
		t.Fatalf("existing card after migration: %+v %v", record, err)
	}
	if len(record.State.Diagnoses) != 0 || len(record.State.Prescriptions) != 0 {
		t.Fatal("migration changed existing medical data")
	}
}
