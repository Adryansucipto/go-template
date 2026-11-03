package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Config struct {
	AppConfig AppConfig       `yaml:"app"`
	Database  Database        `yaml:"database"`
	JWT       JwtCnfg         `yaml:"jwt"`
	Kafka     KafkaConfig     `yaml:"kafka"`
	Scheduler SchedulerConfig `yaml:"scheduler"`
}

func New() (Config, error) {
	var cfg Config

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}

	configPath := filepath.Join("configs", fmt.Sprintf("%s.yaml", env))
	data, err := os.ReadFile(configPath)
	if err != nil {
		return cfg, fmt.Errorf("failed to read config file (%s): %w", configPath, err)
	}

	expanded := os.ExpandEnv(string(data))
	if err := yaml.Unmarshal([]byte(expanded), &cfg); err != nil {
		return cfg, fmt.Errorf("failed to parse YAML: %w", err)
	}

	return cfg, nil
}
