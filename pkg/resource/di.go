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

	if err := container.Provide(injection.NewKafka); err != nil {
		return fmt.Errorf("[DI] cannot initialize Kafka: %+v", err)
	}

	return nil
}
