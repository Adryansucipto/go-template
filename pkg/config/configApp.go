package config

type AppConfig struct {
	Env      string `yaml:"env"`
	Port     string `yaml:"port"`
	BasePath string `yaml:"base-path"`
}
