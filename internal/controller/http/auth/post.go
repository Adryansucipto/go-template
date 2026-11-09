package auth

import (
	"go-template/internal/controller/model"
	"go-template/pkg/logger"
	"go-template/pkg/response"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (c *Controller) LoginFunction(eCtx echo.Context) error {
	ctx := eCtx.Request().Context()
	log := logger.FromContext(ctx, c.Logger)

	var request model.AuthRequest
	if err := eCtx.Bind(&request); err != nil {
		res := response.ResponseGenerate(http.StatusBadRequest, err, err.Error(), nil, nil)
		response.ErrorHandler(log, tag, res.ResponseMessage)
		return response.HttpResponses(eCtx, res)
	}

	if err := eCtx.Validate(&request); err != nil {
		res := response.ResponseGenerate(http.StatusBadRequest, err, err.Error(), nil, nil)
		response.ErrorHandler(log, tag, res.ResponseMessage)
		return response.HttpResponses(eCtx, res)
	}

	res := c.Auth.LoginHandler(ctx, request)
	if res.ResponseCode != http.StatusOK {
		response.ErrorHandler(log, tag, res.ResponseMessage)
		return response.HttpResponses(eCtx, res)
	}

	return response.HttpResponses(eCtx, res)
}

func (c *Controller) LogoutFunction(eCtx echo.Context) error {
	ctx := eCtx.Request().Context()
	log := logger.FromContext(ctx, c.Logger)

	record, res := c.Auth.GetRecordByToken(ctx, eCtx.Request().Header.Get(echo.HeaderAuthorization))
	if res.ResponseCode != http.StatusOK {
		response.ErrorHandler(log, tag, res.ResponseMessage)
		return response.HttpResponses(eCtx, res)
	}

	res = c.Auth.DeleteSession(ctx, record)
	if res.ResponseCode != http.StatusOK {
		response.ErrorHandler(log, tag, res.ResponseMessage)
		return response.HttpResponses(eCtx, res)
	}

	return response.HttpResponses(eCtx, res)
}
