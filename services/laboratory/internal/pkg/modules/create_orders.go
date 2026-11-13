package modules

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	api "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/laboratory/api"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/laboratory/internal/models"
	"github.com/google/uuid"
)

const maxOrders = 20

func (m *Module) CreateOrders(ctx context.Context, staffID uuid.UUID, input api.LabOrderInput) (api.FHIRBundle, error) {
	if len(input.TestIds) == 0 || len(input.TestIds) > maxOrders {
		return api.FHIRBundle{}, models.ErrInvalidInput
	}
	patient, err := m.Directories.Patient(ctx, input.PatientID, true)
	if err != nil {
		return api.FHIRBundle{}, err
	}
	doctor, err := m.Directories.Practitioner(ctx, staffID)
	if err != nil {
		return api.FHIRBundle{}, err
	}
	tests, err := m.Store.Tests(ctx, false)
	if err != nil {
		return api.FHIRBundle{}, err
	}
	byID := map[uuid.UUID]models.Test{}
	for _, t := range tests {
		byID[t.ID] = t
	}
	result := []models.Request{}
	for _, id := range input.TestIds {
		t, ok := byID[id]
		if !ok {
			return api.FHIRBundle{}, fmt.Errorf("%w: Анализ отсутствует или отключён", models.ErrInvalidInput)
		}
		name, practitioner := PersonName(patient), PersonName(doctor)
		request := api.FHIRServiceRequest{ID: uuid.New(), ResourceType: api.ServiceRequest, Status: api.Active, Intent: api.Order, AuthoredOn: time.Now().UTC(), Subject: api.FHIRReference{Reference: "Patient/" + input.PatientID.String(), Display: &name}, Requester: api.FHIRReference{Reference: "Practitioner/" + staffID.String(), Display: &practitioner}, Code: api.FHIRCodeableReference{Concept: api.FHIRCodeableConcept{Coding: []api.FHIRCoding{{System: api.HTTPLoincOrg, Code: t.Code, Display: t.Name}}}}}
		if input.Comment != nil && strings.TrimSpace(*input.Comment) != "" {
			if utf8.RuneCountInString(*input.Comment) > maxResultLength {
				return api.FHIRBundle{}, models.ErrInvalidInput
			}
			notes := []api.FHIRAnnotation{{Text: strings.TrimSpace(*input.Comment)}}
			request.Note = &notes
		}
		result = append(result, request)
	}
	if err = m.Store.CreateRequests(ctx, input.PatientID, result); err != nil {
		return api.FHIRBundle{}, err
	}
	return Bundle(result), nil
}

func PersonName(person api.FHIRResource) string {
	if person.Name == nil || len(*person.Name) == 0 {
		return ""
	}
	n := (*person.Name)[0]
	return strings.TrimSpace(n.Family + " " + strings.Join(n.Given, " "))
}

func Bundle[T any](resources []T) api.FHIRBundle {
	result := api.FHIRBundle{ResourceType: api.Bundle, Type: api.Searchset, Total: len(resources), Entry: []api.FHIRBundleEntry{}}
	for _, resource := range resources {
		data, _ := json.Marshal(resource)
		var value map[string]interface{}
		_ = json.Unmarshal(data, &value)
		result.Entry = append(result.Entry, api.FHIRBundleEntry{Resource: value})
	}
	return result
}
