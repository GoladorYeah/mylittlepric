package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"mylittleprice/internal/app"
	"mylittleprice/internal/config"
	"mylittleprice/internal/container"
	"mylittleprice/internal/handlers"
	"mylittleprice/internal/jobs"
	"mylittleprice/internal/middleware"
	"mylittleprice/internal/utils"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize structured logger
	utils.InitLogger(cfg.LogLevel, cfg.LogFormat, cfg.LokiEnabled, cfg.LokiURL, cfg.LokiServiceName)
	logger := utils.GetLogger()
	ctx := context.Background()

	logger.Info("Starting MyLittlePrice Backend",
		slog.String("version", "1.0.0"),
		slog.String("env", cfg.Env),
		slog.String("log_level", cfg.LogLevel),
		slog.String("log_format", cfg.LogFormat),
		slog.Bool("loki_enabled", cfg.LokiEnabled),
		slog.String("loki_url", cfg.LokiURL),
	)

	c, err := container.NewContainer(cfg)
	if err != nil {
		logger.Error("Failed to initialize container", err)
		os.Exit(1)
	}
	defer c.Close()

	logger.Info("Container initialized successfully")

	// Initialize and start cleanup job
	cleanupJob := jobs.NewCleanupJob(c.SearchHistoryService)
	cleanupJob.Start()
	defer cleanupJob.Stop()

	logger.Info("Cleanup job started")

	fiberApp := fiber.New(fiber.Config{
		AppName:      "MyLittlePrice API",
		ServerHeader: "Fiber",
		ErrorHandler: customErrorHandler,
	})

	fiberApp.Use(recover.New())
	fiberApp.Use(fiberlogger.New(fiberlogger.Config{
		Format:     "${time} | ${status} | ${latency} | ${method} ${path}\n",
		TimeFormat: "15:04:05",
	}))
	fiberApp.Use(middleware.RequestContext())

	// CORS configuration with dynamic origin validation
	fiberApp.Use(cors.New(cors.Config{
		AllowOriginsFunc: func(origin string) bool {
			// Check if origin is in the allowed list
			for _, allowedOrigin := range cfg.CORSOrigins {
				if origin == allowedOrigin {
					// Log successful CORS validation in development
					if cfg.Env == "development" {
						logger.Debug("CORS allowed origin", slog.String("origin", origin))
					}
					return true
				}
			}

			// Always log rejected origins for debugging (even in production)
			logger.Warn("CORS rejected - origin not in allowed list",
				slog.String("origin", origin),
				slog.Any("allowed_origins", cfg.CORSOrigins),
			)

			return false
		},
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: true,
		MaxAge:           86400, // 24 hours for preflight cache
	}))

	// Health check endpoints
	healthHandler := handlers.NewHealthHandler(c)
	fiberApp.Get("/health/live", healthHandler.Liveness)
	fiberApp.Get("/health/ready", healthHandler.Readiness)
	fiberApp.Get("/health", healthHandler.Health)

	logger.Info("Health check endpoints initialized",
		slog.String("health_endpoint", "/health"),
	)

	app.SetupRoutes(fiberApp, c)

	port := cfg.Port
	logger.Info("Server starting",
		slog.String("port", port),
		slog.String("env", cfg.Env),
		slog.Any("allowed_origins", cfg.CORSOrigins),
	)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-quit
		logger.Info("Shutting down server...")

		// Stop cleanup job first
		cleanupJob.Stop()

		if err := fiberApp.Shutdown(); err != nil {
			utils.LogError(ctx, "Server shutdown error", err)
		}

		// Close Loki writer to flush remaining logs
		if err := utils.CloseLoki(); err != nil {
			logger.Error("Failed to close Loki writer", err)
		}

		logger.Info("Server stopped gracefully")
	}()

	if err := fiberApp.Listen(fmt.Sprintf(":%s", port)); err != nil {
		logger.Error("Failed to start server", err)
		os.Exit(1)
	}
}

func customErrorHandler(ctx *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	return ctx.Status(code).JSON(fiber.Map{
		"error":   true,
		"message": message,
		"code":    code,
	})
}
