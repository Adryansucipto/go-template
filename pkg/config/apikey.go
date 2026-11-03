package config

import (
	"crypto/subtle"
	"strings"
)

type APIKeyConfig struct {
	Enabled bool              `yaml:"enabled"`
	Header  string            `yaml:"header"`
	Key     string            `yaml:"key"`
	Clients map[string]string `yaml:"clients,omitempty"`
}

func (a *APIKeyConfig) HeaderName() string {
	if strings.TrimSpace(a.Header) == "" {
		return "X-API-Key"
	}
	return a.Header
}

// Validate checks the incoming key against the primary key or named client keys using constant-time comparison.
// Returns (clientName, isValid).
func (a *APIKeyConfig) Validate(incomingKey string) (string, bool) {
	if !a.Enabled {
		return "disabled", true
	}

	trimmed := strings.TrimSpace(incomingKey)
	if trimmed == "" {
		return "", false
	}

	if a.Key != "" && subtle.ConstantTimeCompare([]byte(trimmed), []byte(a.Key)) == 1 {
		return "default", true
	}

	for clientName, secretKey := range a.Clients {
		if secretKey != "" && subtle.ConstantTimeCompare([]byte(trimmed), []byte(secretKey)) == 1 {
			return clientName, true
		}
	}

	return "", false
}

