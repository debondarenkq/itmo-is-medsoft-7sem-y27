package api

import (
	"context"
	"encoding/json"
)

func Specification() []byte {
	spec, err := GetSwagger()
	if err != nil {
		panic(err)
	}
	spec.InternalizeRefs(context.Background(), nil)
	data, err := json.Marshal(spec)
	if err != nil {
		panic(err)
	}
	return data
}
