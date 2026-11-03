package client

import (
	partnerGRPC "go-template/internal/repository/client/grpc/partner"
	partnerGRPCImpl "go-template/internal/repository/client/grpc/partner/impl"
	partnerHTTP "go-template/internal/repository/client/http/partner"
	partnerHTTPImpl "go-template/internal/repository/client/http/partner/impl"

	"go.uber.org/dig"
)

type Repository struct {
	dig.In

	PartnerHTTP partnerHTTP.Repository
	PartnerGRPC partnerGRPC.Repository
}

func Register(container *dig.Container) error {
	if err := container.Provide(partnerHTTPImpl.New); err != nil {
		return err
	}
	if err := container.Provide(partnerGRPCImpl.New); err != nil {
		return err
	}
	return nil
}

