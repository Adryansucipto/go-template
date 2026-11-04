package resource

import (
	"fmt"
	"go-template/pkg/resource/injection"

	"go.uber.org/dig"
)

func Register(container *dig.Container) error {
	if err := container.Provide(injection.NewDatabase); err != nil {
		return fmt.Errorf("[DI] cannot initialize Database: %+v", err)
	}

	if err := container.Provide(injection.NewRedis); err != nil {
		return fmt.Errorf("[DI] cannot initialize Redis: %+v", err)
	}

	if err := container.Provide(injection.NewKafka); err != nil {
		return fmt.Errorf("[DI] cannot initialize Kafka: %+v", err)
	}

	if err := container.Provide(injection.NewHTTPClient); err != nil {
		return fmt.Errorf("[DI] cannot initialize HTTPClient: %+v", err)
	}

	if err := container.Provide(injection.NewGRPCClient); err != nil {
		return fmt.Errorf("[DI] cannot initialize GRPCClient: %+v", err)
	}

	return nil
}
