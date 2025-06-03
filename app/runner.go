package app

import (
	"go-template/internal/controller"

	"github.com/labstack/echo/v4"
)

func Run() error {
	container := BuildContainer()

	// Start echo server from injected controller
	return container.Invoke(func(ctrl controller.Controller) error {
		e := echo.New()

		ctrl.HTTP.Routes(e)

		// Start the HTTP server
		return e.Start(":8080")
	})
}
