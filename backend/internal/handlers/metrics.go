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
	// Create the Prometheus handler once at initialization
	promHandler := promhttp.Handler()
	fiberHandler := adaptor.HTTPHandler(promHandler)

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
	utils.LogDebug(c.Context(), "📊 Metrics requested")

	// Directly call the pre-created handler
	err := h.handler(c)

	if err != nil {
		log.Printf("❌ Error from metrics handler: %v", err)
	}

	return err
}
