package user

import (
	"go-template/util"

	"github.com/labstack/echo/v4"
)

// GET /users
func (c *Controller) GetUsers(eCtx echo.Context) error {
	ctx := eCtx.Request().Context()

	res, err := c.User.GetUsers(ctx)
	if err != nil {
		return util.HttpResponses(eCtx, res)
	}
	return util.HttpResponses(eCtx, res)
}
