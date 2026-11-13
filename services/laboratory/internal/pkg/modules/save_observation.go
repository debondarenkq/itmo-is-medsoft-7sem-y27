package modules

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"
	"unicode/utf8"

	api "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/laboratory/api"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/laboratory/internal/models"
	"github.com/google/uuid"
)

const maxResultLength = 512

func (m *Module) SaveObservation(ctx context.Context, result api.FHIRObservation) (api.FHIRObservation, error) {
	if result.ResourceType != api.Observation || result.ID == uuid.Nil || result.Status != api.Final || len(result.BasedOn) != 1 || strings.TrimSpace(result.ValueString) == "" || utf8.RuneCountInString(result.ValueString) > maxResultLength || result.EffectiveDateTime.IsZero() || result.Issued.IsZero() {
		return result, models.ErrInvalidInput
	}
	requestID, err := ReferenceID(result.BasedOn[0].Reference, "ServiceRequest")
	if err != nil {
		return result, err
	}
	patientID, err := ReferenceID(result.Subject.Reference, "Patient")
	if err != nil {
		return result, err
	}
	err = m.Store.SaveObservation(ctx, requestID, patientID, result, func(request models.Request) error {
		if request.Status == api.Completed {
			return models.ErrConflict
		}
		if request.Subject.Reference != result.Subject.Reference || !reflect.DeepEqual(request.Code.Concept, result.Code) {
			return fmt.Errorf("%w: Результат не соответствует пациенту или коду назначения", models.ErrInvalidInput)
		}
		if result.EffectiveDateTime.Before(request.AuthoredOn) || result.EffectiveDateTime.After(time.Now().Add(time.Minute)) || result.Issued.Before(result.EffectiveDateTime) {
			return fmt.Errorf("%w: Некорректная дата результата", models.ErrInvalidInput)
		}
		return nil
	})
	return result, err
}

func ReferenceID(reference, kind string) (uuid.UUID, error) {
	prefix := kind + "/"
	if !strings.HasPrefix(reference, prefix) {
		return uuid.Nil, models.ErrInvalidInput
	}
	id, err := uuid.Parse(strings.TrimPrefix(reference, prefix))
	if err != nil || id == uuid.Nil {
		return uuid.Nil, models.ErrInvalidInput
	}
	return id, nil
}
