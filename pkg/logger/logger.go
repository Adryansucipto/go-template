package logger

import (
	"fmt"
	"go-template/pkg/config"

	"go.uber.org/dig"
	"go.uber.org/zap"
)

func New(cfg config.Config) (*zap.Logger, error) {
	if cfg.AppConfig.Env == "prod" {
		return zap.NewProduction()
	}
	return zap.NewDevelopment()
}

func Register(container *dig.Container) error {
	if err := container.Provide(New); err != nil {
		return fmt.Errorf("[DI] cannot initialize logger: %+v", err)
	}
	return nil
}

