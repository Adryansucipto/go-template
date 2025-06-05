package http

import (
	"go-template/internal/controller/http/user"

	"github.com/labstack/echo/v4"
	"go.uber.org/dig"
)

type Controller struct {
	dig.In

	User user.Controller
}

func (c *Controller) Routes(ec *echo.Echo) {

	base := ec.Group("")
	v1 := base.Group("/v1")

	// user
	userGroup := v1.Group("/users")
	userGroup.GET("", c.User.GetUsers)
	//userGroup.POST("", c.User.CreateUser)

}
