package usecase

import (
	"fmt"
	user "go-template/internal/usecase/user"
	userImpl "go-template/internal/usecase/user/impl"

	auth "go-template/internal/usecase/auth"
	authImpl "go-template/internal/usecase/auth/impl"

	role "go-template/internal/usecase/role"
	roleImpl "go-template/internal/usecase/role/impl"

	"go.uber.org/dig"
)

type Usecase struct {
	dig.In

	User user.Usecase
	Auth auth.Usecase
	Role role.Usecase
}

func Register(container *dig.Container) error {
	if err := container.Provide(userImpl.New); err != nil {
		return fmt.Errorf("[DI] error provide user usecase: %+v", err)
	}
	if err := container.Provide(authImpl.New); err != nil {
		return fmt.Errorf("[DI] error provide auth usecase: %+v", err)
	}
	if err := container.Provide(roleImpl.New); err != nil {
		return fmt.Errorf("[DI] error provide role usecase: %+v", err)
	}
	return nil
}
