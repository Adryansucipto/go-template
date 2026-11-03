package outbox

import (
	"context"
	"go-template/internal/repository/model"
)

type Repository interface {
	CreateOutboxEvent(ctx context.Context, event model.OutboxEvent) error
	FetchPendingEvents(ctx context.Context, limit int) ([]model.OutboxEvent, error)
	MarkProcessed(ctx context.Context, id int64) error
	MarkFailed(ctx context.Context, id int64, retryCount int, maxRetries int, lastErr string) error
}

