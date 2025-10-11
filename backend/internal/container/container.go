package container

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"

	"mylittleprice/internal/config"
	"mylittleprice/internal/services"
	"mylittleprice/internal/utils"
)

// Container holds all application dependencies
type Container struct {
	// Configuration
	Config *config.Config

	// Infrastructure
	Redis *redis.Client
	ctx   context.Context

	// Key Rotators
	GeminiRotator *utils.KeyRotator
	SerpRotator   *utils.KeyRotator

	// Services (Singletons)
	GeminiService  *services.GeminiService
	SerpService    *services.SerpService
	CacheService   *services.CacheService
	SessionService *services.SessionService
	Optimizer      *services.QueryOptimizer
}

// NewContainer creates and initializes the dependency injection container
func NewContainer(cfg *config.Config) (*Container, error) {
	c := &Container{
		Config: cfg,
		ctx:    context.Background(),
	}

	// Initialize in correct order
	if err := c.initRedis(); err != nil {
		return nil, fmt.Errorf("failed to initialize Redis: %w", err)
	}

	if err := c.initKeyRotators(); err != nil {
		return nil, fmt.Errorf("failed to initialize key rotators: %w", err)
	}

	if err := c.initServices(); err != nil {
		return nil, fmt.Errorf("failed to initialize services: %w", err)
	}

	log.Println("✅ Dependency container initialized successfully")
	return c, nil
}

// initRedis initializes Redis connection
func (c *Container) initRedis() error {
	c.Redis = redis.NewClient(&redis.Options{
		Addr:     c.Config.RedisURL,
		Password: c.Config.RedisPassword,
		DB:       c.Config.RedisDB,
	})

	if err := c.Redis.Ping(c.ctx).Err(); err != nil {
		return fmt.Errorf("Redis ping failed: %w", err)
	}

	log.Println("✅ Connected to Redis")
	return nil
}

// initKeyRotators initializes API key rotators
func (c *Container) initKeyRotators() error {
	c.GeminiRotator = utils.NewKeyRotator(
		c.ctx,
		"gemini",
		c.Config.GeminiAPIKeys,
		c.Redis,
	)

	c.SerpRotator = utils.NewKeyRotator(
		c.ctx,
		"serp",
		c.Config.SerpAPIKeys,
		c.Redis,
	)

	log.Printf("✅ Gemini Key Rotator: %d keys", c.GeminiRotator.GetTotalKeys())
	log.Printf("✅ SERP Key Rotator: %d keys", c.SerpRotator.GetTotalKeys())

	return nil
}

// initServices initializes all application services
func (c *Container) initServices() error {
	// Session Service
	c.SessionService = services.NewSessionService(
		c.Redis,
		c.Config.SessionTTL,
		c.Config.MaxMessagesPerSession,
	)

	// Cache Service
	c.CacheService = services.NewCacheServiceWithClient(c.Redis, c.Config)

	// Gemini Service (with smart grounding)
	c.GeminiService = services.NewGeminiService(c.GeminiRotator, c.Config)
	log.Printf("🎯 Smart Grounding: '%s' mode", c.Config.GeminiGroundingMode)
	if c.Config.GeminiUseGrounding {
		log.Println("🔍 Grounding: ENABLED (selective usage)")
	} else {
		log.Println("💬 Grounding: DISABLED globally")
	}

	// SERP Service
	c.SerpService = services.NewSerpService(c.SerpRotator, c.Config)

	// Query Optimizer
	c.Optimizer = services.NewQueryOptimizer()

	log.Println("✅ All services initialized")
	return nil
}

// Close gracefully shuts down all resources
func (c *Container) Close() error {
	log.Println("🛑 Shutting down container...")

	if err := c.Redis.Close(); err != nil {
		return fmt.Errorf("failed to close Redis: %w", err)
	}

	log.Println("✅ Container closed gracefully")
	return nil
}

// HealthCheck verifies all dependencies are working
func (c *Container) HealthCheck() map[string]interface{} {
	health := map[string]interface{}{
		"redis": c.checkRedis(),
		"gemini_keys": map[string]interface{}{
			"count":  c.GeminiRotator.GetTotalKeys(),
			"status": "ok",
		},
		"serp_keys": map[string]interface{}{
			"count":  c.SerpRotator.GetTotalKeys(),
			"status": "ok",
		},
		"grounding": map[string]interface{}{
			"mode":    c.Config.GeminiGroundingMode,
			"enabled": c.Config.GeminiUseGrounding,
		},
	}

	return health
}

// checkRedis checks Redis connection health
func (c *Container) checkRedis() map[string]interface{} {
	if err := c.Redis.Ping(c.ctx).Err(); err != nil {
		return map[string]interface{}{
			"status":  "error",
			"message": err.Error(),
		}
	}
	return map[string]interface{}{
		"status": "ok",
	}
}
