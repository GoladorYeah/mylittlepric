package handlers

import (
	"log"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"mylittleprice/internal/utils"
)

// MetricsHandler handles Prometheus metrics endpoint
type MetricsHandler struct {
	// Pre-created handler to avoid recreating on every request
	handler fiber.Handler
}

// NewMetricsHandler creates a new metrics handler
func NewMetricsHandler() *MetricsHandler {
	log.Printf("🔧 Creating new MetricsHandler")
	// Create the Prometheus handler once at initialization
	promHandler := promhttp.Handler()
	fiberHandler := adaptor.HTTPHandler(promHandler)

	log.Printf("✅ MetricsHandler created successfully")
	return &MetricsHandler{
		handler: fiberHandler,
	}
}

// GetMetrics returns Prometheus metrics
// @Summary Get Prometheus metrics
// @Description Returns Prometheus metrics in text format
// @Tags monitoring
// @Produce plain
// @Success 200 {string} string "Prometheus metrics"
// @Router /metrics [get]
func (h *MetricsHandler) GetMetrics(c *fiber.Ctx) error {
	log.Printf("📊 Metrics requested - starting handler")
	utils.LogDebug(c.Context(), "📊 Metrics requested")

	// Debug: Check metrics before gathering
	metrics, err := prometheus.DefaultGatherer.Gather()
	if err != nil {
		log.Printf("❌ Failed to pre-gather metrics: %v", err)
	} else {
		log.Printf("📊 Pre-gather: Total metric families: %d", len(metrics))
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
			log.Printf("⚠️ PRE-GATHER WARNING: http_requests_total appears %d times!", httpRequestsCount)
		}
		if httpDurationCount > 1 {
			log.Printf("⚠️ PRE-GATHER WARNING: http_request_duration_seconds appears %d times!", httpDurationCount)
		}
	}

	// Directly call the pre-created handler
	err = h.handler(c)

	// Log response status and any errors
	statusCode := c.Response().StatusCode()
	log.Printf("📊 Metrics handler finished - status: %d, error: %v", statusCode, err)

	// Log response body if status is 500 (even when err is nil, promhttp sets 500 on collection errors)
	if statusCode >= 500 || err != nil {
		if err != nil {
			log.Printf("❌ Error from metrics handler: %v", err)
		}
		bodyBytes := c.Response().Body()
		if len(bodyBytes) > 0 {
			log.Printf("❌ Response body (first 1000 chars): %s", string(bodyBytes[:min(len(bodyBytes), 1000)]))
		} else {
			log.Printf("⚠️ Response body is empty despite status %d", statusCode)
		}
	}

	return err
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
