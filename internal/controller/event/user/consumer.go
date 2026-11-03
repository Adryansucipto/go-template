package user

import (
	"context"
	"encoding/json"
	"errors"
	"go-template/internal/repository"
	userUsecase "go-template/internal/usecase/user"
	"go-template/pkg/event"
	"go-template/pkg/resource/model"
	"time"

	"github.com/segmentio/kafka-go"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

const TopicUserEvents = "user.events"

type Consumer struct {
	dig.In

	User       userUsecase.Usecase
	Repository repository.Repository
	Kafka      model.Kafka
	Logger     *zap.Logger
}

func (c *Consumer) Start(ctx context.Context) {
	if !c.Kafka.Enabled {
		c.Logger.Info("[UserEventConsumer] Kafka is disabled, skipping consumer", zap.String("topic", TopicUserEvents))
		return
	}

	reader := c.Kafka.NewReader(TopicUserEvents)
	defer func() {
		_ = reader.Close()
	}()

	c.Logger.Info("[UserEventConsumer] Listening for events",
		zap.String("topic", TopicUserEvents),
		zap.String("group_id", c.Kafka.GroupID),
	)

	for {
		msg, err := reader.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				c.Logger.Info("[UserEventConsumer] Consumer stopped cleanly")
				return
			}
			c.Logger.Error("[UserEventConsumer] failed fetching message", zap.Error(err))
			time.Sleep(time.Second)
			continue
		}

		c.processWithRetryAndDLQ(ctx, reader, msg)
	}
}

func (c *Consumer) processWithRetryAndDLQ(ctx context.Context, reader *kafka.Reader, msg kafka.Message) {
	var lastErr error
	maxRetries := c.Kafka.MaxRetries
	if maxRetries <= 0 {
		maxRetries = 3
	}
	backoff := time.Duration(c.Kafka.RetryBackoffMs) * time.Millisecond

	for attempt := 1; attempt <= maxRetries; attempt++ {
		if err := c.handleMessage(ctx, msg); err != nil {
			lastErr = err
			c.Logger.Warn("[UserEventConsumer] handler failed, retrying...",
				zap.String("topic", msg.Topic),
				zap.Int("attempt", attempt),
				zap.Int("max_retries", maxRetries),
				zap.Error(err),
			)
			select {
			case <-ctx.Done():
				return
			case <-time.After(backoff * time.Duration(attempt)):
			}
			continue
		}

		// Successfully processed -> commit offset
		if err := reader.CommitMessages(ctx, msg); err != nil {
			c.Logger.Error("[UserEventConsumer] failed to commit offset", zap.Error(err))
		}
		return
	}

	// Exhausted all retries -> forward to Dead Letter Queue (DLQ)
	c.Logger.Error("[UserEventConsumer] max retries exhausted, publishing to DLQ",
		zap.String("topic", msg.Topic),
		zap.String("key", string(msg.Key)),
		zap.Error(lastErr),
	)

	if err := c.Repository.MessagingRepository.Producer.PublishDLQ(
		ctx,
		msg.Topic,
		string(msg.Key),
		msg.Value,
		lastErr.Error(),
	); err != nil {
		c.Logger.Error("[UserEventConsumer] failed publishing message to DLQ", zap.Error(err))
		return
	}

	// Commit original message after successfully moving to DLQ
	if err := reader.CommitMessages(ctx, msg); err != nil {
		c.Logger.Error("[UserEventConsumer] failed to commit offset after DLQ", zap.Error(err))
	}
}

func (c *Consumer) handleMessage(ctx context.Context, msg kafka.Message) error {
	var envelope event.CloudEvent
	if err := json.Unmarshal(msg.Value, &envelope); err != nil {
		return err
	}

	c.Logger.Info("[UserEventConsumer] Received CloudEvent",
		zap.String("event_id", envelope.ID),
		zap.String("event_type", envelope.Type),
		zap.String("source", envelope.Source),
	)

	return nil
}

