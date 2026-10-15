package modules

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/models"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/validation"
)

var (
	errPatientHasRecord = fmt.Errorf("%w: Нельзя удалить пациента, для которого создана ЭМК", models.ErrConflict)
	errRecordExists     = fmt.Errorf("%w: У пациента уже есть ЭМК", models.ErrConflict)
)

type Storage interface {
	ListPatients(context.Context, models.Page) ([]models.Patient, int, error)
	GetPatient(context.Context, string) (models.Patient, error)
	CreatePatient(context.Context, models.PatientData) (models.Patient, error)
	UpdatePatient(context.Context, string, models.PatientData) (models.Patient, error)
	DeletePatient(context.Context, string) error
	CreateRecord(context.Context, string, models.Actor) (models.Record, error)
	GetPatientRecord(context.Context, string) (models.Record, error)
	GetRecord(context.Context, string) (models.Record, error)
	Save(context.Context, string, models.Actor, func(models.Record) (models.State, []models.Event, error)) (models.Record, error)
	History(context.Context, string, models.Page) ([]models.Event, int, error)
	Events(context.Context, string, models.StateQuery) ([]models.Event, error)
}

type Staff interface {
	GetActor(context.Context, string) (models.Actor, error)
}
type Catalog interface {
	GetDiagnosis(context.Context, string) (models.DiagnosisSnapshot, error)
}
type Deps struct {
	Storage Storage
	Staff   Staff
	Catalog Catalog
}
type Module struct {
	Deps
}

func New(deps Deps) *Module {
	return &Module{Deps: deps}
}

func ValidatePatient(patientData models.PatientData) (models.PatientData, error) {
	patientData.FirstName = strings.TrimSpace(patientData.FirstName)
	patientData.LastName = strings.TrimSpace(patientData.LastName)
	if err := validation.Text("first_name", patientData.FirstName, 1, 64); err != nil {
		return patientData, err
	}
	if err := validation.Text("last_name", patientData.LastName, 1, 64); err != nil {
		return patientData, err
	}
	if patientData.MiddleName != nil {
		middleName := strings.TrimSpace(*patientData.MiddleName)
		if middleName == "" {
			patientData.MiddleName = nil
		} else {
			patientData.MiddleName = &middleName
			if err := validation.Text("middle_name", middleName, 1, 64); err != nil {
				return patientData, err
			}
		}
	}
	birthDate, err := time.Parse(time.DateOnly, patientData.BirthDate)
	if err != nil || birthDate.After(time.Now().UTC()) {
		return patientData, fmt.Errorf("%w: Дата рождения обязательна, формат YYYY-MM-DD, не в будущем", models.ErrInvalidInput)
	}
	if patientData.AdministrativeSex != "M" && patientData.AdministrativeSex != "F" && patientData.AdministrativeSex != "UNKNOWN" {
		return patientData, fmt.Errorf("%w: Пол должен быть M, F или UNKNOWN", models.ErrInvalidInput)
	}
	if err := validation.Text("comment", patientData.Comment, 0, 32); err != nil {
		return patientData, err
	}
	return patientData, nil
}

func (module *Module) ListPatients(ctx context.Context, page models.Page) ([]models.Patient, int, error) {
	return module.Storage.ListPatients(ctx, page)
}

func (module *Module) GetPatient(ctx context.Context, id string) (models.Patient, error) {
	return module.Storage.GetPatient(ctx, id)
}

func (module *Module) CreatePatient(ctx context.Context, patientData models.PatientData) (models.Patient, error) {
	patientData, err := ValidatePatient(patientData)
	if err != nil {
		return models.Patient{}, err
	}
	return module.Storage.CreatePatient(ctx, patientData)
}

func (module *Module) UpdatePatient(ctx context.Context, id string, patientData models.PatientData) (models.Patient, error) {
	patientData, err := ValidatePatient(patientData)
	if err != nil {
		return models.Patient{}, err
	}
	return module.Storage.UpdatePatient(ctx, id, patientData)
}

func (module *Module) DeletePatient(ctx context.Context, id string) error {
	patient, err := module.Storage.GetPatient(ctx, id)
	if err != nil {
		return err
	}
	if patient.HasRecord {
		return errPatientHasRecord
	}
	err = module.Storage.DeletePatient(ctx, id)
	if errors.Is(err, models.ErrReferenced) {
		return errPatientHasRecord
	}
	return err
}

func (module *Module) GetPatientRecord(ctx context.Context, id string) (models.Record, error) {
	return module.Storage.GetPatientRecord(ctx, id)
}

func (module *Module) GetRecord(ctx context.Context, id string) (models.Record, error) {
	return module.Storage.GetRecord(ctx, id)
}

func (module *Module) CreateRecord(ctx context.Context, patientID, staffID string) (models.Record, error) {
	staffID, err := validation.ParseID(staffID)
	if err != nil {
		return models.Record{}, fmt.Errorf("%w: Укажите медработника, создающего ЭМК", models.ErrInvalidInput)
	}
	patient, err := module.Storage.GetPatient(ctx, patientID)
	if err != nil {
		return models.Record{}, err
	}
	if patient.HasRecord {
		return models.Record{}, errRecordExists
	}
	actor, err := module.Staff.GetActor(ctx, staffID)
	if err != nil {
		return models.Record{}, err
	}
	record, err := module.Storage.CreateRecord(ctx, patientID, actor)
	if errors.Is(err, models.ErrAlreadyExists) {
		return models.Record{}, errRecordExists
	}
	return record, err
}

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
		if command.Type == "add_diagnosis" {
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

func (module *Module) History(ctx context.Context, id string, page models.Page) ([]models.Event, int, error) {
	if _, err := module.Storage.GetRecord(ctx, id); err != nil {
		return nil, 0, err
	}
	return module.Storage.History(ctx, id, page)
}

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
