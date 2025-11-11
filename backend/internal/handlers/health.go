package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"mylittleprice/internal/container"
)

type HealthHandler struct {
	container *container.Container
}

func NewHealthHandler(c *container.Container) *HealthHandler {
	return &HealthHandler{container: c}
}

type HealthResponse struct {
	Status  string           `json:"status"`
	Checks  map[string]Check `json:"checks"`
	Version string           `json:"version"`
	Uptime  int64            `json:"uptime_seconds"`
}

type Check struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

// Liveness - простая проверка для Kubernetes liveness probe
// GET /health/live
func (h *HealthHandler) Liveness(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status": "ok",
	})
}

// Readiness - проверка зависимостей для Kubernetes readiness probe
// GET /health/ready
func (h *HealthHandler) Readiness(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	checks := make(map[string]Check)
	healthy := true

	// Check PostgreSQL
	if err := h.container.Ent.DB().PingContext(ctx); err != nil {
		checks["postgresql"] = Check{Status: "unhealthy", Message: err.Error()}
		healthy = false
	} else {
		checks["postgresql"] = Check{Status: "healthy"}
	}

	// Check Redis
	if err := h.container.Redis.Ping(ctx).Err(); err != nil {
		checks["redis"] = Check{Status: "unhealthy", Message: err.Error()}
		healthy = false
	} else {
		checks["redis"] = Check{Status: "healthy"}
	}

	status := "ok"
	statusCode := fiber.StatusOK
	if !healthy {
		status = "degraded"
		statusCode = fiber.StatusServiceUnavailable
	}

	return c.Status(statusCode).JSON(HealthResponse{
		Status:  status,
		Checks:  checks,
		Version: "1.0.0",
		Uptime:  int64(time.Since(h.container.StartTime).Seconds()),
	})
}

// Health - детальная информация о здоровье системы
// GET /health
func (h *HealthHandler) Health(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	checks := make(map[string]Check)
	healthy := true

	// Check PostgreSQL
	if err := h.container.Ent.DB().PingContext(ctx); err != nil {
		checks["postgresql"] = Check{Status: "unhealthy", Message: err.Error()}
		healthy = false
	} else {
		checks["postgresql"] = Check{Status: "healthy"}
	}

	// Check Redis
	if err := h.container.Redis.Ping(ctx).Err(); err != nil {
		checks["redis"] = Check{Status: "unhealthy", Message: err.Error()}
		healthy = false
	} else {
		poolStats := h.container.Redis.PoolStats()
		checks["redis"] = Check{
			Status:  "healthy",
			Message: fmt.Sprintf("pool: %d active, %d idle", poolStats.TotalConns-poolStats.IdleConns, poolStats.IdleConns),
		}
	}

	status := "ok"
	if !healthy {
		status = "degraded"
	}

	return c.JSON(HealthResponse{
		Status:  status,
		Checks:  checks,
		Version: "1.0.0",
		Uptime:  int64(time.Since(h.container.StartTime).Seconds()),
	})
}
