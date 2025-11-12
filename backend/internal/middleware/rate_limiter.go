package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

// RateLimiterConfig holds the configuration for rate limiting
type RateLimiterConfig struct {
	Redis         *redis.Client
	Max           int           // Maximum number of requests
	Window        time.Duration // Time window
	KeyPrefix     string        // Redis key prefix
	SkipFailOpen  bool          // If true, allow requests when Redis is down
	Message       string        // Custom error message
	StatusCode    int           // HTTP status code for rate limit exceeded
	KeyGenerator  func(*fiber.Ctx) string // Custom key generator
}

// DefaultRateLimiterConfig returns default configuration
func DefaultRateLimiterConfig(redis *redis.Client) RateLimiterConfig {
	return RateLimiterConfig{
		Redis:        redis,
		Max:          100,
		Window:       1 * time.Minute,
		KeyPrefix:    "rate_limit:",
		SkipFailOpen: true, // Allow requests if Redis is down
		Message:      "Too many requests, please try again later",
		StatusCode:   fiber.StatusTooManyRequests,
		KeyGenerator: func(c *fiber.Ctx) string {
			// Use IP address as default key
			return c.IP()
		},
	}
}

// RateLimiter creates a new rate limiting middleware
func RateLimiter(config RateLimiterConfig) fiber.Handler {
	// Set defaults if not provided
	if config.Max <= 0 {
		config.Max = 100
	}
	if config.Window <= 0 {
		config.Window = 1 * time.Minute
	}
	if config.KeyPrefix == "" {
		config.KeyPrefix = "rate_limit:"
	}
	if config.Message == "" {
		config.Message = "Too many requests, please try again later"
	}
	if config.StatusCode == 0 {
		config.StatusCode = fiber.StatusTooManyRequests
	}
	if config.KeyGenerator == nil {
		config.KeyGenerator = func(c *fiber.Ctx) string {
			return c.IP()
		}
	}

	return func(c *fiber.Ctx) error {
		// Generate unique key for this client
		key := config.KeyPrefix + config.KeyGenerator(c)

		ctx := context.Background()

		// Increment counter and get current count
		pipe := config.Redis.Pipeline()
		incrCmd := pipe.Incr(ctx, key)
		pipe.Expire(ctx, key, config.Window)

		_, err := pipe.Exec(ctx)
		if err != nil {
			// Redis error - fail open if configured
			if config.SkipFailOpen {
				fmt.Printf("⚠️ Rate limiter Redis error (failing open): %v\n", err)
				return c.Next()
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "rate_limiter_error",
				"message": "Internal server error",
			})
		}

		// Get the count
		count := incrCmd.Val()

		// Check if limit exceeded
		if count > int64(config.Max) {
			// Get TTL for Retry-After header
			ttl, _ := config.Redis.TTL(ctx, key).Result()

			c.Set("Retry-After", fmt.Sprintf("%d", int(ttl.Seconds())))
			c.Set("X-RateLimit-Limit", fmt.Sprintf("%d", config.Max))
			c.Set("X-RateLimit-Remaining", "0")
			c.Set("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(ttl).Unix()))

			return c.Status(config.StatusCode).JSON(fiber.Map{
				"error": "rate_limit_exceeded",
				"message": config.Message,
				"retry_after": int(ttl.Seconds()),
			})
		}

		// Set rate limit headers
		c.Set("X-RateLimit-Limit", fmt.Sprintf("%d", config.Max))
		c.Set("X-RateLimit-Remaining", fmt.Sprintf("%d", config.Max-int(count)))
		c.Set("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(config.Window).Unix()))

		return c.Next()
	}
}

// WebSocketRateLimiter creates a rate limiter specifically for WebSocket connections
// This checks connection rate, not message rate
func WebSocketRateLimiter(redis *redis.Client, maxConnectionsPerMinute int) fiber.Handler {
	config := RateLimiterConfig{
		Redis:      redis,
		Max:        maxConnectionsPerMinute,
		Window:     1 * time.Minute,
		KeyPrefix:  "ws_conn_limit:",
		Message:    "Too many WebSocket connection attempts, please try again later",
		StatusCode: fiber.StatusTooManyRequests,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
	}

	return RateLimiter(config)
}

// AuthRateLimiter creates a rate limiter for authentication endpoints
func AuthRateLimiter(redis *redis.Client) fiber.Handler {
	config := RateLimiterConfig{
		Redis:      redis,
		Max:        10, // 10 attempts per 5 minutes
		Window:     5 * time.Minute,
		KeyPrefix:  "auth_limit:",
		Message:    "Too many authentication attempts, please try again later",
		StatusCode: fiber.StatusTooManyRequests,
		KeyGenerator: func(c *fiber.Ctx) string {
			// Rate limit by IP for auth endpoints
			return c.IP()
		},
	}

	return RateLimiter(config)
}
