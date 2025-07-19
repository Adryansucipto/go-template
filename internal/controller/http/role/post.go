package role

import (
	"go-template/internal/controller/model"
	"go-template/util"
	CONST "go-template/util/consts"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (c *Controller) CreateRole(eCtx echo.Context) error {
	var response util.Response
	var request model.CreateRoleRequest
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

	if err := eCtx.Bind(&request); err != nil {
		response.ResponseCode = http.StatusBadRequest
		response.ResponseMessage = CONST.ErrorMessages[http.StatusBadRequest]
		response.Errors = "Error Bind Request"
		util.ErrorHandler(tag, response.ResponseMessage)
		return util.HttpResponses(eCtx, response)
	}

	response, err = c.Role.CreateRole(ctx, request, record.Username)
	if response.ResponseCode != 200 {
		util.ErrorHandler(tag, err)
		return util.HttpResponses(eCtx, response)
	}

	return util.HttpResponses(eCtx, response)
}
