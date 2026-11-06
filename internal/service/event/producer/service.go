package producer

import "context"

type Service interface {
	PublishEvent(ctx context.Context, topic, key string, payload any) error
	PublishRaw(ctx context.Context, topic, key string, rawPayload []byte) error
	PublishDLQ(ctx context.Context, originalTopic, key string, rawPayload []byte, reason string) error
}
