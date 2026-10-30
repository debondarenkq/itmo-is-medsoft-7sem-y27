package main

import (
	"encoding/json"
	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/api/security"
	"os"
)

func main() {
	spec, err := contract.GetSwagger()
	if err != nil {
		panic(err)
	}
	if err = json.NewEncoder(os.Stdout).Encode(spec); err != nil {
		panic(err)
	}
}
