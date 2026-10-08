package models

import "time"

type Patient struct {
	ID                string    `json:"id"`
	FirstName         string    `json:"first_name"`
	LastName          string    `json:"last_name"`
	MiddleName        *string   `json:"middle_name"`
	BirthDate         string    `json:"birth_date"`
	AdministrativeSex string    `json:"administrative_sex"`
	Comment           string    `json:"comment"`
	HasRecord         bool      `json:"has_record"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// Snapshots deliberately omit mutable service metadata.
type Actor struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Position  string `json:"position"`
}

type RecordDiagnosis struct {
	ID          string `json:"id"`
	DiagnosisID string `json:"diagnosis_id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
}

type Prescription struct {
	ID     string `json:"id"`
	Text   string `json:"text"`
	Status string `json:"status"`
}

type State struct {
	Diagnoses     []RecordDiagnosis `json:"diagnoses"`
	Prescriptions []Prescription    `json:"prescriptions"`
}

func EmptyState() State {
	return State{Diagnoses: []RecordDiagnosis{}, Prescriptions: []Prescription{}}
}

type Record struct {
	ID        string    `json:"id"`
	PatientID string    `json:"patient_id"`
	Version   int64     `json:"version"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	State     State     `json:"state"`
}
