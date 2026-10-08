package api

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/api"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/models"
)

func stringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
func optionalID(value *uuid.UUID) string {
	if value == nil {
		return ""
	}
	return value.String()
}
func fromPatientInput(value contract.PatientInput) models.PatientData {
	return models.PatientData{FirstName: value.FirstName, LastName: value.LastName, MiddleName: value.MiddleName, BirthDate: value.BirthDate.Time.Format(time.DateOnly), AdministrativeSex: string(value.AdministrativeSex), Comment: stringValue(value.Comment)}
}
func fromChangesInput(value contract.ChangesInput) models.RecordChanges {
	commands := make([]models.Command, len(value.Commands))
	for index, command := range value.Commands {
		commands[index] = models.Command{Type: string(command.Type), ID: optionalID(command.ID), DiagnosisID: optionalID(command.DiagnosisID), Text: stringValue(command.Text)}
	}
	return models.RecordChanges{StaffID: value.StaffID.String(), ExpectedVersion: &value.ExpectedVersion, Commands: commands}
}
func toPatient(value models.Patient) contract.Patient {
	birthDate, err := time.Parse(time.DateOnly, value.BirthDate)
	if err != nil {
		panic(err)
	}
	return contract.Patient{ID: uuid.MustParse(value.ID), FirstName: value.FirstName, LastName: value.LastName, MiddleName: value.MiddleName, BirthDate: openapi_types.Date{Time: birthDate}, AdministrativeSex: contract.AdministrativeSex(value.AdministrativeSex), Comment: value.Comment, HasRecord: value.HasRecord, CreatedAt: value.CreatedAt, UpdatedAt: value.UpdatedAt}
}
func toRecord(value models.Record) contract.Record {
	diagnoses := make([]contract.RecordDiagnosis, len(value.State.Diagnoses))
	for index, diagnosis := range value.State.Diagnoses {
		diagnoses[index] = contract.RecordDiagnosis{ID: uuid.MustParse(diagnosis.ID), DiagnosisID: uuid.MustParse(diagnosis.DiagnosisID), Code: diagnosis.Code, Name: diagnosis.Name}
	}
	prescriptions := make([]contract.Prescription, len(value.State.Prescriptions))
	for index, prescription := range value.State.Prescriptions {
		prescriptions[index] = contract.Prescription{ID: uuid.MustParse(prescription.ID), Text: prescription.Text, Status: contract.PrescriptionStatus(prescription.Status)}
	}
	return contract.Record{ID: uuid.MustParse(value.ID), PatientID: uuid.MustParse(value.PatientID), Version: value.Version, CreatedAt: value.CreatedAt, UpdatedAt: value.UpdatedAt, State: contract.State{Diagnoses: diagnoses, Prescriptions: prescriptions}}
}
func snapshot(value json.RawMessage) *map[string]interface{} {
	var result *map[string]interface{}
	if err := json.Unmarshal(value, &result); err != nil {
		panic(err)
	}
	return result
}
func toEvent(value models.Event) contract.Event {
	return contract.Event{ID: uuid.MustParse(value.ID), RecordID: uuid.MustParse(value.RecordID), Sequence: value.Sequence, Version: value.Version, OccurredAt: value.OccurredAt, Actor: contract.Actor{ID: uuid.MustParse(value.Actor.ID), FirstName: value.Actor.FirstName, LastName: value.Actor.LastName, Position: value.Actor.Position}, Type: contract.EventType(value.Type), EntityID: uuid.MustParse(value.EntityID), Before: snapshot(value.Before), After: snapshot(value.After)}
}
