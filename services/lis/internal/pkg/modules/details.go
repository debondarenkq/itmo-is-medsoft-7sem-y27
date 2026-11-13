package modules

import (
	"context"
	"strings"

	api "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/lis/api"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/lis/internal/models"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/lis/internal/pkg/adapters"
	his "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/lis/internal/pkg/adapters/laboratoryapi"
	"github.com/google/uuid"
)

func (m *Module) Details(ctx context.Context, id uuid.UUID) (api.LISRequestDetails, error) {
	history, err := m.Store.ByRequest(ctx, id)
	if err != nil {
		return api.LISRequestDetails{}, err
	}
	if history != nil {
		return api.LISRequestDetails{Request: history.Request, Patient: history.Patient, Practitioner: history.Practitioner, History: history}, nil
	}
	client, _, err := m.client(ctx)
	if err != nil {
		return api.LISRequestDetails{}, err
	}
	response, err := client.GetFHIRServiceRequestWithResponse(ctx, id)
	if err != nil {
		return api.LISRequestDetails{}, err
	}
	if response.ApplicationFhirJSON200 == nil {
		return api.LISRequestDetails{}, adapters.ResponseError(response.StatusCode(), response.Body)
	}
	request := response.ApplicationFhirJSON200
	patientID, err := referenceID(request.Subject.Reference, "Patient")
	if err != nil {
		return api.LISRequestDetails{}, err
	}
	staffID, err := referenceID(request.Requester.Reference, "Practitioner")
	if err != nil {
		return api.LISRequestDetails{}, err
	}
	patient, err := client.GetFHIRPatientWithResponse(ctx, patientID)
	if err != nil {
		return api.LISRequestDetails{}, err
	}
	if patient.ApplicationFhirJSON200 == nil {
		return api.LISRequestDetails{}, adapters.ResponseError(patient.StatusCode(), patient.Body)
	}
	practitioner, err := client.GetFHIRPractitionerWithResponse(ctx, staffID)
	if err != nil {
		return api.LISRequestDetails{}, err
	}
	if practitioner.ApplicationFhirJSON200 == nil {
		return api.LISRequestDetails{}, adapters.ResponseError(practitioner.StatusCode(), practitioner.Body)
	}
	details := map[string]any{"request": request, "patient": patient.ApplicationFhirJSON200, "practitioner": practitioner.ApplicationFhirJSON200}
	if request.Status == his.Completed {
		patient := patientID.String()
		results, err := client.SearchFHIRObservationWithResponse(ctx, &his.SearchFHIRObservationParams{Patient: &patient})
		if err != nil {
			return api.LISRequestDetails{}, err
		}
		if results.ApplicationFhirJSON200 == nil {
			return api.LISRequestDetails{}, adapters.ResponseError(results.StatusCode(), results.Body)
		}
		for _, entry := range results.ApplicationFhirJSON200.Entry {
			observation, err := adapters.Convert[api.FHIRObservation](entry.Resource)
			if err != nil {
				return api.LISRequestDetails{}, err
			}
			for _, reference := range observation.BasedOn {
				if reference.Reference == "ServiceRequest/"+id.String() {
					details["result"] = observation
				}
			}
		}
		if details["result"] == nil {
			return api.LISRequestDetails{}, models.ErrNotFound
		}
	}
	return adapters.Convert[api.LISRequestDetails](details)
}

func referenceID(reference, kind string) (uuid.UUID, error) {
	if !strings.HasPrefix(reference, kind+"/") {
		return uuid.Nil, models.ErrInvalid
	}
	id, err := uuid.Parse(strings.TrimPrefix(reference, kind+"/"))
	return id, err
}
