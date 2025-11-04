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
	"mylittleprice/internal/middleware"
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

	// Authentication routes (public)
	setupAuthRoutes(api, c)

	// WebSocket chat (REQUIRES authentication)
	wsHandler := handlers.NewWSHandler(c)
	authMiddleware := middleware.AuthMiddleware(c.JWTService)

	app.Use("/ws", func(ctx *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(ctx) {
			// Try to get token from query parameter or Authorization header
			var token string

			// First, try query parameter (for WebSocket compatibility)
			queryToken := ctx.Query("token")
			if queryToken != "" {
				token = queryToken
			} else {
				// Fallback to Authorization header
				authHeader := ctx.Get("Authorization")
				if authHeader != "" && len(authHeader) > 7 {
					token = authHeader[7:] // Remove "Bearer " prefix
				}
			}

			if token == "" {
				return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "authentication required",
				})
			}

			// Validate token
			claims, err := c.JWTService.ValidateAccessToken(token)
			if err != nil {
				return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "invalid or expired token",
				})
			}

			// Store user info in locals for WebSocket handler
			ctx.Locals("user_id", claims.UserID)
			ctx.Locals("user_email", claims.Email)
			ctx.Locals("allowed", true)
			return ctx.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws", websocket.New(func(conn *websocket.Conn) {
		wsHandler.HandleWebSocket(conn)
	}))

	// Chat endpoints (REQUIRE authentication)
	chatHandler := handlers.NewChatHandler(c)
	api.Post("/chat", authMiddleware, chatHandler.HandleChat)
	api.Get("/chat/messages", authMiddleware, chatHandler.GetSessionMessages)

	productHandler := handlers.NewProductHandler(c)
	api.Post("/product-details", productHandler.HandleProductDetails)

	// Search history routes (REQUIRE authentication)
	setupSearchHistoryRoutes(api, c, authMiddleware)

	setupStatsRoutes(api, c)
}

func setupAuthRoutes(api fiber.Router, c *container.Container) {
	auth := api.Group("/auth")
	authHandler := handlers.NewAuthHandler(c)

	// Public routes
	auth.Post("/signup", authHandler.Signup)
	auth.Post("/login", authHandler.Login)
	auth.Post("/google", authHandler.GoogleLogin)
	auth.Post("/refresh", authHandler.RefreshToken)
	auth.Post("/logout", authHandler.Logout)

	// Protected routes (require authentication)
	authMiddleware := middleware.AuthMiddleware(c.JWTService)
	auth.Get("/me", authMiddleware, authHandler.GetMe)
	auth.Post("/claim-sessions", authMiddleware, authHandler.ClaimSessions)
}

func setupSearchHistoryRoutes(api fiber.Router, c *container.Container, authMiddleware fiber.Handler) {
	historyHandler := handlers.NewSearchHistoryHandler(c)

	// Search history routes (REQUIRE authentication)
	api.Get("/search-history", authMiddleware, historyHandler.GetSearchHistory)
	api.Delete("/search-history/:id", authMiddleware, historyHandler.DeleteSearchHistory)
	api.Post("/search-history/:id/click", authMiddleware, historyHandler.TrackProductClick)
	api.Delete("/search-history", authMiddleware, historyHandler.DeleteAllSearchHistory)
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
