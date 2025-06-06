package main

import (
	"go-template/internal/controller"
	pkgDi "go-template/pkg/di"

	"github.com/labstack/echo/v4"
	"go.uber.org/dig"
)

type (
	ContainerCall func() (*dig.Container, error)
	Invoke        func(container *dig.Container) error
	InvokeError   func(container *dig.Container, err error)
)

func Run(containerCall ContainerCall, invoke Invoke, onError InvokeError) error {
	container, err := containerCall()
	if err != nil {
		return err
	}

	if err := invoke(container); err != nil {
		onError(container, err)
	}
	return nil
}

func onRun(container *dig.Container) error {

	// Start echo server from injected controller
	return container.Invoke(func(ctrl controller.Controller) error {
		e := echo.New()

		ctrl.HTTP.Routes(e)

		// Start the HTTP server
		return e.Start(":9001")
	})
}

func OnError(container *dig.Container, err error) {
	// Bisa log error di sini atau shutdown service
	if container != nil && err != nil {
		// Log / cleanup if needed
		panic(err)
	}
}

func main() {
	if err := Run(pkgDi.Container, onRun, OnError); err != nil {
		panic(err)
	}
}
