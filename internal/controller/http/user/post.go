package user

import (
	"go-template/internal/controller/model"
	"go-template/pkg/response"
	"net/http"

	"github.com/labstack/echo/v4"
)

// POST /v1/register
func (c *Controller) CreateUser(eCtx echo.Context) error {
	ctx := eCtx.Request().Context()

	var req model.CreateUserRequest
	if err := eCtx.Bind(&req); err != nil {
		res := response.ResponseGenerate(http.StatusBadRequest, err, err.Error(), nil, nil)
		response.ErrorHandler(c.Logger, tag, res.ResponseMessage)
		return response.HttpResponses(eCtx, res)
	}

	if err := eCtx.Validate(&req); err != nil {
		res := response.ResponseGenerate(http.StatusBadRequest, err, err.Error(), nil, nil)
		response.ErrorHandler(c.Logger, tag, res.ResponseMessage)
		return response.HttpResponses(eCtx, res)
	}

	res, err := c.User.CreateUser(ctx, req)
	if err != nil {
		response.ErrorHandler(c.Logger, tag, res.ResponseMessage)
		return response.HttpResponses(eCtx, res)
	}

	return response.HttpResponses(eCtx, res)
}
