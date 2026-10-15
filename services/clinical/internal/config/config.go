package config

import "os"

type Config struct {
	GatewayURL string
}

func New() Config {
	return Config{GatewayURL: os.Getenv("GATEWAY_URL")}
}
