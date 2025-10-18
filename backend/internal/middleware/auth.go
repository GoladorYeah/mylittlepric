package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"mylittleprice/internal/utils"
)

// AuthMiddleware validates JWT tokens and adds user context
func AuthMiddleware(jwtService *utils.JWTService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Extract token from Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "missing authorization header",
			})
		}

		// Check for Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid authorization header format",
			})
		}

		token := parts[1]

		// Validate token
		claims, err := jwtService.ValidateAccessToken(token)
		if err != nil {
			if err == utils.ErrExpiredToken {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "token expired",
				})
			}
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid token",
			})
		}

		// Add user info to context
		c.Locals("user_id", claims.UserID)
		c.Locals("user_email", claims.Email)

		return c.Next()
	}
}

// OptionalAuthMiddleware validates JWT if present, but doesn't require it
func OptionalAuthMiddleware(jwtService *utils.JWTService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Next()
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Next()
		}

		token := parts[1]
		claims, err := jwtService.ValidateAccessToken(token)
		if err == nil {
			c.Locals("user_id", claims.UserID)
			c.Locals("user_email", claims.Email)
		}

		return c.Next()
	}
}

// GetUserID retrieves user ID from context
func GetUserID(c *fiber.Ctx) (uuid.UUID, bool) {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	return userID, ok
}

// GetUserEmail retrieves user email from context
func GetUserEmail(c *fiber.Ctx) (string, bool) {
	email, ok := c.Locals("user_email").(string)
	return email, ok
}
