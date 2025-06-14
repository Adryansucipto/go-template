package auth

import (
	"go-template/internal/controller/model"
	"go-template/util"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (c *Controller) LoginFunction(eCtx echo.Context) error {
	var response util.Response
	ctx := eCtx.Request().Context()

	var request model.AuthRequest
	if err := eCtx.Bind(&request); err != nil {
		response.ResponseCode = http.StatusUnauthorized
		response.ResponseError = err
		return util.HttpResponses(eCtx, response)
	}

	res, err := c.Auth.LoginHandler(ctx, request)
	if err != nil {
		return util.HttpResponses(eCtx, res)
	}

	return util.HttpResponses(eCtx, res)
}
