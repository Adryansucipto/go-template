package db

import (
	userRepository "go-template/internal/repository/db/user"
	userRepositoryImpl "go-template/internal/repository/db/user/impl"

	authRepository "go-template/internal/repository/db/auth"
	authRepositoryImpl "go-template/internal/repository/db/auth/impl"

	roleRepository "go-template/internal/repository/db/role"
	roleRepositoryImpl "go-template/internal/repository/db/role/impl"

	"go.uber.org/dig"
)

type Repository struct {
	dig.In
	UserRepository userRepository.Repository
	AuthRepository authRepository.Repository
	RoleRepository roleRepository.Repository
}

func Register(container *dig.Container) (err error) {
	if err := container.Provide(userRepositoryImpl.New); err != nil {
		return err
	}
	if err := container.Provide(authRepositoryImpl.New); err != nil {
		return err
	}
	if err := container.Provide(roleRepositoryImpl.New); err != nil {
		return err
	}
	return nil
}
