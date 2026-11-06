package impl

import (
	"context"
	"fmt"
	producerSvc "go-template/internal/service/event/producer"
	"go-template/pkg/event"
	"go-template/pkg/resource/model"
	"time"

	"github.com/segmentio/kafka-go"
)

type Service struct {
	Kafka model.Kafka
}

func New(k model.Kafka) producerSvc.Service {
	return &Service{
		Kafka: k,
	}
}

func (s *Service) PublishEvent(ctx context.Context, topic, key string, payload any) error {
	eventID := fmt.Sprintf("%s-%d", key, time.Now().UnixNano())
	_, rawBytes, err := event.NewCloudEvent(eventID, topic, "go-template", payload)
	if err != nil {
		return err
	}
	return s.Kafka.Publish(ctx, topic, []byte(key), rawBytes)
}

func (s *Service) PublishRaw(ctx context.Context, topic, key string, rawPayload []byte) error {
	return s.Kafka.Publish(ctx, topic, []byte(key), rawPayload)
}

func (s *Service) PublishDLQ(ctx context.Context, originalTopic, key string, rawPayload []byte, reason string) error {
	dlqTopic := originalTopic + s.Kafka.DLQTopicSuffix
	return s.Kafka.Publish(
		ctx,
		dlqTopic,
		[]byte(key),
		rawPayload,
		kafka.Header{Key: "x-original-topic", Value: []byte(originalTopic)},
		kafka.Header{Key: "x-error-reason", Value: []byte(reason)},
		kafka.Header{Key: "x-failed-at", Value: []byte(time.Now().UTC().Format(time.RFC3339))},
	)
}
