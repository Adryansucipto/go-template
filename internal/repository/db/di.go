package db

import (
	userRepository "go-template/internal/repository/db/user"
	userRepositoryImpl "go-template/internal/repository/db/user/impl"

	"go.uber.org/dig"
)

type Repository struct {
	dig.In
	UserRepository userRepository.Repository
}

func Register(container *dig.Container) (err error) {
	if err := container.Provide(userRepositoryImpl.New); err != nil {
		return err
	}
	return nil
}
