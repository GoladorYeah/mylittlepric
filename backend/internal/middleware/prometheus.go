package middleware

import (
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
)

// Prometheus metrics - using default registry
var (
	// HTTP request metrics
	httpRequestsTotal *prometheus.CounterVec
	httpRequestDuration *prometheus.HistogramVec
	httpRequestsInFlight prometheus.Gauge

	// Rate limiting metrics
	rateLimitExceeded *prometheus.CounterVec
	rateLimiterRedisErrors prometheus.Counter
)

// init registers metrics once
func init() {
	// Try to register metrics, ignore if already registered
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "handler", "status"},
	)
	if err := prometheus.Register(httpRequestsTotal); err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			httpRequestsTotal = are.ExistingCollector.(*prometheus.CounterVec)
			log.Printf("⚠️ http_requests_total already registered, reusing existing")
		} else {
			log.Printf("❌ Failed to register http_requests_total: %v", err)
		}
	}

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request latencies in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "handler"},
	)
	if err := prometheus.Register(httpRequestDuration); err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			httpRequestDuration = are.ExistingCollector.(*prometheus.HistogramVec)
			log.Printf("⚠️ http_request_duration_seconds already registered, reusing existing")
		} else {
			log.Printf("❌ Failed to register http_request_duration_seconds: %v", err)
		}
	}

	httpRequestsInFlight = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_requests_in_flight",
			Help: "Current number of HTTP requests being processed",
		},
	)
	if err := prometheus.Register(httpRequestsInFlight); err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			httpRequestsInFlight = are.ExistingCollector.(prometheus.Gauge)
			log.Printf("⚠️ http_requests_in_flight already registered, reusing existing")
		} else {
			log.Printf("❌ Failed to register http_requests_in_flight: %v", err)
		}
	}

	rateLimitExceeded = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "rate_limit_exceeded_total",
			Help: "Total number of requests that exceeded rate limit",
		},
		[]string{"endpoint"},
	)
	if err := prometheus.Register(rateLimitExceeded); err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			rateLimitExceeded = are.ExistingCollector.(*prometheus.CounterVec)
			log.Printf("⚠️ rate_limit_exceeded_total already registered, reusing existing")
		} else {
			log.Printf("❌ Failed to register rate_limit_exceeded_total: %v", err)
		}
	}

	rateLimiterRedisErrors = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "rate_limiter_redis_errors_total",
			Help: "Total number of Redis errors in rate limiter",
		},
	)
	if err := prometheus.Register(rateLimiterRedisErrors); err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			rateLimiterRedisErrors = are.ExistingCollector.(prometheus.Counter)
			log.Printf("⚠️ rate_limiter_redis_errors_total already registered, reusing existing")
		} else {
			log.Printf("❌ Failed to register rate_limiter_redis_errors_total: %v", err)
		}
	}
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
