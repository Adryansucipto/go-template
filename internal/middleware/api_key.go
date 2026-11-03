package middleware

import (
	"errors"
	"go-template/pkg/config"
	"go-template/pkg/response"
	"net/http"

	"github.com/labstack/echo/v4"
)

func APIKeyAuth(apiKeyCfg config.APIKeyConfig) echo.MiddlewareFunc {
	headerName := apiKeyCfg.HeaderName()

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			incomingKey := c.Request().Header.Get(headerName)
			clientName, valid := apiKeyCfg.Validate(incomingKey)
			if !valid {
				res := response.ResponseGenerate(
					http.StatusUnauthorized,
					errors.New("invalid or missing API key"),
					nil,
					nil,
					nil,
				)
				return response.HttpResponses(c, res)
			}

			c.Set("api_client", clientName)
			return next(c)
		}
	}
}

