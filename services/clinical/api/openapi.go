package api

import "encoding/json"

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
