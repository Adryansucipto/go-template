package db

import (
	userRepository "go-template/internal/repository/db/user"
	userRepositoryImpl "go-template/internal/repository/db/user/impl"

	authRepository "go-template/internal/repository/db/auth"
	authRepositoryImpl "go-template/internal/repository/db/auth/impl"

	outboxRepository "go-template/internal/repository/db/outbox"
	outboxRepositoryImpl "go-template/internal/repository/db/outbox/impl"

	"go-template/internal/repository/db/tx"

	"go.uber.org/dig"
)

type Repository struct {
	dig.In
	TxManager        tx.Manager
	UserRepository   userRepository.Repository
	AuthRepository   authRepository.Repository
	OutboxRepository outboxRepository.Repository
}

func Register(container *dig.Container) (err error) {
	if err := container.Provide(tx.New); err != nil {
		return err
	}
	if err := container.Provide(userRepositoryImpl.New); err != nil {
		return err
	}
	if err := container.Provide(authRepositoryImpl.New); err != nil {
		return err
	}
	if err := container.Provide(outboxRepositoryImpl.New); err != nil {
		return err
	}
	return nil
}
