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
		response.ResponseMessage = err
		util.ErrorHandler(tag, response.ResponseMessage)
		return util.HttpResponses(eCtx, response)
	}

	res := c.Auth.LoginHandler(ctx, request)
	if res.ResponseCode != 200 {
		util.ErrorHandler(tag, res.ResponseMessage)
		return util.HttpResponses(eCtx, res)
	}

	return util.HttpResponses(eCtx, res)
}

func (c *Controller) LogoutFunction(eCtx echo.Context) error {
	var response util.Response
	ctx := eCtx.Request().Context()

	record, response := c.Auth.GetRecordByToken(ctx, eCtx.Request().Header.Get(echo.HeaderAuthorization))
	if response.ResponseCode != 200 {
		util.ErrorHandler(tag, response.ResponseMessage)
		return util.HttpResponses(eCtx, response)
	}
	response = c.Auth.DeleteSession(ctx, record)
	if response.ResponseCode != 200 {
		util.ErrorHandler(tag, response.ResponseMessage)
		return util.HttpResponses(eCtx, response)
	}

	return util.HttpResponses(eCtx, response)
}
