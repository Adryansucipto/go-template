package service

import (
	"fmt"
	auth "go-template/internal/service/auth"
	authImpl "go-template/internal/service/auth/impl"
	"go-template/internal/service/event/consumer"
	"go-template/internal/service/event/producer"
	producerImpl "go-template/internal/service/event/producer/impl"
	"go-template/internal/service/scheduler"
	user "go-template/internal/service/user"
	userImpl "go-template/internal/service/user/impl"

	"go.uber.org/dig"
)

// Service aggregates business logic services, EDA Event (Producer & Consumer), and Cron Scheduler.
type Service struct {
	dig.In

	User      user.Service
	Auth      auth.Service
	Producer  producer.Service
	Consumer  consumer.Service
	Scheduler scheduler.Service
}

func Register(container *dig.Container) error {
	if err := container.Provide(producerImpl.New); err != nil {
		return fmt.Errorf("[DI] error provide event producer service: %+v", err)
	}
	if err := container.Provide(userImpl.New); err != nil {
		return fmt.Errorf("[DI] error provide user service: %+v", err)
	}
	if err := container.Provide(authImpl.New); err != nil {
		return fmt.Errorf("[DI] error provide auth service: %+v", err)
	}
	return nil
}
