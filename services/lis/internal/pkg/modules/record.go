package modules

import (
	"context"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	api "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/lis/api"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/lis/internal/models"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/lis/internal/pkg/adapters"
	his "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/lis/internal/pkg/adapters/laboratoryapi"
	"github.com/google/uuid"
)

const maxResultLength = 512

func (m *Module) Record(ctx context.Context, input api.LISResultInput) (models.History, error) {
	value := strings.TrimSpace(input.Value)
	if value == "" || utf8.RuneCountInString(value) > maxResultLength {
		return models.History{}, models.ErrInvalid
	}
	details, err := m.Details(ctx, input.RequestID)
	if err != nil {
		return models.History{}, err
	}
	if details.History != nil || details.Request.Status != api.Active {
		return models.History{}, models.ErrConflict
	}
	now := time.Now().UTC()
	h := models.History{ID: uuid.New(), Request: details.Request, Patient: details.Patient, Practitioner: details.Practitioner, Result: api.FHIRObservation{ID: uuid.New(), ResourceType: api.Observation, Status: api.Final, Code: details.Request.Code.Concept, Subject: details.Request.Subject, BasedOn: []api.FHIRReference{{Reference: "ServiceRequest/" + input.RequestID.String()}}, EffectiveDateTime: now, Issued: now, ValueString: value}}
	client, _, err := m.client(ctx)
	if err != nil {
		return h, err
	}
	if err = m.Store.Record(ctx, h); err != nil {
		return h, err
	}
	body, err := adapters.Convert[his.FHIRObservation](h.Result)
	if err != nil {
		return h, err
	}
	response, err := client.CreateFHIRObservationWithApplicationFhirPlusJSONBodyWithResponse(ctx, body)
	if err != nil {
		return h, fmt.Errorf("Результат сохранён в LIS; доставка в HIS не подтверждена: %w", err)
	}
	if response.ApplicationFhirJSON201 == nil {
		return h, fmt.Errorf("Результат сохранён в LIS; доставка в HIS не подтверждена: %w", adapters.ResponseError(response.StatusCode(), response.Body))
	}
	if err = m.Store.Delivered(ctx, h.ID); err != nil {
		return h, err
	}
	h.Delivered = true
	return h, nil
}
