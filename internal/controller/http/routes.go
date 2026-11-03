package http

import (
	"go-template/docs"
	"go-template/internal/controller/http/auth"
	"go-template/internal/controller/http/user"
	"go-template/internal/middleware"
	"go-template/pkg/config"
	"go-template/pkg/resource/model"
	"net/http"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

type Controller struct {
	dig.In

	User     user.Controller
	Auth     auth.Controller
	Config   config.Config
	Database model.Database
	Logger   *zap.Logger
}

func (c *Controller) Routes(ec *echo.Echo) {
	ec.Use(echoMiddleware.Recover())
	ec.Use(echoMiddleware.RequestID())
	ec.Use(echoMiddleware.CORS())
	ec.Use(middleware.ZapLogger(c.Logger))

	// Register Swagger UI (/swagger) and OpenAPI YAML (/docs/openapi.yaml)
	docs.RegisterRoutes(ec)

	ec.GET("/", func(eCtx echo.Context) error {
		return eCtx.String(http.StatusOK, "Go API is up and running!")
	})

	ec.GET("/health", func(eCtx echo.Context) error {
		sqlDB, err := c.Database.PostgreSql.DB()
		if err != nil {
			return eCtx.JSON(http.StatusServiceUnavailable, echo.Map{"status": "unhealthy", "error": err.Error()})
		}
		if err := sqlDB.PingContext(eCtx.Request().Context()); err != nil {
			return eCtx.JSON(http.StatusServiceUnavailable, echo.Map{"status": "unhealthy", "error": err.Error()})
		}
		return eCtx.JSON(http.StatusOK, echo.Map{"status": "healthy"})
	})

	v1 := ec.Group("/v1")

	// 1. Public routes (No Auth)
	v1.POST("/register", c.User.CreateUser)
	v1.POST("/login", c.Auth.LoginFunction)

	// 2. JWT Protected routes (Authorization: Bearer <token>)
	jwtProtected := v1.Group("")
	jwtProtected.Use(middleware.JWTAuth(c.Config.JWT))
	jwtProtected.GET("/users", c.User.GetUsers)
	jwtProtected.POST("/logout", c.Auth.LogoutFunction)

	// 3. API Key Protected routes (X-API-Key: <key>) - Service-to-Service / Internal APIs
	internalGroup := v1.Group("/internal")
	internalGroup.Use(middleware.APIKeyAuth(c.Config.APIKey))
	internalGroup.GET("/users", c.User.GetUsers)
}
