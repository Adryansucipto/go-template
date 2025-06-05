package di

import (
	"go-template/internal/repository"
	"go-template/internal/usecase"
	"go-template/pkg/config"
	"go-template/pkg/resource"
	"sync"

	"go.uber.org/dig"
)

var (
	container *dig.Container
	once      sync.Once
)

// DI ini untuk menyimpan semua Usecase, Repository yang sudah di Register
// Sedangkan Controller tidak perlu karena tidak perlu DI, maka tidak perlu di Register

func Container() (*dig.Container, error) {
	var outer error
	once.Do(func() {
		container = dig.New()

		if err := config.Register(container); err != nil {
			outer = err
			return
		}

		if err := resource.Register(container); err != nil {
			outer = err
			return
		}

		if err := usecase.Register(container); err != nil {
			outer = err
			return
		}

		if err := repository.Register(container); err != nil {
			outer = err
			return
		}
	})

	if outer != nil {
		return nil, outer
	}
	return container, nil
}
