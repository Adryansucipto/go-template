package repository

import (
	authRepository "go-template/internal/repository/auth"
	authRepositoryImpl "go-template/internal/repository/auth/impl"
	eventQueueRepository "go-template/internal/repository/event_queue"
	eventQueueRepositoryImpl "go-template/internal/repository/event_queue/impl"
	"go-template/internal/repository/transaction"
	userRepository "go-template/internal/repository/user"
	userRepositoryImpl "go-template/internal/repository/user/impl"

	"go.uber.org/dig"
)

// Repository strictly contains Database access repositories and DB Transaction management.
type Repository struct {
	dig.In

	TransactionManager   transaction.Manager
	UserRepository       userRepository.Repository
	AuthRepository       authRepository.Repository
	EventQueueRepository eventQueueRepository.Repository
}

func Register(container *dig.Container) error {
	if err := container.Provide(transaction.New); err != nil {
		return err
	}
	if err := container.Provide(userRepositoryImpl.New); err != nil {
		return err
	}
	if err := container.Provide(authRepositoryImpl.New); err != nil {
		return err
	}
	if err := container.Provide(eventQueueRepositoryImpl.New); err != nil {
		return err
	}
	return nil
}
