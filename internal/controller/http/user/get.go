package user

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// GET /users
func (c *Controller) GetUsers(eCtx echo.Context) error {
	ctx := eCtx.Request().Context()
	users, err := c.User.GetUsers(ctx)
	if err != nil {
		return eCtx.JSON(http.StatusInternalServerError, nil)
	}
	return eCtx.JSON(http.StatusOK, users)
}
