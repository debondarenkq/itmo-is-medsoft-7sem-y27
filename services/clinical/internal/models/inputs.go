package models

import (
	"encoding/json"
	"time"
)

type PatientData struct {
	FirstName         string
	LastName          string
	MiddleName        *string
	BirthDate         string
	AdministrativeSex AdministrativeSex
	Comment           string
}

type CommandType string

const (
	CommandAddDiagnosis       CommandType = "add_diagnosis"
	CommandEditDiagnosis      CommandType = "edit_diagnosis"
	CommandRemoveDiagnosis    CommandType = "remove_diagnosis"
	CommandAddPrescription    CommandType = "add_prescription"
	CommandEditPrescription   CommandType = "edit_prescription"
	CommandCancelPrescription CommandType = "cancel_prescription"
	CommandRemovePrescription CommandType = "remove_prescription"
)

type Command struct {
	Type        CommandType
	DiagnosisID string
	ID          string
	Text        string
}

type RecordChanges struct {
	StaffID         string
	ExpectedVersion *int64
	Commands        []Command
}

type EventType string

const (
	EventRecordCreated         EventType = "record_created"
	EventDiagnosisAdded        EventType = "diagnosis_added"
	EventDiagnosisUpdated      EventType = "diagnosis_updated"
	EventDiagnosisRemoved      EventType = "diagnosis_removed"
	EventPrescriptionAdded     EventType = "prescription_added"
	EventPrescriptionUpdated   EventType = "prescription_updated"
	EventPrescriptionCancelled EventType = "prescription_cancelled"
	EventPrescriptionRemoved   EventType = "prescription_removed"
)

type Event struct {
	ID         string          `json:"id"`
	RecordID   string          `json:"record_id"`
	Sequence   int64           `json:"sequence"`
	Version    int64           `json:"version"`
	OccurredAt time.Time       `json:"occurred_at"`
	Actor      Actor           `json:"actor"`
	Type       EventType       `json:"type"`
	EntityID   string          `json:"entity_id"`
	Before     json.RawMessage `json:"before"`
	After      json.RawMessage `json:"after"`
}

type DiagnosisSnapshot struct {
	ID   string `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type StateQuery struct {
	At      *time.Time
	Version *int64
}
