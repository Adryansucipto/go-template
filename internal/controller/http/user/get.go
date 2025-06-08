package user

import (
	"go-template/util"
	"net/http"

	"github.com/labstack/echo/v4"
)

// GET /users
func (c *Controller) GetUsers(eCtx echo.Context) error {
	var response util.Response
	ctx := eCtx.Request().Context()
	users, err := c.User.GetUsers(ctx)
	if err != nil {
		response.ResponseCode = http.StatusOK
		response.ResponseError = err
		return util.HttpResponses(eCtx, response)
	}
	response.ResponseCode = http.StatusOK
	response.ResponseError = nil
	response.Data = users

	return util.HttpResponses(eCtx, response)
}
