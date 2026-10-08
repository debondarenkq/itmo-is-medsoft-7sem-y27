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
	AdministrativeSex string
	Comment           string
}

type Command struct {
	Type        string
	DiagnosisID string
	ID          string
	Text        string
}

type RecordChanges struct {
	StaffID         string
	ExpectedVersion *int64
	Commands        []Command
}

type Event struct {
	ID         string          `json:"id"`
	RecordID   string          `json:"record_id"`
	Sequence   int64           `json:"sequence"`
	Version    int64           `json:"version"`
	OccurredAt time.Time       `json:"occurred_at"`
	Actor      Actor           `json:"actor"`
	Type       string          `json:"type"`
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
