package eventqueue

import (
	"context"
	"go-template/internal/repository"
	"go-template/internal/service/event/producer"
	"go-template/pkg/config"

	"go.uber.org/dig"
	"go.uber.org/zap"
)

// RelayJob reads PENDING rows from EVENT_QUEUE and publishes them to Kafka.
type RelayJob struct {
	dig.In

	Repository repository.Repository
	Producer   producer.Service
	Config     config.Config
	Logger     *zap.Logger
}

func (j *RelayJob) ProcessBatch(ctx context.Context, batchSize int) {
	if !j.Config.Kafka.Enabled {
		return
	}

	if batchSize <= 0 {
		batchSize = j.Config.Kafka.EventQueueBatchSize
	}
	if batchSize <= 0 {
		batchSize = 50
	}

	maxRetries := j.Config.Kafka.MaxRetries
	if maxRetries <= 0 {
		maxRetries = 3
	}

	events, err := j.Repository.EventQueueRepository.FetchPendingEvents(ctx, batchSize)
	if err != nil {
		j.Logger.Error("[EventQueueRelayJob] failed to fetch pending events", zap.Error(err))
		return
	}

	for _, evt := range events {
		err := j.Producer.PublishRaw(ctx, evt.Topic, evt.AggregateID, []byte(evt.Payload))
		if err != nil {
			nextRetry := evt.RetryCount + 1
			j.Logger.Warn("[EventQueueRelayJob] failed to publish event",
				zap.Int64("event_queue_id", evt.ID),
				zap.String("topic", evt.Topic),
				zap.Int("retry_count", nextRetry),
				zap.Error(err),
			)
			_ = j.Repository.EventQueueRepository.MarkFailed(ctx, evt.ID, nextRetry, maxRetries, err.Error())
			continue
		}

		if err := j.Repository.EventQueueRepository.MarkProcessed(ctx, evt.ID); err != nil {
			j.Logger.Error("[EventQueueRelayJob] failed to mark event processed",
				zap.Int64("event_queue_id", evt.ID),
				zap.Error(err),
			)
		}
	}
}

