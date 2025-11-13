package handlers

import (
	"log"
	"runtime/debug"

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
	// This is more efficient and avoids potential issues with recreating handlers
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
	// Wrap handler with panic recovery and detailed logging
	defer func() {
		if r := recover(); r != nil {
			stack := debug.Stack()
			log.Printf("❌ Panic in metrics handler: %v\nStack trace:\n%s", r, string(stack))
			c.Status(fiber.StatusInternalServerError).SendString("Error collecting metrics")
		}
	}()

	// Use adaptor to convert standard http.Handler to fiber handler
	handler := adaptor.HTTPHandler(promhttp.Handler())
	err := handler(c)

	if err != nil {
		log.Printf("❌ Error from metrics handler: %v", err)
		return err
	}

	return nil
}
