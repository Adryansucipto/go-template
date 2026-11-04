package repository

import (
	cacheRepo "go-template/internal/repository/cache"
	clientRepo "go-template/internal/repository/client"
	dbRepo "go-template/internal/repository/db"
	msgRepo "go-template/internal/repository/messaging"

	"go.uber.org/dig"
)

type Repository struct {
	dig.In

	DBRepository        dbRepo.Repository
	CacheRepository     cacheRepo.Repository
	MessagingRepository msgRepo.Repository
	ClientRepository    clientRepo.Repository
}

func Register(container *dig.Container) (err error) {
	if err := dbRepo.Register(container); err != nil {
		return err
	}

	if err := cacheRepo.Register(container); err != nil {
		return err
	}

	if err := msgRepo.Register(container); err != nil {
		return err
	}

	if err := clientRepo.Register(container); err != nil {
		return err
	}

	return nil
}
