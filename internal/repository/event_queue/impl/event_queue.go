package impl

import (
	"context"
	"go-template/internal/repository/model"
	"go-template/internal/repository/transaction"
	"go-template/pkg/response"
	"time"
)

func (r *Repository) CreateEventQueue(ctx context.Context, event model.EventQueue) error {
	db := transaction.ExtractTx(ctx, r.DbGorm)
	if event.Status == "" {
		event.Status = model.EventQueueStatusPending
	}
	if event.CreatedAt.IsZero() {
		event.CreatedAt = time.Now().UTC()
	}
	if err := db.Create(&event).Error; err != nil {
		return response.ErrorGenerate(tag, err)
	}
	return nil
}

func (r *Repository) FetchPendingEvents(ctx context.Context, limit int) ([]model.EventQueue, error) {
	var events []model.EventQueue
	db := transaction.ExtractTx(ctx, r.DbGorm)

	if err := db.Where("status = ?", model.EventQueueStatusPending).
		Order("created_at ASC").
		Limit(limit).
		Find(&events).Error; err != nil {
		return nil, response.ErrorGenerate(tag, err)
	}
	return events, nil
}

func (r *Repository) MarkProcessed(ctx context.Context, id int64) error {
	now := time.Now().UTC()
	db := transaction.ExtractTx(ctx, r.DbGorm)
	if err := db.Model(&model.EventQueue{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"status":       model.EventQueueStatusProcessed,
			"processed_at": &now,
		}).Error; err != nil {
		return response.ErrorGenerate(tag, err)
	}
	return nil
}

func (r *Repository) MarkFailed(ctx context.Context, id int64, retryCount int, maxRetries int, lastErr string) error {
	status := model.EventQueueStatusPending
	if retryCount >= maxRetries {
		status = model.EventQueueStatusFailed
	}

	db := transaction.ExtractTx(ctx, r.DbGorm)
	if err := db.Model(&model.EventQueue{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"status":      status,
			"retry_count": retryCount,
			"last_error":  lastErr,
		}).Error; err != nil {
		return response.ErrorGenerate(tag, err)
	}
	return nil
}

