package main

import (
	"fmt"
	"go-template/internal/controller"
	pkgDi "go-template/pkg/di"
	"os"

	"github.com/labstack/echo/v4"
	"go.uber.org/dig"
	"go.uber.org/zap"
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

		logger, _ := zap.NewProduction()
		defer logger.Sync()

		ctrl.HTTP.Routes(e, logger)

		// Start the HTTP server
		port := os.Getenv("PORT")
		if port == "" {
			port = "9001" // fallback saat lokal
		}

		return e.Start(fmt.Sprintf(":%s", port))
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
