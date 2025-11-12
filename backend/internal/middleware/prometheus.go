package middleware

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Prometheus metrics
var (
	// HTTP request metrics
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "handler", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request latencies in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "handler"},
	)

	httpRequestsInFlight = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_requests_in_flight",
			Help: "Current number of HTTP requests being processed",
		},
	)

	// Rate limiting metrics
	rateLimitExceeded = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "rate_limit_exceeded_total",
			Help: "Total number of requests that exceeded rate limit",
		},
		[]string{"endpoint"},
	)

	rateLimiterRedisErrors = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "rate_limiter_redis_errors_total",
			Help: "Total number of Redis errors in rate limiter",
		},
	)
)

// PrometheusMiddleware creates a Fiber middleware that collects HTTP metrics
func PrometheusMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Increment in-flight requests
		httpRequestsInFlight.Inc()
		defer httpRequestsInFlight.Dec()

		// Process request
		err := c.Next()

		// Calculate duration
		duration := time.Since(start).Seconds()

		// Get route path (or use raw path if route not found)
		path := c.Route().Path
		if path == "" {
			path = c.Path()
		}

		method := c.Method()
		status := strconv.Itoa(c.Response().StatusCode())

		// Record metrics
		httpRequestsTotal.WithLabelValues(method, path, status).Inc()
		httpRequestDuration.WithLabelValues(method, path).Observe(duration)

		return err
	}
}

// RecordRateLimitExceeded increments the rate limit counter for an endpoint
func RecordRateLimitExceeded(endpoint string) {
	rateLimitExceeded.WithLabelValues(endpoint).Inc()
}

// RecordRateLimiterRedisError increments the Redis error counter
func RecordRateLimiterRedisError() {
	rateLimiterRedisErrors.Inc()
}
