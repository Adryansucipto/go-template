package di

import (
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
	})

	if outer != nil {
		return nil, outer
	}
	return container, nil
}
