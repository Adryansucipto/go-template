package repository

import (
	dbRepo "go-template/internal/repository/db"

	"go.uber.org/dig"
)

type Repository struct {
	dig.In

	Repository dbRepo.Repository
}

func Register(container *dig.Container) (err error) {
	if err := dbRepo.Register(container); err != nil {
		return err
	}

	return nil
}
