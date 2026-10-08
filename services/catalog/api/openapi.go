package api

import "encoding/json"

// Specification is JSON derived from the generated, embedded source contract.
func Specification() []byte {
	spec, err := GetSwagger()
	if err != nil {
		panic(err)
	}
	data, err := json.Marshal(spec)
	if err != nil {
		panic(err)
	}
	return data
}
