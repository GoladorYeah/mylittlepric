package middleware

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"mylittleprice/internal/metrics"
)

// PrometheusMiddleware records HTTP request metrics for Prometheus
func PrometheusMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Process request
		err := c.Next()

		// Record metrics
		duration := time.Since(start).Seconds()
		status := c.Response().StatusCode()
		method := c.Method()
		path := c.Path()

		// Record request count with labels
		metrics.HTTPRequestsTotal.WithLabelValues(
			method,
			path,
			fmt.Sprintf("%d", status),
		).Inc()

		// Record request duration
		metrics.HTTPRequestDuration.WithLabelValues(
			method,
			path,
		).Observe(duration)

		return err
	}
}
