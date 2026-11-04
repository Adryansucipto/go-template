package impl

import (
	"context"
	"go-template/internal/repository/db/tx"
	"go-template/internal/repository/model"
	"go-template/pkg/response"
	"time"
)

func (r *Repository) CreateOutboxEvent(ctx context.Context, event model.OutboxEvent) error {
	db := tx.ExtractTx(ctx, r.DbGorm)
	if event.Status == "" {
		event.Status = model.OutboxStatusPending
	}
	if event.CreatedAt.IsZero() {
		event.CreatedAt = time.Now().UTC()
	}
	if err := db.Create(&event).Error; err != nil {
		return response.ErrorGenerate(tag, err)
	}
	return nil
}

func (r *Repository) FetchPendingEvents(ctx context.Context, limit int) ([]model.OutboxEvent, error) {
	var events []model.OutboxEvent
	db := tx.ExtractTx(ctx, r.DbGorm)

	if err := db.Where("status = ?", model.OutboxStatusPending).
		Order("created_at ASC").
		Limit(limit).
		Find(&events).Error; err != nil {
		return nil, response.ErrorGenerate(tag, err)
	}
	return events, nil
}

func (r *Repository) MarkProcessed(ctx context.Context, id int64) error {
	now := time.Now().UTC()
	db := tx.ExtractTx(ctx, r.DbGorm)
	if err := db.Model(&model.OutboxEvent{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"status":       model.OutboxStatusProcessed,
			"processed_at": &now,
		}).Error; err != nil {
		return response.ErrorGenerate(tag, err)
	}
	return nil
}

func (r *Repository) MarkFailed(ctx context.Context, id int64, retryCount int, maxRetries int, lastErr string) error {
	status := model.OutboxStatusPending
	if retryCount >= maxRetries {
		status = model.OutboxStatusFailed
	}

	db := tx.ExtractTx(ctx, r.DbGorm)
	if err := db.Model(&model.OutboxEvent{}).
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
