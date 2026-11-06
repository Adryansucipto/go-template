package cache

import (
	userCache "go-template/internal/outbound/cache/user"
	userCacheImpl "go-template/internal/outbound/cache/user/impl"

	"go.uber.org/dig"
)

type Outbound struct {
	dig.In

	User userCache.Outbound
}

func Register(container *dig.Container) error {
	if err := container.Provide(userCacheImpl.New); err != nil {
		return err
	}
	return nil
}
