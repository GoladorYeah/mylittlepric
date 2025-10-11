package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"mylittleprice/internal/config"
	"mylittleprice/internal/container"
	"mylittleprice/internal/handlers"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize dependency container
	c, err := container.NewContainer(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize container: %v", err)
	}
	defer c.Close()

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "MyLittlePrice API",
		ServerHeader: "Fiber",
		ErrorHandler: customErrorHandler,
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format:     "${time} | ${status} | ${latency} | ${method} ${path}\n",
		TimeFormat: "15:04:05",
	}))

	// CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins:     joinStrings(cfg.CORSOrigins, ","),
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: true,
	}))

	// Health check
	app.Get("/health", func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{
			"status":    "ok",
			"timestamp": time.Now(),
			"services":  c.HealthCheck(),
		})
	})

	// Setup routes
	setupRoutes(app, c)

	// Start server
	port := cfg.Port
	log.Printf("🚀 Server starting on port %s", port)
	log.Printf("🔒 Environment: %s", cfg.Env)
	log.Printf("🌍 Allowed origins: %v", cfg.CORSOrigins)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("🛑 Shutting down server...")

		if err := app.Shutdown(); err != nil {
			log.Printf("❌ Server shutdown error: %v", err)
		}

		log.Println("✅ Server stopped gracefully")
	}()

	// Start server
	if err := app.Listen(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}

// setupRoutes configures all application routes
func setupRoutes(app *fiber.App, c *container.Container) {
	api := app.Group("/api")

	// Chat routes
	chatHandler := handlers.NewChatHandler(c)
	api.Post("/chat", chatHandler.HandleChat)

	// Product details
	productHandler := handlers.NewProductHandler(c)
	api.Post("/product-details", productHandler.HandleProductDetails)

	// Statistics endpoints
	setupStatsRoutes(api, c)
}

// setupStatsRoutes configures statistics endpoints
func setupStatsRoutes(api fiber.Router, c *container.Container) {
	// Key rotation stats
	api.Get("/stats/keys", func(ctx *fiber.Ctx) error {
		geminiStats, _ := c.GeminiRotator.GetAllStats()
		serpStats, _ := c.SerpRotator.GetAllStats()

		return ctx.JSON(fiber.Map{
			"gemini": geminiStats,
			"serp":   serpStats,
		})
	})

	// Grounding statistics
	api.Get("/stats/grounding", func(ctx *fiber.Ctx) error {
		stats := c.GeminiService.GetGroundingStats()

		groundingPercentage := float32(0)
		if stats.TotalDecisions > 0 {
			groundingPercentage = float32(stats.GroundingEnabled) / float32(stats.TotalDecisions) * 100
		}

		return ctx.JSON(fiber.Map{
			"total_decisions":      stats.TotalDecisions,
			"grounding_enabled":    stats.GroundingEnabled,
			"grounding_disabled":   stats.GroundingDisabled,
			"grounding_percentage": fmt.Sprintf("%.1f%%", groundingPercentage),
			"reason_breakdown":     stats.ReasonCounts,
			"average_confidence":   fmt.Sprintf("%.2f", stats.AverageConfidence),
			"mode":                 c.Config.GeminiGroundingMode,
			"config": fiber.Map{
				"enabled":   c.Config.GeminiUseGrounding,
				"min_words": c.Config.GeminiGroundingMinWords,
			},
		})
	})

	// Token usage statistics
	api.Get("/stats/tokens", func(ctx *fiber.Ctx) error {
		tokenStats := c.GeminiService.GetTokenStats()

		return ctx.JSON(fiber.Map{
			"token_usage": tokenStats,
			"timestamp":   time.Now(),
		})
	})

	// Combined statistics
	api.Get("/stats/all", func(ctx *fiber.Ctx) error {
		geminiStats, _ := c.GeminiRotator.GetAllStats()
		serpStats, _ := c.SerpRotator.GetAllStats()
		groundingStats := c.GeminiService.GetGroundingStats()
		tokenStats := c.GeminiService.GetTokenStats()

		groundingPercentage := float32(0)
		if groundingStats.TotalDecisions > 0 {
			groundingPercentage = float32(groundingStats.GroundingEnabled) / float32(groundingStats.TotalDecisions) * 100
		}

		return ctx.JSON(fiber.Map{
			"api_keys": fiber.Map{
				"gemini": geminiStats,
				"serp":   serpStats,
			},
			"grounding": fiber.Map{
				"total_decisions":      groundingStats.TotalDecisions,
				"grounding_enabled":    groundingStats.GroundingEnabled,
				"grounding_disabled":   groundingStats.GroundingDisabled,
				"grounding_percentage": fmt.Sprintf("%.1f%%", groundingPercentage),
				"reason_breakdown":     groundingStats.ReasonCounts,
				"average_confidence":   fmt.Sprintf("%.2f", groundingStats.AverageConfidence),
				"mode":                 c.Config.GeminiGroundingMode,
			},
			"tokens":    tokenStats,
			"timestamp": time.Now(),
		})
	})
}

// customErrorHandler handles all errors
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

// joinStrings joins string slice with separator
func joinStrings(slice []string, sep string) string {
	result := ""
	for i, s := range slice {
		if i > 0 {
			result += sep
		}
		result += s
	}
	return result
}
