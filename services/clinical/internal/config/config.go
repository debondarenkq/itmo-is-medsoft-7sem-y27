package config

import "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/pkg/platform"

type Config struct{ StaffURL, CatalogURL string }

func New() Config {
	return Config{StaffURL: platform.Env("STAFF_URL", "http://localhost:8081"), CatalogURL: platform.Env("CATALOG_URL", "http://localhost:8082")}
}
