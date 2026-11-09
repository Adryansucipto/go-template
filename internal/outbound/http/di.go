package http

import (
	partnerHTTP "go-template/internal/outbound/http/partner"
	partnerHTTPImpl "go-template/internal/outbound/http/partner/impl"

	"go.uber.org/dig"
)

type Outbound struct {
	dig.In

	Partner partnerHTTP.Outbound
}

func Register(container *dig.Container) error {
	if err := container.Provide(partnerHTTPImpl.New); err != nil {
		return err
	}
	return nil
}
