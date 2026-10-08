package modules

import (
	"context"
	"strings"
	"time"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/models"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/pkg/platform"
)

type Storage interface {
	ListPatients(context.Context, platform.Page) ([]models.Patient, int, error)
	GetPatient(context.Context, string) (models.Patient, error)
	CreatePatient(context.Context, models.PatientInput) (models.Patient, error)
	UpdatePatient(context.Context, string, models.PatientInput) (models.Patient, error)
	DeletePatient(context.Context, string) error
	CreateRecord(context.Context, string, models.Actor) (models.Record, error)
	GetPatientRecord(context.Context, string) (models.Record, error)
	GetRecord(context.Context, string) (models.Record, error)
	Save(context.Context, string, int64, models.Actor, []models.Command, map[string]models.DiagnosisSnapshot) (models.Record, error)
	History(context.Context, string, platform.Page) ([]models.Event, int, error)
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

func ValidatePatient(v models.PatientInput) (models.PatientInput, error) {
	v.FirstName = strings.TrimSpace(v.FirstName)
	v.LastName = strings.TrimSpace(v.LastName)
	if err := platform.Text("first_name", v.FirstName, 1, 64); err != nil {
		return v, err
	}
	if err := platform.Text("last_name", v.LastName, 1, 64); err != nil {
		return v, err
	}
	if v.MiddleName != nil {
		s := strings.TrimSpace(*v.MiddleName)
		if s == "" {
			v.MiddleName = nil
		} else {
			v.MiddleName = &s
			if err := platform.Text("middle_name", s, 1, 64); err != nil {
				return v, err
			}
		}
	}
	d, err := time.Parse(time.DateOnly, v.BirthDate)
	if err != nil || d.After(time.Now().UTC()) {
		return v, platform.Bad("BIRTH_DATE", "Дата рождения обязательна, формат YYYY-MM-DD, не в будущем")
	}
	if v.AdministrativeSex != "M" && v.AdministrativeSex != "F" && v.AdministrativeSex != "UNKNOWN" {
		return v, platform.Bad("SEX", "Пол должен быть M, F или UNKNOWN")
	}
	if err := platform.Text("comment", v.Comment, 0, 32); err != nil {
		return v, err
	}
	return v, nil
}

func (m *Module) ListPatients(ctx context.Context, p platform.Page) ([]models.Patient, int, error) {
	return m.Storage.ListPatients(ctx, p)
}

func (m *Module) GetPatient(ctx context.Context, id string) (models.Patient, error) {
	return m.Storage.GetPatient(ctx, id)
}

func (m *Module) CreatePatient(ctx context.Context, v models.PatientInput) (models.Patient, error) {
	v, err := ValidatePatient(v)
	if err != nil {
		return models.Patient{}, err
	}
	return m.Storage.CreatePatient(ctx, v)
}

func (m *Module) UpdatePatient(ctx context.Context, id string, v models.PatientInput) (models.Patient, error) {
	v, err := ValidatePatient(v)
	if err != nil {
		return models.Patient{}, err
	}
	return m.Storage.UpdatePatient(ctx, id, v)
}

func (m *Module) DeletePatient(ctx context.Context, id string) error {
	return m.Storage.DeletePatient(ctx, id)
}

func (m *Module) GetPatientRecord(ctx context.Context, id string) (models.Record, error) {
	return m.Storage.GetPatientRecord(ctx, id)
}

func (m *Module) GetRecord(ctx context.Context, id string) (models.Record, error) {
	return m.Storage.GetRecord(ctx, id)
}

func (m *Module) CreateRecord(ctx context.Context, patientID, staffID string) (models.Record, error) {
	staffID, err := platform.ParseID(staffID)
	if err != nil {
		return models.Record{}, platform.Bad("STAFF_REQUIRED", "Укажите медработника, создающего ЭМК")
	}
	if _, err = m.Storage.GetPatient(ctx, patientID); err != nil {
		return models.Record{}, err
	}
	actor, err := m.Staff.GetActor(ctx, staffID)
	if err != nil {
		return models.Record{}, err
	}
	return m.Storage.CreateRecord(ctx, patientID, actor)
}

func (m *Module) Save(ctx context.Context, id string, in models.ChangesInput) (models.Record, error) {
	staffID, err := platform.ParseID(in.StaffID)
	if err != nil {
		return models.Record{}, platform.Bad("STAFF_REQUIRED", "Укажите медработника, изменяющего ЭМК")
	}
	if in.ExpectedVersion == nil || *in.ExpectedVersion < 1 {
		return models.Record{}, platform.Bad("VERSION_REQUIRED", "Требуется expected_version >= 1")
	}
	if err = ValidateCommands(in.Commands); err != nil {
		return models.Record{}, err
	}
	if _, err = m.Storage.GetRecord(ctx, id); err != nil {
		return models.Record{}, err
	}
	actor, err := m.Staff.GetActor(ctx, staffID)
	if err != nil {
		return models.Record{}, err
	}
	diagnoses := map[string]models.DiagnosisSnapshot{}
	for _, c := range in.Commands {
		if c.Type == "add_diagnosis" {
			did, _ := platform.ParseID(c.DiagnosisID)
			if _, ok := diagnoses[did]; ok {
				continue
			}
			d, err := m.Catalog.GetDiagnosis(ctx, did)
			if err != nil {
				return models.Record{}, err
			}
			diagnoses[did] = d
		}
	}
	return m.Storage.Save(ctx, id, *in.ExpectedVersion, actor, in.Commands, diagnoses)
}

func (m *Module) History(ctx context.Context, id string, p platform.Page) ([]models.Event, int, error) {
	if _, err := m.Storage.GetRecord(ctx, id); err != nil {
		return nil, 0, err
	}
	return m.Storage.History(ctx, id, p)
}

func (m *Module) StateAt(ctx context.Context, id string, q models.StateQuery) (models.Record, error) {
	record, err := m.Storage.GetRecord(ctx, id)
	if err != nil {
		return record, err
	}
	if q.Version != nil && (*q.Version < 1 || *q.Version > record.Version) {
		return record, platform.NotFound()
	}
	if q.At != nil && q.At.Before(record.CreatedAt) {
		return record, platform.NotFound()
	}
	events, err := m.Storage.Events(ctx, id, q)
	if err != nil {
		return record, err
	}
	if len(events) == 0 {
		return record, platform.NotFound()
	}
	record.State, err = Replay(events)
	last := events[len(events)-1]
	record.Version = last.Version
	record.UpdatedAt = last.OccurredAt
	return record, err
}
