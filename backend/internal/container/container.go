package container

import (
	"context"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/redis/go-redis/v9/maintnotifications"
	"google.golang.org/genai"

	"mylittleprice/internal/config"
	"mylittleprice/internal/services"
	"mylittleprice/internal/utils"
)

type Container struct {
	Config *config.Config
	DB     *sqlx.DB
	Redis  *redis.Client
	ctx    context.Context

	GeminiRotator *utils.KeyRotator
	SerpRotator   *utils.KeyRotator
	JWTService    *utils.JWTService

	EmbeddingService     *services.EmbeddingService
	GeminiService        *services.GeminiService
	SerpService          *services.SerpService
	CacheService         *services.CacheService
	SessionService       *services.SessionService
	GoogleOAuthService   *services.GoogleOAuthService
	AuthService          *services.AuthService
	SearchHistoryService *services.SearchHistoryService
	PreferencesService   *services.PreferencesService
}

func NewContainer(cfg *config.Config) (*Container, error) {
	c := &Container{
		Config: cfg,
		ctx:    context.Background(),
	}

	if err := c.initDatabase(); err != nil {
		return nil, fmt.Errorf("failed to initialize Database: %w", err)
	}

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

func (c *Container) initDatabase() error {
	db, err := sqlx.Connect("postgres", c.Config.DatabaseURL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	c.DB = db
	log.Println("✅ Connected to PostgreSQL")
	return nil
}

func (c *Container) initRedis() error {
	c.Redis = redis.NewClient(&redis.Options{
		Addr:     c.Config.RedisURL,
		Password: c.Config.RedisPassword,
		DB:       c.Config.RedisDB,
		// Disable maintenance notifications (only needed for Redis Enterprise/Cloud)
		MaintNotificationsConfig: &maintnotifications.Config{
			Mode: maintnotifications.ModeDisabled,
		},
	})

	if err := c.Redis.Ping(c.ctx).Err(); err != nil {
		return fmt.Errorf("Redis ping failed: %w", err)
	}

	log.Println("✅ Connected to Redis")
	return nil
}

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

func (c *Container) initServices() error {
	// Initialize JWT Service
	c.JWTService = utils.NewJWTService(
		c.Config.JWTAccessSecret,
		c.Config.JWTRefreshSecret,
		c.Config.JWTAccessTTL,
		c.Config.JWTRefreshTTL,
	)
	log.Println("🔐 JWT Service initialized")

	// Initialize Google OAuth Service
	c.GoogleOAuthService = services.NewGoogleOAuthService(c.Config)
	log.Println("🔑 Google OAuth Service initialized")

	// Initialize Auth Service
	c.AuthService = services.NewAuthService(c.DB, c.Redis, c.JWTService, c.GoogleOAuthService)
	log.Println("🔑 Auth Service initialized")

	c.SessionService = services.NewSessionService(
		c.Redis,
		c.DB,
		c.Config.SessionTTL,
		c.Config.MaxMessagesPerSession,
	)

	apiKey, _, _ := c.GeminiRotator.GetNextKey()
	geminiClient, _ := genai.NewClient(c.ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})

	c.EmbeddingService = services.NewEmbeddingService(geminiClient, c.Redis, c.Config)
	log.Println("🧠 Embedding Service initialized")

	c.CacheService = services.NewCacheService(c.Redis, c.Config, c.EmbeddingService)

	c.GeminiService = services.NewGeminiService(c.GeminiRotator, c.Config, c.EmbeddingService)
	log.Printf("🎯 Smart Grounding: '%s' mode", c.Config.GeminiGroundingMode)
	if c.Config.GeminiUseGrounding {
		log.Println("🔍 Grounding: ENABLED (selective usage)")
	} else {
		log.Println("💬 Grounding: DISABLED globally")
	}

	c.SerpService = services.NewSerpService(c.SerpRotator, c.Config)

	c.SearchHistoryService = services.NewSearchHistoryService(c.DB)
	log.Println("📜 Search History Service initialized")

	c.PreferencesService = services.NewPreferencesService(c.DB, c.AuthService)
	log.Println("⚙️ Preferences Service initialized")

	log.Println("✅ All services initialized")
	return nil
}

func (c *Container) Close() error {
	log.Println("🛑 Shutting down container...")

	if c.DB != nil {
		if err := c.DB.Close(); err != nil {
			log.Printf("⚠️ Failed to close database: %v", err)
		}
	}

	if err := c.Redis.Close(); err != nil {
		return fmt.Errorf("failed to close Redis: %w", err)
	}

	log.Println("✅ Container closed gracefully")
	return nil
}

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
		"embedding": map[string]interface{}{
			"status": "ok",
		},
	}

	return health
}

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
