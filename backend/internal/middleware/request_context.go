package middleware

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"mylittleprice/internal/utils"
)

// RequestContext adds request_id, user_id, and session_id to the request context
// This enables context-aware structured logging throughout the application
func RequestContext() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Generate a unique request ID
		requestID := uuid.New().String()

		// Get the base context from Fiber
		ctx := c.UserContext()

		// Add request ID to context
		ctx = utils.WithRequestID(ctx, requestID)

		// Add user ID to context if authenticated
		if userID := c.Locals("userID"); userID != nil {
			if userIDStr, ok := userID.(string); ok {
				ctx = utils.WithUserID(ctx, userIDStr)
			}
		}

		// Add session ID to context if available
		if sessionID := c.Locals("sessionID"); sessionID != nil {
			if sessionIDStr, ok := sessionID.(string); ok {
				ctx = utils.WithSessionID(ctx, sessionIDStr)
			}
		}

		// Set the enriched context back to Fiber
		c.SetUserContext(ctx)

		// Add request ID to response headers for tracing
		c.Set("X-Request-ID", requestID)

		return c.Next()
	}
}

// GetContext returns the context from Fiber context, or a background context if not available
func GetContext(c *fiber.Ctx) context.Context {
	if c == nil {
		return context.Background()
	}
	ctx := c.UserContext()
	if ctx == nil {
		return context.Background()
	}
	return ctx
}
