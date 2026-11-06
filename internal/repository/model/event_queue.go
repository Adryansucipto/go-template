package model

import "time"

const (
	EventQueueStatusPending   = "PENDING"
	EventQueueStatusProcessed = "PROCESSED"
	EventQueueStatusFailed    = "FAILED"
)

type EventQueue struct {
	ID            int64      `gorm:"column:id;primaryKey;autoIncrement"`
	AggregateType string     `gorm:"column:aggregate_type"`
	AggregateID   string     `gorm:"column:aggregate_id"`
	Topic         string     `gorm:"column:topic"`
	Payload       string     `gorm:"column:payload"`
	Status        string     `gorm:"column:status"`
	RetryCount    int        `gorm:"column:retry_count"`
	LastError     string     `gorm:"column:last_error"`
	CreatedAt     time.Time  `gorm:"column:created_at"`
	ProcessedAt   *time.Time `gorm:"column:processed_at"`
}

func (t *EventQueue) TableName() string {
	return `"public"."EVENT_QUEUE"`
}
