package middleware

import (
	"errors"
	"go-template/pkg/config"
	"go-template/pkg/response"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func JWTAuth(jwtCfg config.JwtCnfg) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get(echo.HeaderAuthorization)
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				res := response.ResponseGenerate(http.StatusUnauthorized, errors.New("missing or malformed bearer token"), nil, nil, nil)
				return response.HttpResponses(c, res)
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := jwtCfg.VerifyToken(tokenString)
			if err != nil {
				res := response.ResponseGenerate(http.StatusUnauthorized, err, nil, nil, nil)
				return response.HttpResponses(c, res)
			}

			c.Set("user_claims", claims)
			if username, ok := claims["username"].(string); ok {
				c.Set("username", username)
			}

			return next(c)
		}
	}
}

