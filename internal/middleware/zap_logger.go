package middleware

import (
	"go-template/pkg/logger"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func ZapLogger(zapLog *zap.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			req := c.Request()
			res := c.Response()

			requestID := req.Header.Get(echo.HeaderXRequestID)
			if requestID == "" {
				requestID = res.Header().Get(echo.HeaderXRequestID)
			}

			// Inject request_id into standard context.Context so it flows to Service & Repo
			ctx := logger.WithRequestID(req.Context(), requestID)
			req = req.WithContext(ctx)
			c.SetRequest(req)

			err := next(c)
			stop := time.Now()

			log := logger.FromContext(ctx, zapLog)
			log.Info("HTTP request",
				zap.Int("status", res.Status),
				zap.String("method", req.Method),
				zap.String("path", req.URL.Path),
				zap.String("query", req.URL.RawQuery),
				zap.String("client_ip", c.RealIP()),
				zap.Duration("latency", stop.Sub(start)),
			)

			return err
		}
	}
}
