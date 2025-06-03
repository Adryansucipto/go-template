package app

import (
	"go-template/internal/controller"

	"go.uber.org/dig"
)

func BuildContainer() *dig.Container {
	container := dig.New()

	// Provide all dependencies
	container.Provide(func() controller.Controller {
		return controller.Controller{}
	})

	return container
}
