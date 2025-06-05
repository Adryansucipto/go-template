package config

type Config struct {
	AppConfig AppConfig `yaml:"app"`
	Database  Database  `yaml:"database"`
}

func New() (Config, error) {
	return Config{}, nil
}
