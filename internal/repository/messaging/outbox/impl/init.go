package impl

import (
	"context"
	"fmt"
	msgOutbox "go-template/internal/repository/messaging/outbox"
	"go-template/pkg/event"
	"go-template/pkg/resource/model"
	"time"

	"github.com/segmentio/kafka-go"
)

type Repository struct {
	Kafka model.Kafka
}

func New(k model.Kafka) msgOutbox.Repository {
	return &Repository{
		Kafka: k,
	}
}

func (r *Repository) PublishEvent(ctx context.Context, topic, key string, payload any) error {
	eventID := fmt.Sprintf("%s-%d", key, time.Now().UnixNano())
	_, rawBytes, err := event.NewCloudEvent(eventID, topic, "go-template", payload)
	if err != nil {
		return err
	}
	return r.Kafka.Publish(ctx, topic, []byte(key), rawBytes)
}

func (r *Repository) PublishRaw(ctx context.Context, topic, key string, rawPayload []byte) error {
	return r.Kafka.Publish(ctx, topic, []byte(key), rawPayload)
}

func (r *Repository) PublishDLQ(ctx context.Context, originalTopic, key string, rawPayload []byte, reason string) error {
	dlqTopic := originalTopic + r.Kafka.DLQTopicSuffix
	return r.Kafka.Publish(
		ctx,
		dlqTopic,
		[]byte(key),
		rawPayload,
		kafka.Header{Key: "x-original-topic", Value: []byte(originalTopic)},
		kafka.Header{Key: "x-error-reason", Value: []byte(reason)},
		kafka.Header{Key: "x-failed-at", Value: []byte(time.Now().UTC().Format(time.RFC3339))},
	)
}

