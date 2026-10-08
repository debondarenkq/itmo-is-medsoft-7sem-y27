package models

import (
	"encoding/json"
	"time"
)

type PatientInput struct {
	FirstName         string  `json:"first_name"`
	LastName          string  `json:"last_name"`
	MiddleName        *string `json:"middle_name"`
	BirthDate         string  `json:"birth_date"`
	AdministrativeSex string  `json:"administrative_sex"`
	Comment           string  `json:"comment"`
}

type CreateRecordInput struct {
	StaffID string `json:"staff_id"`
}

type Command struct {
	Type        string `json:"type"`
	DiagnosisID string `json:"diagnosis_id,omitempty"`
	ID          string `json:"id,omitempty"`
	Text        string `json:"text,omitempty"`
}

type ChangesInput struct {
	StaffID         string    `json:"staff_id"`
	ExpectedVersion *int64    `json:"expected_version"`
	Commands        []Command `json:"commands"`
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
