package api

import "embed"

//go:embed openapi.json
var files embed.FS

func Specification() []byte {
	data, _ := files.ReadFile("openapi.json")
	return data
}
