package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/api"
)

func main() {
	output := flag.String("output", "openapi.json", "Generated JSON documentation path")
	flag.Parse()
	var spec any
	if err := json.Unmarshal(contract.Specification(), &spec); err != nil {
		log.Fatal(err)
	}
	data, err := json.MarshalIndent(spec, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	if err = os.WriteFile(*output, append(data, '\n'), 0644); err != nil {
		log.Fatal(err)
	}
}
