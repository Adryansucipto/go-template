package repository

import (
	dbRepo "go-template/internal/repository/db"
	msgRepo "go-template/internal/repository/messaging"

	"go.uber.org/dig"
)

type Repository struct {
	dig.In

	DBRepository        dbRepo.Repository
	MessagingRepository msgRepo.Repository
}

func Register(container *dig.Container) (err error) {
	if err := dbRepo.Register(container); err != nil {
		return err
	}

	if err := msgRepo.Register(container); err != nil {
		return err
	}

	return nil
}
