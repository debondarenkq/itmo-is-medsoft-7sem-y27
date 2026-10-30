package config

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func SigningKey() (ed25519.PrivateKey, error) {
	path := os.Getenv("SIGNING_KEY_FILE")
	if path == "" {
		return nil, fmt.Errorf("SIGNING_KEY_FILE is required")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("invalid signing key PEM")
	}
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	ed, ok := key.(ed25519.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("Ed25519 signing key required")
	}
	return ed, nil
}
