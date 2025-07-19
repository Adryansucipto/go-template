package auth

import (
	"go-template/internal/controller/model"
	"go-template/util"
	CONST "go-template/util/consts"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (c *Controller) LoginFunction(eCtx echo.Context) error {
	var response util.Response
	ctx := eCtx.Request().Context()

	var request model.AuthRequest
	if err := eCtx.Bind(&request); err != nil {
		response.ResponseCode = http.StatusUnauthorized
		response.ResponseMessage = CONST.ErrorMessages[http.StatusUnauthorized]
		response.Errors = err
		util.ErrorHandler(tag, err)
		return util.HttpResponses(eCtx, response)
	}

	res, err := c.Auth.LoginHandler(ctx, request)
	if res.ResponseCode != 200 {
		util.ErrorHandler(tag, err)
		return util.HttpResponses(eCtx, res)
	}

	return util.HttpResponses(eCtx, res)
}

func (c *Controller) LogoutFunction(eCtx echo.Context) error {
	var response util.Response
	ctx := eCtx.Request().Context()

	token, response, err := c.Auth.ValidateToken(ctx, eCtx.Request().Header.Get(echo.HeaderAuthorization))
	if response.ResponseCode != 200 {
		util.ErrorHandler(tag, err)
		return util.HttpResponses(eCtx, response)
	}

	record, response, err := c.Auth.AuthorizeUser(ctx, token)
	if response.ResponseCode != 200 {
		util.ErrorHandler(tag, err)
		return util.HttpResponses(eCtx, response)
	}

	response = c.Auth.DeleteSession(ctx, record.Username)
	if response.ResponseCode != 200 {
		util.ErrorHandler(tag, response.Errors)
		return util.HttpResponses(eCtx, response)
	}

	return util.HttpResponses(eCtx, response)
}

func (c *Controller) RefreshFunction(eCtx echo.Context) error {
	var response util.Response
	ctx := eCtx.Request().Context()

	username, response, err := c.Auth.ValidateRefreshToken(ctx, eCtx.Request().Header.Get(echo.HeaderAuthorization))
	if response.ResponseCode != 200 {
		util.ErrorHandler(tag, err)
		return util.HttpResponses(eCtx, response)
	}

	response, err = c.Auth.RefreshNewToken(ctx, username)
	if response.ResponseCode != 200 {
		util.ErrorHandler(tag, err)
		return util.HttpResponses(eCtx, response)
	}

	return util.HttpResponses(eCtx, response)
}

func (c *Controller) AuthorizeFunction(eCtx echo.Context) error {
	var response util.Response
	ctx := eCtx.Request().Context()

	token, response, err := c.Auth.ValidateToken(ctx, eCtx.Request().Header.Get(echo.HeaderAuthorization))
	if response.ResponseCode != 200 {
		util.ErrorHandler(tag, err)
		return util.HttpResponses(eCtx, response)
	}

	_, response, err = c.Auth.AuthorizeUser(ctx, token)
	if response.ResponseCode != 200 {
		util.ErrorHandler(tag, err)
		return util.HttpResponses(eCtx, response)
	}

	return util.HttpResponses(eCtx, response)
}
