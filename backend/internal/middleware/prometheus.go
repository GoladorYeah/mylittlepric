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
	log.Printf("🔵 RegisterMetrics() called")
	metricsOnce.Do(func() {
		log.Printf("🔧 Registering HTTP middleware metrics (inside sync.Once)")

		httpRequestsTotal = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"method", "handler", "status"},
		)
		log.Printf("📍 httpRequestsTotal created at %p", httpRequestsTotal)
		prometheus.MustRegister(httpRequestsTotal)
		log.Printf("✅ Registered http_requests_total at %p", httpRequestsTotal)

		httpRequestDuration = prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "HTTP request latencies in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "handler"},
		)
		prometheus.MustRegister(httpRequestDuration)
		log.Printf("✅ Registered http_request_duration_seconds")

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

		// Debug: Check what's registered in DefaultRegistry
		metrics, err := prometheus.DefaultGatherer.Gather()
		if err != nil {
			log.Printf("❌ Failed to gather metrics: %v", err)
		} else {
			log.Printf("📊 Total metric families in DefaultRegistry: %d", len(metrics))
			httpRequestsCount := 0
			httpDurationCount := 0
			for _, m := range metrics {
				if m.GetName() == "http_requests_total" {
					httpRequestsCount++
					log.Printf("  - http_requests_total [#%d]: %d metrics", httpRequestsCount, len(m.GetMetric()))
				}
				if m.GetName() == "http_request_duration_seconds" {
					httpDurationCount++
					log.Printf("  - http_request_duration_seconds [#%d]: %d metrics", httpDurationCount, len(m.GetMetric()))
				}
			}
			if httpRequestsCount > 1 {
				log.Printf("⚠️ WARNING: http_requests_total appears %d times in registry!", httpRequestsCount)
			}
			if httpDurationCount > 1 {
				log.Printf("⚠️ WARNING: http_request_duration_seconds appears %d times in registry!", httpDurationCount)
			}
		}
	})
}

// PrometheusMiddleware creates a Fiber middleware that collects HTTP metrics
func PrometheusMiddleware() fiber.Handler {
	log.Printf("🟢 PrometheusMiddleware() called - creating new handler")
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

		// CRITICAL FIX: Make immutable copies of Fiber strings before using as Prometheus labels
		// Fiber reuses byte slices for performance, which causes race conditions with Prometheus
		// See: https://github.com/prometheus/client_golang/issues/1269
		method := string([]byte(c.Method()))
		pathCopy := string([]byte(path))
		status := strconv.Itoa(c.Response().StatusCode())

		// Debug logging for /api/user/preferences
		if pathCopy == "/api/user/preferences" && method == "PUT" {
			log.Printf("📊 Recording metrics: method=%s, path=%s, status=%s, duration=%.4fs, request_id=%p", method, pathCopy, status, duration, c)
		}

		// Record metrics
		log.Printf("🔹 About to record: method=%s, path=%s, status=%s (req=%p)", method, pathCopy, status, c)
		httpRequestsTotal.WithLabelValues(method, pathCopy, status).Inc()
		httpRequestDuration.WithLabelValues(method, pathCopy).Observe(duration)
		log.Printf("✅ Recorded: method=%s, path=%s, status=%s (req=%p)", method, pathCopy, status, c)

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
