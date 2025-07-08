package http

import (
	"go-template/internal/controller/http/auth"
	"go-template/internal/controller/http/user"
	"go-template/middleware"

	"github.com/labstack/echo/v4"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

type Controller struct {
	dig.In

	User user.Controller
	Auth auth.Controller
}

func (c *Controller) Routes(ec *echo.Echo, logger *zap.Logger) {

	ec.Use(middleware.ZapLogger(logger))

	ec.GET("/", func(c echo.Context) error {
		return c.String(200, "Go API is up and running!")
	})
	ec.GET("/ping", func(c echo.Context) error {
		return c.String(200, "Test Ping")
	})

	base := ec.Group("")
	v1 := base.Group("/v1")

	// user
	userGroup := v1.Group("/users")
	userGroup.GET("", c.User.GetUsers)
	//userGroup.POST("", c.User.CreateUser)

	// auth
	authGroup := v1.Group("/auth")
	authGroup.POST("/login", c.Auth.LoginFunction)
	authGroup.POST("/logout", c.Auth.LogoutFunction)
	authGroup.POST("/refresh-token", c.Auth.RefreshFunction)
}
