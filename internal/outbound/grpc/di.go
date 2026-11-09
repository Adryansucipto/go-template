package grpc

import (
	partnerGRPC "go-template/internal/outbound/grpc/partner"
	partnerGRPCImpl "go-template/internal/outbound/grpc/partner/impl"

	"go.uber.org/dig"
)

type Outbound struct {
	dig.In

	Partner partnerGRPC.Outbound
}

func Register(container *dig.Container) error {
	if err := container.Provide(partnerGRPCImpl.New); err != nil {
		return err
	}
	return nil
}
