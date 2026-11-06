package consumer

import (
	"context"
	"go-template/internal/service/event/consumer/user"

	"go.uber.org/dig"
	"go.uber.org/zap"
)

type Service struct {
	dig.In

	UserConsumer user.Consumer
	Logger       *zap.Logger
}

func (s *Service) Start(ctx context.Context) {
	go s.UserConsumer.Start(ctx)
}
