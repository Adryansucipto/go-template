package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	AppConfig AppConfig `yaml:"app"`
	Database  Database  `yaml:"database"`
}

func New() (Config, error) {
	var cfg Config

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}
	filepath := fmt.Sprintf("./variables/%s.yaml", env)
	data, err := os.ReadFile(filepath)
	if err != nil {
		return cfg, fmt.Errorf("failed to read config file: %w", err)
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, fmt.Errorf("failed to parse YAML: %w", err)
	}

	return cfg, nil
}
