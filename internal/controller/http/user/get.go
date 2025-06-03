package user

import (
	"go-template/internal/controller/model"
	"net/http"

	"github.com/labstack/echo/v4"
)

var users = []model.User{
	{ID: 1, Name: "Alice"},
	{ID: 2, Name: "Bob"},
}

// GET /users
func (c *Controller) GetUsers(eCtx echo.Context) error {
	return eCtx.JSON(http.StatusOK, users)
}
