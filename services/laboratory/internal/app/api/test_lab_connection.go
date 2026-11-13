package api

import (
	"context"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/laboratory/api"
)

func (i *Implementation) TestLabConnection(context.Context, contract.TestLabConnectionRequestObject) (contract.TestLabConnectionResponseObject, error) {
	return contract.TestLabConnection200JSONResponse{Online: true, FhirVersion: "5.0.0"}, nil
}
