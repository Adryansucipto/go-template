package http

import (
	"go-template/internal/controller/http/auth"
	"go-template/internal/controller/http/user"

	"github.com/labstack/echo/v4"
	"go.uber.org/dig"
)

type Controller struct {
	dig.In

	User user.Controller
	Auth auth.Controller
}

func (c *Controller) Routes(ec *echo.Echo) {

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
	authGroup := v1.Group("/login")
	authGroup.POST("", c.Auth.LoginFunction)

}
