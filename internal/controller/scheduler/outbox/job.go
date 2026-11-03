package outbox

import (
	"context"
	"go-template/internal/repository"
	"go-template/pkg/config"

	"go.uber.org/dig"
	"go.uber.org/zap"
)

type RelayJob struct {
	dig.In

	Repository repository.Repository
	Config     config.Config
	Logger     *zap.Logger
}

func (j *RelayJob) ProcessBatch(ctx context.Context, batchSize int) {
	if !j.Config.Kafka.Enabled {
		return
	}

	if batchSize <= 0 {
		batchSize = j.Config.Kafka.OutboxBatchSize
	}
	if batchSize <= 0 {
		batchSize = 50
	}

	maxRetries := j.Config.Kafka.MaxRetries
	if maxRetries <= 0 {
		maxRetries = 3
	}

	events, err := j.Repository.DBRepository.OutboxRepository.FetchPendingEvents(ctx, batchSize)
	if err != nil {
		j.Logger.Error("[OutboxRelayJob] failed to fetch pending outbox events", zap.Error(err))
		return
	}

	for _, evt := range events {
		err := j.Repository.MessagingRepository.Producer.PublishRaw(
			ctx,
			evt.Topic,
			evt.AggregateID,
			[]byte(evt.Payload),
		)
		if err != nil {
			nextRetry := evt.RetryCount + 1
			j.Logger.Warn("[OutboxRelayJob] failed to publish outbox event",
				zap.Int64("outbox_id", evt.ID),
				zap.String("topic", evt.Topic),
				zap.Int("retry_count", nextRetry),
				zap.Error(err),
			)
			_ = j.Repository.DBRepository.OutboxRepository.MarkFailed(ctx, evt.ID, nextRetry, maxRetries, err.Error())
			continue
		}

		if err := j.Repository.DBRepository.OutboxRepository.MarkProcessed(ctx, evt.ID); err != nil {
			j.Logger.Error("[OutboxRelayJob] failed to mark outbox event processed",
				zap.Int64("outbox_id", evt.ID),
				zap.Error(err),
			)
		}
	}
}

