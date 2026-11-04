package cache

import (
	userCache "go-template/internal/repository/cache/user"
	userCacheImpl "go-template/internal/repository/cache/user/impl"

	"go.uber.org/dig"
)

type Repository struct {
	dig.In

	UserCache userCache.Repository
}

func Register(container *dig.Container) error {
	if err := container.Provide(userCacheImpl.New); err != nil {
		return err
	}
	return nil
}

