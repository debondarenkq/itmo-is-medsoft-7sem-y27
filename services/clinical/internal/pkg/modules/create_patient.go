package modules

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/models"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/validation"
)

const (
	patientNameMinLength    = 1
	patientNameMaxLength    = 64
	patientCommentMinLength = 0
	patientCommentMaxLength = 32
)

func (module *Module) CreatePatient(ctx context.Context, patientData models.PatientData) (models.Patient, error) {
	patientData, err := ValidatePatient(patientData)
	if err != nil {
		return models.Patient{}, err
	}
	return module.Storage.CreatePatient(ctx, patientData)
}

func ValidatePatient(patientData models.PatientData) (models.PatientData, error) {
	patientData.FirstName = strings.TrimSpace(patientData.FirstName)
	patientData.LastName = strings.TrimSpace(patientData.LastName)
	if err := validation.Text("first_name", patientData.FirstName, patientNameMinLength, patientNameMaxLength); err != nil {
		return patientData, err
	}
	if err := validation.Text("last_name", patientData.LastName, patientNameMinLength, patientNameMaxLength); err != nil {
		return patientData, err
	}
	if patientData.MiddleName != nil {
		middleName := strings.TrimSpace(*patientData.MiddleName)
		if middleName == "" {
			patientData.MiddleName = nil
		} else {
			patientData.MiddleName = &middleName
			if err := validation.Text("middle_name", middleName, patientNameMinLength, patientNameMaxLength); err != nil {
				return patientData, err
			}
		}
	}
	birthDate, err := time.Parse(time.DateOnly, patientData.BirthDate)
	if err != nil || birthDate.After(time.Now().UTC()) {
		return patientData, fmt.Errorf("%w: Дата рождения обязательна, формат YYYY-MM-DD, не в будущем", models.ErrInvalidInput)
	}
	switch patientData.AdministrativeSex {
	case models.SexMale, models.SexFemale, models.SexUnknown:
	default:
		return patientData, fmt.Errorf("%w: Пол должен быть M, F или UNKNOWN", models.ErrInvalidInput)
	}
	if err := validation.Text("comment", patientData.Comment, patientCommentMinLength, patientCommentMaxLength); err != nil {
		return patientData, err
	}
	return patientData, nil
}
