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
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

type Controller struct {
	dig.In

	User     user.Controller
	Auth     auth.Controller
	Config   config.Config
	Database model.Database
	Redis    model.Redis
	Logger   *zap.Logger
}

func (c *Controller) Routes(ec *echo.Echo) {
	ec.Use(echoMiddleware.Recover())
	ec.Use(echoMiddleware.RequestID())
	ec.Use(echoMiddleware.CORS())
	ec.Use(middleware.ZapLogger(c.Logger))
	ec.Use(middleware.PrometheusMetrics())

	// Register Swagger UI (/swagger) and OpenAPI YAML (/docs/openapi.yaml)
	docs.RegisterRoutes(ec)

	// Expose Prometheus Metrics (/metrics)
	ec.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	ec.GET("/", func(eCtx echo.Context) error {
		return eCtx.String(http.StatusOK, "Go API is up and running!")
	})

	// General Health
	ec.GET("/health", func(eCtx echo.Context) error {
		return eCtx.JSON(http.StatusOK, echo.Map{"status": "healthy"})
	})

	// K8s Liveness Probe: process is alive
	ec.GET("/health/liveness", func(eCtx echo.Context) error {
		return eCtx.JSON(http.StatusOK, echo.Map{"status": "alive"})
	})

	// K8s Readiness Probe: dependencies (PostgreSQL, Redis) are healthy
	ec.GET("/health/readiness", func(eCtx echo.Context) error {
		ctx := eCtx.Request().Context()
		checks := echo.Map{}
		isHealthy := true

		// Check PostgreSQL
		sqlDB, err := c.Database.PostgreSql.DB()
		if err != nil || sqlDB.PingContext(ctx) != nil {
			isHealthy = false
			if err != nil {
				checks["database"] = err.Error()
			} else {
				checks["database"] = sqlDB.PingContext(ctx).Error()
			}
		} else {
			checks["database"] = "up"
		}

		// Check Redis
		if err := c.Redis.Ping(ctx); err != nil {
			isHealthy = false
			checks["redis"] = err.Error()
		} else {
			checks["redis"] = "up"
		}

		if !isHealthy {
			return eCtx.JSON(http.StatusServiceUnavailable, echo.Map{
				"status": "unhealthy",
				"checks": checks,
			})
		}

		return eCtx.JSON(http.StatusOK, echo.Map{
			"status": "ready",
			"checks": checks,
		})
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
