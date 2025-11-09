package app

import (
	"fmt"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"

	"mylittleprice/internal/container"
	"mylittleprice/internal/handlers"
	"mylittleprice/internal/middleware"
)

// SetupRoutes configures all application routes
func SetupRoutes(app *fiber.App, c *container.Container) {
	// Health check
	app.Get("/health", func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{
			"status":    "ok",
			"timestamp": time.Now(),
			"services":  c.HealthCheck(),
		})
	})

	api := app.Group("/api")

	// Authentication routes (public)
	setupAuthRoutes(api, c)

	// WebSocket chat (optional authentication)
	setupWebSocketRoutes(app, c)

	// Chat endpoints (optional authentication)
	setupChatRoutes(api, c)

	// Product routes
	setupProductRoutes(api, c)

	// Search history routes (optional authentication)
	setupSearchHistoryRoutes(api, c)

	// Stats routes
	setupStatsRoutes(api, c)
}

func setupAuthRoutes(api fiber.Router, c *container.Container) {
	auth := api.Group("/auth")
	authHandler := handlers.NewAuthHandler(c)
	authMiddleware := middleware.AuthMiddleware(c.JWTService)

	// Public routes
	auth.Post("/signup", authHandler.Signup)
	auth.Post("/login", authHandler.Login)
	auth.Post("/google", authHandler.GoogleLogin)
	auth.Post("/refresh", authHandler.RefreshToken)
	auth.Post("/logout", authHandler.Logout)

	// Protected routes
	auth.Get("/me", authMiddleware, authHandler.GetMe)
	auth.Post("/claim-sessions", authMiddleware, authHandler.ClaimSessions)
}

func setupWebSocketRoutes(app *fiber.App, c *container.Container) {
	wsHandler := handlers.NewWSHandler(c)

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

			// If token is provided, validate it
			if token != "" {
				claims, err := c.JWTService.ValidateAccessToken(token)
				if err == nil {
					// Store user info in locals for WebSocket handler
					ctx.Locals("user_id", claims.UserID)
					ctx.Locals("user_email", claims.Email)
				}
				// If token validation fails, we just proceed without authentication
			}

			ctx.Locals("allowed", true)
			return ctx.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws", websocket.New(func(conn *websocket.Conn) {
		wsHandler.HandleWebSocket(conn)
	}))
}

func setupChatRoutes(api fiber.Router, c *container.Container) {
	chatHandler := handlers.NewChatHandler(c)
	optionalAuthMiddleware := middleware.OptionalAuthMiddleware(c.JWTService)

	api.Post("/chat", optionalAuthMiddleware, chatHandler.HandleChat)
	api.Get("/chat/messages", optionalAuthMiddleware, chatHandler.GetSessionMessages)
}

func setupProductRoutes(api fiber.Router, c *container.Container) {
	productHandler := handlers.NewProductHandler(c)
	api.Post("/product-details", productHandler.HandleProductDetails)
}

func setupSearchHistoryRoutes(api fiber.Router, c *container.Container) {
	historyHandler := handlers.NewSearchHistoryHandler(c)
	authMiddleware := middleware.AuthMiddleware(c.JWTService)
	optionalAuthMiddleware := middleware.OptionalAuthMiddleware(c.JWTService)

	// Get search history - supports both authenticated and anonymous users
	api.Get("/search-history", optionalAuthMiddleware, historyHandler.GetSearchHistory)

	// Delete operations - support both authenticated and anonymous users
	api.Delete("/search-history/:id", optionalAuthMiddleware, historyHandler.DeleteSearchHistory)
	api.Post("/search-history/:id/click", optionalAuthMiddleware, historyHandler.TrackProductClick)

	// Delete all - requires authentication
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
