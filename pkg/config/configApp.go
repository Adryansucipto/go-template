package config

type AppConfig struct {
	BasePath string `yaml:"base-path" default:"sample-project"`
}
