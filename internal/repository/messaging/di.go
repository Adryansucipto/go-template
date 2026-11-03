package messaging

import (
	outboxMsg "go-template/internal/repository/messaging/outbox"
	outboxMsgImpl "go-template/internal/repository/messaging/outbox/impl"

	"go.uber.org/dig"
)

type Repository struct {
	dig.In
	Producer outboxMsg.Repository
}

func Register(container *dig.Container) error {
	if err := container.Provide(outboxMsgImpl.New); err != nil {
		return err
	}
	return nil
}

