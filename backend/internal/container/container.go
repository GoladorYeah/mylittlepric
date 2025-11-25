package container

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/redis/go-redis/v9/maintnotifications"
	"google.golang.org/genai"

	"mylittleprice/ent"
	"mylittleprice/internal/config"
	"mylittleprice/internal/metrics"
	"mylittleprice/internal/middleware"
	"mylittleprice/internal/services"
	"mylittleprice/internal/utils"
)

type Container struct {
	Config    *config.Config
	StartTime time.Time // Application start time for uptime tracking
	EntDB     *sql.DB     // SQL DB for Ent
	Ent       *ent.Client // Ent ORM client
	Redis     *redis.Client
	ctx       context.Context

	GeminiRotator *utils.KeyRotator
	SerpRotator   *utils.KeyRotator
	JWTService    *utils.JWTService

	EmbeddingService        *services.EmbeddingService
	GeminiService           *services.GeminiService
	SerpService             *services.SerpService
	CacheService            *services.CacheService
	SessionService          *services.SessionService
	MessageService          *services.MessageService
	CycleService            *services.CycleService
	GoogleOAuthService      *services.GoogleOAuthService
	AuthService             *services.AuthService
	EmailService            *services.EmailService
	SearchHistoryService    *services.SearchHistoryService
	PreferencesService      *services.PreferencesService
	CleanupService          *services.CleanupService
	SessionOwnershipChecker *middleware.SessionOwnershipValidator
}

