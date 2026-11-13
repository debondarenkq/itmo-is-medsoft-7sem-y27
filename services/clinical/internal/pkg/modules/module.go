package modules

import (
	"context"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/models"
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
