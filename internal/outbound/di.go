package outbound

import (
	"go-template/internal/outbound/cache"
	"go-template/internal/outbound/grpc"
	"go-template/internal/outbound/http"

	"go.uber.org/dig"
)

// Outbound represents the peer layer for external communications (Redis Cache, Cross-Service HTTP & gRPC).
type Outbound struct {
	dig.In

	Cache cache.Outbound
	HTTP  http.Outbound
	GRPC  grpc.Outbound
}

func Register(container *dig.Container) error {
	if err := cache.Register(container); err != nil {
		return err
	}
	if err := http.Register(container); err != nil {
		return err
	}
	if err := grpc.Register(container); err != nil {
		return err
	}
	return nil
}