func NewContainer(cfg *config.Config) (*Container, error) {
	c := &Container{
		Config:    cfg,
		StartTime: time.Now(),
		ctx:       context.Background(),
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

	utils.LogInfo(c.ctx, "dependency container initialized successfully")
	return c, nil
}

func (c *Container) initDatabase() error {
	// Initialize Ent client
	sqlDB, err := sql.Open("postgres", c.Config.DatabaseURL)
	if err != nil {
		return fmt.Errorf("failed to open database for Ent: %w", err)
	}

	// Set connection pool settings for Ent DB
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	c.EntDB = sqlDB

	// Create Ent client with the SQL driver
	drv := entsql.OpenDB(dialect.Postgres, sqlDB)
	entClient := ent.NewClient(ent.Driver(drv))

	// Run auto migrations
	if err := entClient.Schema.Create(c.ctx); err != nil {
		utils.LogWarn(c.ctx, "failed to create Ent schema (tables may already exist)", slog.Any("error", err))
	}

	c.Ent = entClient
	utils.LogInfo(c.ctx, "connected to PostgreSQL (Ent ORM)")
	return nil
}

func (c *Container) initRedis() error {
	c.Redis = redis.NewClient(&redis.Options{
		Addr:     c.Config.RedisURL,
		Password: c.Config.RedisPassword,
		DB:       c.Config.RedisDB,

		// Connection pool settings
		PoolSize:     c.Config.RedisPoolSize,
		MinIdleConns: c.Config.RedisMinIdle,
		MaxIdleConns: c.Config.RedisMaxIdle,

		// Timeouts
		DialTimeout:  c.Config.RedisDialTimeout,
		ReadTimeout:  c.Config.RedisReadTimeout,
		WriteTimeout: c.Config.RedisWriteTimeout,
		PoolTimeout:  c.Config.RedisPoolTimeout,

		// Buffer sizes (for high-throughput)
		ReadBufferSize:  c.Config.RedisReadBufferSize,
		WriteBufferSize: c.Config.RedisWriteBufferSize,

		// Retry configuration
		MaxRetries:      c.Config.RedisMaxRetries,
		MinRetryBackoff: c.Config.RedisMinRetryBackoff,
		MaxRetryBackoff: c.Config.RedisMaxRetryBackoff,

		// Disable maintenance notifications (only needed for Redis Enterprise/Cloud)
		MaintNotificationsConfig: &maintnotifications.Config{
			Mode: maintnotifications.ModeDisabled,
		},
	})

	// Health check with context timeout
	ctx, cancel := context.WithTimeout(c.ctx, 5*time.Second)
	defer cancel()

	if err := c.Redis.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("Redis ping failed: %w", err)
	}

	utils.LogInfo(c.ctx, "connected to Redis with optimized configuration",
		slog.Int("pool_size", c.Config.RedisPoolSize),
		slog.Int("min_idle", c.Config.RedisMinIdle),
		slog.Int("max_idle", c.Config.RedisMaxIdle),
	)
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

	utils.LogInfo(c.ctx, "Gemini key rotator initialized", slog.Int("total_keys", c.GeminiRotator.GetTotalKeys()))
	utils.LogInfo(c.ctx, "SERP key rotator initialized", slog.Int("total_keys", c.SerpRotator.GetTotalKeys()))

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
	utils.LogInfo(c.ctx, "JWT service initialized")

	// Initialize Google OAuth Service
	c.GoogleOAuthService = services.NewGoogleOAuthService(c.Config)
	utils.LogInfo(c.ctx, "Google OAuth service initialized")

	// Initialize Email Service
	c.EmailService = services.NewEmailService(c.Config)
	utils.LogInfo(c.ctx, "Email service initialized")

	// Initialize Auth Service
	c.AuthService = services.NewAuthService(c.Ent, c.Redis, c.JWTService, c.GoogleOAuthService)
	utils.LogInfo(c.ctx, "Auth service initialized")

	// Initialize CycleService (no dependencies)
	c.CycleService = services.NewCycleService()
	utils.LogInfo(c.ctx, "Cycle service initialized")

	// Initialize MessageService (depends on Redis and Ent)
	c.MessageService = services.NewMessageService(c.Redis, c.Ent, c.Config.SessionTTL)
	utils.LogInfo(c.ctx, "Message service initialized with PostgreSQL persistence")

	// Initialize SessionService (depends on CycleService)
	c.SessionService = services.NewSessionService(
		c.Redis,
		c.Ent,
		c.CycleService,
		c.Config.SessionTTL,
		c.Config.MaxMessagesPerSession,
	)
	c.SessionService.SetAuthService(c.AuthService)
	utils.LogInfo(c.ctx, "Session service initialized")

	apiKey, _, _ := c.GeminiRotator.GetNextKey()
	geminiClient, _ := genai.NewClient(c.ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})

	c.EmbeddingService = services.NewEmbeddingService(geminiClient, c.Redis, c.Config)
	utils.LogInfo(c.ctx, "Embedding service initialized")

	c.CacheService = services.NewCacheService(c.Redis, c.Config, c.EmbeddingService)

	c.GeminiService = services.NewGeminiService(c.GeminiRotator, c.Config, c.EmbeddingService)
	utils.LogInfo(c.ctx, "Smart grounding configured",
		slog.String("mode", c.Config.GeminiGroundingMode),
		slog.Bool("enabled", c.Config.GeminiUseGrounding),
	)

	c.SerpService = services.NewSerpService(c.SerpRotator, c.Config)

	c.SearchHistoryService = services.NewSearchHistoryService(c.Ent)
	utils.LogInfo(c.ctx, "Search history service initialized")

	c.PreferencesService = services.NewPreferencesService(c.Ent, c.AuthService)
	utils.LogInfo(c.ctx, "Preferences service initialized")

	c.CleanupService = services.NewCleanupService(c.Ent)
	utils.LogInfo(c.ctx, "Cleanup service initialized")

	// Start periodic cleanup (runs daily at 3 AM)
	c.CleanupService.StartPeriodicCleanup(24 * time.Hour)

	// Initialize Session Ownership Validator
	c.SessionOwnershipChecker = middleware.NewSessionOwnershipValidator(&services.SessionAdapter{SessionService: c.SessionService}, c.Config.JWTAccessSecret)
	utils.LogInfo(c.ctx, "Session ownership validation initialized")

	utils.LogInfo(c.ctx, "all services initialized")
	return nil
}

func (c *Container) Close() error {
	utils.LogInfo(c.ctx, "shutting down container")

	// Close Ent client
	if c.Ent != nil {
		if err := c.Ent.Close(); err != nil {
			utils.LogWarn(c.ctx, "failed to close Ent client", slog.Any("error", err))
		}
	}

	// Close EntDB connection
	if c.EntDB != nil {
		if err := c.EntDB.Close(); err != nil {
			utils.LogWarn(c.ctx, "failed to close Ent database", slog.Any("error", err))
		}
	}

	if err := c.Redis.Close(); err != nil {
		return fmt.Errorf("failed to close Redis: %w", err)
	}

	utils.LogInfo(c.ctx, "container closed gracefully")
	return nil
}

// RegisterMetrics registers all WebSocket and Session metrics
func (c *Container) RegisterMetrics() {
	metrics.RegisterWebSocketMetrics()
	metrics.RegisterSessionMetrics()
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
