package middleware

import (
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
)

// Prometheus metrics
var (
	// HTTP request metrics
	httpRequestsTotal      *prometheus.CounterVec
	httpRequestDuration    *prometheus.HistogramVec
	httpRequestsInFlight   prometheus.Gauge
	rateLimitExceeded      *prometheus.CounterVec
	rateLimiterRedisErrors prometheus.Counter

	// Ensure metrics are registered only once
	metricsOnce sync.Once
	metricsRegistered bool
)

// RegisterMetrics registers all HTTP middleware metrics to default registry
// Called explicitly to avoid issues with init() being called multiple times
func RegisterMetrics() {
	metricsOnce.Do(func() {
		log.Printf("🔧 Registering HTTP middleware metrics")

		httpRequestsTotal = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"method", "handler", "status"},
		)
		prometheus.MustRegister(httpRequestsTotal)

		httpRequestDuration = prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "HTTP request latencies in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "handler"},
		)
		prometheus.MustRegister(httpRequestDuration)

		httpRequestsInFlight = prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "http_requests_in_flight",
				Help: "Current number of HTTP requests being processed",
			},
		)
		prometheus.MustRegister(httpRequestsInFlight)

		rateLimitExceeded = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "rate_limit_exceeded_total",
				Help: "Total number of requests that exceeded rate limit",
			},
			[]string{"endpoint"},
		)
		prometheus.MustRegister(rateLimitExceeded)

		rateLimiterRedisErrors = prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "rate_limiter_redis_errors_total",
				Help: "Total number of Redis errors in rate limiter",
			},
		)
		prometheus.MustRegister(rateLimiterRedisErrors)

		metricsRegistered = true
		log.Printf("✅ HTTP middleware metrics registered successfully")
	})
}

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

		// Debug logging for /api/user/preferences
		if path == "/api/user/preferences" && method == "PUT" {
			log.Printf("📊 Recording metrics: method=%s, path=%s, status=%s, duration=%.4fs", method, path, status, duration)
		}

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
