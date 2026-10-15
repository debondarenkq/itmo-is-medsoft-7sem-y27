package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	nethttpmiddleware "github.com/oapi-codegen/nethttp-middleware"
	openapi_types "github.com/oapi-codegen/runtime/types"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/api"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/models"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/pkg/modules"
)

type Deps struct{ Module *modules.Module }
type Implementation struct{ Deps }

var _ contract.StrictServerInterface = (*Implementation)(nil)

func New(deps Deps) *Implementation { return &Implementation{Deps: deps} }

func (implementation *Implementation) Register(router *mux.Router) {
	router.HandleFunc("/openapi.json", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, json.RawMessage(contract.Specification()))
	}).Methods(http.MethodGet)
	spec, err := contract.GetSwagger()
	if err != nil {
		panic(err)
	}
	validator := nethttpmiddleware.OapiRequestValidatorWithOptions(spec, &nethttpmiddleware.Options{
		DoNotValidateServers: true,
		ErrorHandlerWithOpts: func(_ context.Context, err error, w http.ResponseWriter, request *http.Request, options nethttpmiddleware.ErrorHandlerOpts) {
			writeValidationError(w, options.StatusCode, err)
		},
	})
	strict := contract.NewStrictHandlerWithOptions(implementation, nil, contract.StrictHTTPServerOptions{
		RequestErrorHandlerFunc:  writeRequestError,
		ResponseErrorHandlerFunc: writeError,
	})
	contract.HandlerWithOptions(strict, contract.GorillaServerOptions{
		BaseRouter:       router,
		Middlewares:      []contract.MiddlewareFunc{validator},
		ErrorHandlerFunc: writeRequestError,
	})
}

func writeError(w http.ResponseWriter, request *http.Request, err error) {
	status, code := http.StatusInternalServerError, "INTERNAL"
	switch {
	case errors.Is(err, models.ErrVersionConflict):
		status, code = http.StatusConflict, "VERSION_CONFLICT"
	case errors.Is(err, models.ErrUnavailable):
		status, code = http.StatusServiceUnavailable, "DEPENDENCY_UNAVAILABLE"
	case errors.Is(err, models.ErrInvalidInput):
		status, code = http.StatusBadRequest, "INVALID_REQUEST"
	case errors.Is(err, models.ErrNotFound):
		status, code = http.StatusNotFound, "NOT_FOUND"
	case errors.Is(err, models.ErrConflict):
		status, code = http.StatusConflict, "CONFLICT"
	}
	message := err.Error()
	if status == http.StatusInternalServerError {
		slog.Error("request failed", "method", request.Method, "path", request.URL.Path, "error", err)
		message = "Внутренняя ошибка сервиса"
	}
	writeJSON(w, status, contract.Error{Error: contract.ErrorDetail{Code: code, Message: message}})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		slog.Warn("encode response", "error", err)
	}
}

func writeRequestError(w http.ResponseWriter, _ *http.Request, _ error) {
	writeJSON(w, http.StatusBadRequest, contract.Error{Error: contract.ErrorDetail{
		Code: "INVALID_REQUEST", Message: "Некорректный формат запроса",
	}})
}

func writeValidationError(w http.ResponseWriter, status int, err error) {
	detail := contract.ErrorDetail{Code: "INVALID_REQUEST", Message: "Проверьте заполнение полей"}
	var schemaError *openapi3.SchemaError
	if errors.As(err, &schemaError) && schemaError.Schema != nil {
		message := "Некорректное значение"
		switch schemaError.SchemaField {
		case "minLength", "maxLength":
			if maximum := schemaError.Schema.MaxLength; maximum != nil {
				message = fmt.Sprintf("От %d до %d символов", schemaError.Schema.MinLength, *maximum)
			}
		}
		if path := schemaError.JSONPointer(); len(path) > 0 {
			fields := map[string]string{strings.Join(path, "."): message}
			detail.Fields = &fields
		}
	}
	writeJSON(w, status, contract.Error{Error: detail})
}

