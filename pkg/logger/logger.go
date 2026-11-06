package logger

import (
	"context"
	"fmt"
	"go-template/pkg/config"

	"go.uber.org/dig"
	"go.uber.org/zap"
)

type contextKey string

const RequestIDKey contextKey = "request_id"

// WithRequestID injects a request/trace ID into context.Context.
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestIDKey, requestID)
}

// RequestIDFromContext extracts request/trace ID from context.Context.
func RequestIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if val, ok := ctx.Value(RequestIDKey).(string); ok {
		return val
	}
	return ""
}

// FromContext returns a logger tagged with request_id if present in context.
func FromContext(ctx context.Context, l *zap.Logger) *zap.Logger {
	if l == nil {
		return zap.NewNop()
	}
	reqID := RequestIDFromContext(ctx)
	if reqID != "" {
		return l.With(zap.String("request_id", reqID))
	}
	return l
}

func New(cfg config.Config) (*zap.Logger, error) {
	if cfg.AppConfig.Env == "prod" {
		return zap.NewProduction()
	}
	return zap.NewDevelopment()
}

func Register(container *dig.Container) error {
	if err := container.Provide(New); err != nil {
		return fmt.Errorf("[DI] cannot initialize logger: %+v", err)
	}
	return nil
}
