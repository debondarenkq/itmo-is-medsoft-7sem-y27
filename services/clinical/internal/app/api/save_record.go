package api

import (
	"context"

	"github.com/google/uuid"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/api"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/models"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/validation"
)

func (implementation *Implementation) SaveRecord(ctx context.Context, request contract.SaveRecordRequestObject) (contract.SaveRecordResponseObject, error) {
	if err := validation.RequireID(request.ID); err != nil {
		return nil, err
	}
	record, err := implementation.Module.Save(ctx, request.ID.String(), fromChangesInput(*request.Body))
	if err != nil {
		return nil, err
	}
	response, err := toRecord(record)
	if err != nil {
		return nil, err
	}
	return contract.SaveRecord200JSONResponse(response), nil
}

func fromChangesInput(value contract.ChangesInput) models.RecordChanges {
	commands := make([]models.Command, len(value.Commands))
	for index, command := range value.Commands {
		commands[index] = models.Command{
			Type:        string(command.Type),
			ID:          optionalID(command.ID),
			DiagnosisID: optionalID(command.DiagnosisID),
			Text:        stringValue(command.Text),
		}
	}
	return models.RecordChanges{
		StaffID:         value.StaffID.String(),
		ExpectedVersion: &value.ExpectedVersion,
		Commands:        commands,
	}
}

func optionalID(value *uuid.UUID) string {
	if value == nil {
		return ""
	}
	return value.String()
}
