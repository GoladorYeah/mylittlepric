package handlers

import (
	"log"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"

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

	// Directly call the pre-created handler
	err := h.handler(c)

	// Log response status and any errors
	statusCode := c.Response().StatusCode()
	log.Printf("📊 Metrics handler finished - status: %d, error: %v", statusCode, err)

	if err != nil {
		log.Printf("❌ Error from metrics handler: %v", err)
		// Also log the response body if there's an error
		bodyBytes := c.Response().Body()
		if len(bodyBytes) > 0 {
			log.Printf("❌ Response body: %s", string(bodyBytes[:min(len(bodyBytes), 500)]))
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
