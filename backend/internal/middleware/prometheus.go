package middleware

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"mylittleprice/internal/metrics"
)

// PrometheusMiddleware собирает метрики HTTP запросов
func PrometheusMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Process request
		err := c.Next()

		// Record metrics
		duration := time.Since(start).Seconds()
		status := c.Response().StatusCode()

		metrics.HTTPRequestsTotal.WithLabelValues(
			c.Method(),
			c.Path(),
			fmt.Sprintf("%d", status),
		).Inc()

		metrics.HTTPRequestDuration.WithLabelValues(
			c.Method(),
			c.Path(),
		).Observe(duration)

		return err
	}
}
