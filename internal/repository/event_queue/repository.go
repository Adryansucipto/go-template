package eventqueue

import (
	"context"
	"go-template/internal/repository/model"
)

// Repository stores events in the EVENT_QUEUE table (Transactional Outbox pattern)
// so they are saved atomically with business data, then relayed to Kafka by the scheduler.
type Repository interface {
	CreateEventQueue(ctx context.Context, event model.EventQueue) error
	FetchPendingEvents(ctx context.Context, limit int) ([]model.EventQueue, error)
	MarkProcessed(ctx context.Context, id int64) error
	MarkFailed(ctx context.Context, id int64, retryCount int, maxRetries int, lastErr string) error
}

