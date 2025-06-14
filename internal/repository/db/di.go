package db

import (
	userRepository "go-template/internal/repository/db/user"
	userRepositoryImpl "go-template/internal/repository/db/user/impl"

	authRepository "go-template/internal/repository/db/auth"
	authRepositoryImpl "go-template/internal/repository/db/auth/impl"

	"go.uber.org/dig"
)

type Repository struct {
	dig.In
	UserRepository userRepository.Repository
	AuthRepository authRepository.Repository
}

func Register(container *dig.Container) (err error) {
	if err := container.Provide(userRepositoryImpl.New); err != nil {
		return err
	}
	if err := container.Provide(authRepositoryImpl.New); err != nil {
		return err
	}
	return nil
}
