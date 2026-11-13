package models

import contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/laboratory/api"

type Test = contract.LabTest
type Request = contract.FHIRServiceRequest
type Observation = contract.FHIRObservation

const FHIRMediaType = "application/fhir+json"

func Outcome(code contract.FHIROperationOutcomeIssueCode, message string) contract.FHIROperationOutcome {
	result := contract.FHIROperationOutcome{ResourceType: contract.OperationOutcome}
	result.Issue = append(result.Issue, struct {
		Code        contract.FHIROperationOutcomeIssueCode     `json:"code"`
		Diagnostics string                                     `json:"diagnostics"`
		Severity    contract.FHIROperationOutcomeIssueSeverity `json:"severity"`
	}{Code: code, Diagnostics: message, Severity: contract.Error})
	return result
}
