package middleware

import (
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests processed, partitioned by method, path, and status code.",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Latency of HTTP requests in seconds, partitioned by method and path.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)
)

// PrometheusMetrics returns an Echo middleware that records HTTP metrics for Prometheus.
func PrometheusMetrics() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Path() == "/metrics" {
				return next(c)
			}

			start := time.Now()
			err := next(c)

			path := c.Path()
			if path == "" {
				path = "unmatched"
			}
			method := c.Request().Method
			status := strconv.Itoa(c.Response().Status)

			httpRequestsTotal.WithLabelValues(method, path, status).Inc()
			httpRequestDuration.WithLabelValues(method, path).Observe(time.Since(start).Seconds())

			return err
		}
	}
}

