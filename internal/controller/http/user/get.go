package user

import (
	"go-template/pkg/response"

	"github.com/labstack/echo/v4"
)

// GET /v1/users
func (c *Controller) GetUsers(eCtx echo.Context) error {
	ctx := eCtx.Request().Context()

	res, err := c.User.GetUsers(ctx)
	if err != nil {
		response.ErrorHandler(c.Logger, tag, res.ResponseMessage)
		return response.HttpResponses(eCtx, res)
	}
	return response.HttpResponses(eCtx, res)
}
