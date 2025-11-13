package handlers

import (
	"log/slog"

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
func (h *MetricsHandler) GetMetrics(c *fiber.Ctx) (err error) {
	logger := utils.GetLogger()

	// Add panic recovery to prevent crashes
	defer func() {
		if r := recover(); r != nil {
			logger.Error("Panic recovered in metrics endpoint",
				slog.Any("panic", r),
				slog.String("path", c.Path()),
			)
			err = fiber.NewError(fiber.StatusInternalServerError, "Failed to generate metrics")
		}
	}()

	// Call the pre-created handler
	err = h.handler(c)
	if err != nil {
		logger.Error("Error serving metrics",
			slog.Any("error", err),
			slog.String("path", c.Path()),
		)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to serve metrics")
	}

	return nil
}
