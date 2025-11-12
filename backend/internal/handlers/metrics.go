package handlers

import (
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// MetricsHandler handles Prometheus metrics endpoint
type MetricsHandler struct{}

// NewMetricsHandler creates a new metrics handler
func NewMetricsHandler() *MetricsHandler {
	return &MetricsHandler{}
}

// GetMetrics returns Prometheus metrics
// @Summary Get Prometheus metrics
// @Description Returns Prometheus metrics in text format
// @Tags monitoring
// @Produce plain
// @Success 200 {string} string "Prometheus metrics"
// @Router /metrics [get]
func (h *MetricsHandler) GetMetrics(c *fiber.Ctx) error {
	handler := adaptor.HTTPHandler(promhttp.Handler())
	return handler(c)
}
