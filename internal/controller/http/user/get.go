package user

import (
	"go-template/pkg/logger"
	"go-template/pkg/response"

	"github.com/labstack/echo/v4"
)

// GET /v1/users
func (c *Controller) GetUsers(eCtx echo.Context) error {
	ctx := eCtx.Request().Context()
	log := logger.FromContext(ctx, c.Logger)

	res, err := c.User.GetUsers(ctx)
	if err != nil {
		response.ErrorHandler(log, tag, res.ResponseMessage)
		return response.HttpResponses(eCtx, res)
	}
	return response.HttpResponses(eCtx, res)
}
