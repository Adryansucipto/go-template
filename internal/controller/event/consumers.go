package event

import (
	"context"
	"go-template/internal/controller/event/user"

	"go.uber.org/dig"
	"go.uber.org/zap"
)

type Controller struct {
	dig.In

	UserConsumer user.Consumer
	Logger       *zap.Logger
}

func (c *Controller) StartConsumers(ctx context.Context) {
	go c.UserConsumer.Start(ctx)
}