func pageFromParams(limit, offset *int, search *string) models.Page {
	page := models.Page{Limit: 50}
	if limit != nil {
		page.Limit = *limit
	}
	if offset != nil {
		page.Offset = *offset
	}
	if search != nil {
		page.Search = strings.TrimSpace(*search)
	}
	return page
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func fromPatientInput(value contract.PatientInput) models.PatientData {
	return models.PatientData{
		FirstName:         value.FirstName,
		LastName:          value.LastName,
		MiddleName:        value.MiddleName,
		BirthDate:         value.BirthDate.Time.Format(time.DateOnly),
		AdministrativeSex: string(value.AdministrativeSex),
		Comment:           stringValue(value.Comment),
	}
}

func toPatient(value models.Patient) (contract.Patient, error) {
	id, err := uuid.Parse(value.ID)
	if err != nil {
		return contract.Patient{}, fmt.Errorf("map patient id: %w", err)
	}
	birthDate, err := time.Parse(time.DateOnly, value.BirthDate)
	if err != nil {
		return contract.Patient{}, fmt.Errorf("map patient birth date: %w", err)
	}
	return contract.Patient{
		ID:                id,
		FirstName:         value.FirstName,
		LastName:          value.LastName,
		MiddleName:        value.MiddleName,
		BirthDate:         openapi_types.Date{Time: birthDate},
		AdministrativeSex: contract.AdministrativeSex(value.AdministrativeSex),
		Comment:           value.Comment,
		HasRecord:         value.HasRecord,
		CreatedAt:         value.CreatedAt,
		UpdatedAt:         value.UpdatedAt,
	}, nil
}

func toRecord(value models.Record) (contract.Record, error) {
	id, err := uuid.Parse(value.ID)
	if err != nil {
		return contract.Record{}, fmt.Errorf("map record id: %w", err)
	}
	patientID, err := uuid.Parse(value.PatientID)
	if err != nil {
		return contract.Record{}, fmt.Errorf("map record patient id: %w", err)
	}
	diagnoses := make([]contract.RecordDiagnosis, len(value.State.Diagnoses))
	for index, diagnosis := range value.State.Diagnoses {
		diagnoses[index], err = toRecordDiagnosis(diagnosis)
		if err != nil {
			return contract.Record{}, fmt.Errorf("map record diagnosis %d: %w", index, err)
		}
	}
	prescriptions := make([]contract.Prescription, len(value.State.Prescriptions))
	for index, prescription := range value.State.Prescriptions {
		prescriptions[index], err = toPrescription(prescription)
		if err != nil {
			return contract.Record{}, fmt.Errorf("map record prescription %d: %w", index, err)
		}
	}
	return contract.Record{
		ID:        id,
		PatientID: patientID,
		Version:   value.Version,
		CreatedAt: value.CreatedAt,
		UpdatedAt: value.UpdatedAt,
		State: contract.State{
			Diagnoses:     diagnoses,
			Prescriptions: prescriptions,
		},
	}, nil
}

func toRecordDiagnosis(value models.RecordDiagnosis) (contract.RecordDiagnosis, error) {
	id, err := uuid.Parse(value.ID)
	if err != nil {
		return contract.RecordDiagnosis{}, fmt.Errorf("id: %w", err)
	}
	diagnosisID, err := uuid.Parse(value.DiagnosisID)
	if err != nil {
		return contract.RecordDiagnosis{}, fmt.Errorf("diagnosis id: %w", err)
	}
	return contract.RecordDiagnosis{
		ID:          id,
		DiagnosisID: diagnosisID,
		Code:        value.Code,
		Name:        value.Name,
	}, nil
}

func toPrescription(value models.Prescription) (contract.Prescription, error) {
	id, err := uuid.Parse(value.ID)
	if err != nil {
		return contract.Prescription{}, fmt.Errorf("id: %w", err)
	}
	return contract.Prescription{
		ID:     id,
		Text:   value.Text,
		Status: contract.PrescriptionStatus(value.Status),
	}, nil
}
