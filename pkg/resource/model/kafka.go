package model

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type Kafka struct {
	Enabled        bool
	Brokers        []string
	GroupID        string
	MaxRetries     int
	RetryBackoffMs int
	DLQTopicSuffix string
	Writer         *kafka.Writer
}

func (k *Kafka) NewReader(topic string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  k.Brokers,
		GroupID:  k.GroupID,
		Topic:    topic,
		MinBytes: 1,
		MaxBytes: 10e6, // 10MB
	})
}

func (k *Kafka) Publish(ctx context.Context, topic string, key, value []byte, headers ...kafka.Header) error {
	if !k.Enabled || k.Writer == nil {
		return nil
	}
	return k.Writer.WriteMessages(ctx, kafka.Message{
		Topic:   topic,
		Key:     key,
		Value:   value,
		Headers: headers,
	})
}

func (k *Kafka) Close() error {
	if k.Writer != nil {
		return k.Writer.Close()
	}
	return nil
}
