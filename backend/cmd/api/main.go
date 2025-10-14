package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"mylittleprice/internal/config"
	"mylittleprice/internal/container"
	"mylittleprice/internal/handlers"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	c, err := container.NewContainer(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize container: %v", err)
	}
	defer c.Close()

	app := fiber.New(fiber.Config{
		AppName:      "MyLittlePrice API",
		ServerHeader: "Fiber",
		ErrorHandler: customErrorHandler,
	})

	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format:     "${time} | ${status} | ${latency} | ${method} ${path}\n",
		TimeFormat: "15:04:05",
	}))

	app.Use(cors.New(cors.Config{
		AllowOrigins:     joinStrings(cfg.CORSOrigins, ","),
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: true,
	}))

	app.Get("/health", func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{
			"status":    "ok",
			"timestamp": time.Now(),
			"services":  c.HealthCheck(),
		})
	})

	setupRoutes(app, c)

	port := cfg.Port
	log.Printf("🚀 Server starting on port %s", port)
	log.Printf("🔒 Environment: %s", cfg.Env)
	log.Printf("🌍 Allowed origins: %v", cfg.CORSOrigins)

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

	if err := app.Listen(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}

func setupRoutes(app *fiber.App, c *container.Container) {
	api := app.Group("/api")

	wsHandler := handlers.NewWSHandler(c)

	app.Use("/ws", func(ctx *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(ctx) {
			ctx.Locals("allowed", true)
			return ctx.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws", websocket.New(func(conn *websocket.Conn) {
		wsHandler.HandleWebSocket(conn)
	}))

	chatHandler := handlers.NewChatHandler(c)
	api.Post("/chat", chatHandler.HandleChat)

	productHandler := handlers.NewProductHandler(c)
	api.Post("/product-details", productHandler.HandleProductDetails)

	setupStatsRoutes(api, c)
}

func setupStatsRoutes(api fiber.Router, c *container.Container) {
	api.Get("/stats/keys", func(ctx *fiber.Ctx) error {
		geminiStats, _ := c.GeminiRotator.GetAllStats()
		serpStats, _ := c.SerpRotator.GetAllStats()

		return ctx.JSON(fiber.Map{
			"gemini": geminiStats,
			"serp":   serpStats,
		})
	})

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

	api.Get("/stats/tokens", func(ctx *fiber.Ctx) error {
		tokenStats := c.GeminiService.GetTokenStats()

		return ctx.JSON(fiber.Map{
			"token_usage": tokenStats,
			"timestamp":   time.Now(),
		})
	})

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
