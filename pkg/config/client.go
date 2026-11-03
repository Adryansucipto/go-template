package config

type GRPCServerConfig struct {
	Port string `yaml:"port"`
}

type OutboundClientConfig struct {
	TimeoutSeconds int                              `yaml:"timeout_seconds"`
	Services       map[string]ExternalServiceConfig `yaml:"services"`
}

type ExternalServiceConfig struct {
	BaseURL    string `yaml:"base_url,omitempty"`
	GRPCAddr   string `yaml:"grpc_addr,omitempty"`
	APIKey     string `yaml:"api_key,omitempty"`
	TimeoutSec int    `yaml:"timeout_seconds,omitempty"`
}

