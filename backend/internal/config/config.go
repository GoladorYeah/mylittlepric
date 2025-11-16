package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	// Server
	Port string
	Env  string

	// Database
	DatabaseURL string

	// Redis
	RedisURL      string
	RedisPassword string
	RedisDB       int

	// Redis Connection Pool Settings
	RedisPoolSize     int
	RedisMinIdle      int
	RedisMaxIdle      int
	RedisDialTimeout  time.Duration
	RedisReadTimeout  time.Duration
	RedisWriteTimeout time.Duration
	RedisPoolTimeout  time.Duration

	// Redis Retry Settings
	RedisMaxRetries      int
	RedisMinRetryBackoff time.Duration
	RedisMaxRetryBackoff time.Duration

	// Redis Buffer Settings
	RedisReadBufferSize  int
	RedisWriteBufferSize int

	// JWT Authentication
	JWTAccessSecret  string
	JWTRefreshSecret string
	JWTAccessTTL     time.Duration
	JWTRefreshTTL    time.Duration

	// Google OAuth
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string

	// Session
	SessionTTL            int
	MaxMessagesPerSession int
	MaxSearchesPerSession int // Added from .env

	// API Keys
	GeminiAPIKeys []string
	SerpAPIKeys   []string

	// Gemini Configuration
	GeminiModel           string
	GeminiFallbackModel   string // Fallback model for retries
	GeminiTemperature     float32
	GeminiMaxOutputTokens int
	GeminiUseGrounding    bool

	// Grounding Strategy Settings
	GeminiGroundingMode     string // "conservative", "balanced", "aggressive"
	GeminiGroundingMinWords int

	// Grounding Strategy Thresholds
	GeminiBrandQueryConfidence     float64
	GeminiGroundingWeightFreshInfo float64
	GeminiGroundingWeightSpecific  float64
	GeminiGroundingWeightDrift     float64
	GeminiGroundingWeightElectron  float64
	GeminiGroundingDecisionThresh  float64
	GeminiBrandQueryMaxWords       int
	GeminiBrandSimilarityThresh    float64
	GeminiDialogueHistoryWindow    int
	GeminiDialogueDriftThresh      float64
	GeminiDriftScoreBonus          float64
	GeminiElectronicsThreshHigh    float64
	GeminiElectronicsScoreHigh     float64
	GeminiCategorySimilarityThresh float64
	GeminiCategoryScore            float64

	// Gemini Translation Settings
	GeminiTranslationTemperature float32
	GeminiTranslationMaxTokens   int

	// Embedding Settings
	GeminiEmbeddingModel             string
	EmbeddingCategoryDetectionThresh float64
	CacheQueryEmbeddingTTL           int

	// SERP Relevance Thresholds
	SerpThresholdExact      float64
	SerpThresholdParameters float64
	SerpThresholdCategory   float64
	SerpLogTopResultsCount  int
	SerpFallbackMinResults  int

	// SERP Scoring Weights
	SerpScorePhraseMatch     float64
	SerpScoreAllWords        float64
	SerpScorePartialWords    float64
	SerpScoreWordOrderWeight float64
	SerpScoreBrandMatch      float64
	SerpScoreModelMatch      float64
	SerpMinWordLength        int
	SerpModelNumberMinLength int

	// SERP Max Products
	SerpMaxProductsExact      int
	SerpMaxProductsParameters int
	SerpMaxProductsCategory   int
	SerpMaxProductsDefault    int

	// Default Values
	DefaultCountry  string
	DefaultLanguage string
	DefaultCurrency string

	// Cache TTL (seconds)
	CacheGeminiTTL    int
	CacheSerpTTL      int
	CacheImmersiveTTL int

	// Rate Limiting
	RateLimitRequests int
	RateLimitWindow   int

	// CORS
	CORSOrigins []string

	// Notifications
	DiscordWebhookURL string

	// Logging
	LogLevel  string
	LogFormat string // "json" or "text"

	// Loki Logging
	LokiEnabled     bool
	LokiURL         string
	LokiServiceName string
}

func Load() (*Config, error) {
	// Load .env file (ignore error if not exists)
	_ = godotenv.Load()

	config := &Config{
		Port:                  getEnv("PORT", "8080"),
		Env:                   getEnv("ENV", "development"),
		DatabaseURL:           getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/mylittleprice?sslmode=disable"),
		RedisURL:              getEnv("REDIS_URL", "localhost:6379"),
		RedisPassword:         getEnv("REDIS_PASSWORD", ""),
		RedisDB:               getEnvAsInt("REDIS_DB", 0),
		RedisPoolSize:         getEnvAsInt("REDIS_POOL_SIZE", 50),
		RedisMinIdle:          getEnvAsInt("REDIS_MIN_IDLE", 10),
		RedisMaxIdle:          getEnvAsInt("REDIS_MAX_IDLE", 20),
		RedisDialTimeout:      time.Duration(getEnvAsInt("REDIS_DIAL_TIMEOUT", 5)) * time.Second,
		RedisReadTimeout:      time.Duration(getEnvAsInt("REDIS_READ_TIMEOUT", 3)) * time.Second,
		RedisWriteTimeout:     time.Duration(getEnvAsInt("REDIS_WRITE_TIMEOUT", 3)) * time.Second,
		RedisPoolTimeout:      time.Duration(getEnvAsInt("REDIS_POOL_TIMEOUT", 4)) * time.Second,
		RedisMaxRetries:       getEnvAsInt("REDIS_MAX_RETRIES", 3),
		RedisMinRetryBackoff:  time.Duration(getEnvAsInt("REDIS_MIN_RETRY_BACKOFF", 8)) * time.Millisecond,
		RedisMaxRetryBackoff:  time.Duration(getEnvAsInt("REDIS_MAX_RETRY_BACKOFF", 512)) * time.Millisecond,
		RedisReadBufferSize:   getEnvAsInt("REDIS_READ_BUFFER_SIZE", 1024*1024),   // 1 MiB
		RedisWriteBufferSize:  getEnvAsInt("REDIS_WRITE_BUFFER_SIZE", 1024*1024),  // 1 MiB
		JWTAccessSecret:       getEnv("JWT_ACCESS_SECRET", ""),
		JWTRefreshSecret:      getEnv("JWT_REFRESH_SECRET", ""),
		JWTAccessTTL:          time.Duration(getEnvAsInt("JWT_ACCESS_TTL", 900)) * time.Second,     // 15 minutes default
		JWTRefreshTTL:         time.Duration(getEnvAsInt("JWT_REFRESH_TTL", 604800)) * time.Second, // 7 days default
		GoogleClientID:        getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret:    getEnv("GOOGLE_CLIENT_SECRET", ""),
		GoogleRedirectURL:     getEnv("GOOGLE_REDIRECT_URL", ""),
		SessionTTL:            getEnvAsInt("SESSION_TTL", 86400),
		MaxMessagesPerSession: getEnvAsInt("MAX_MESSAGES_PER_SESSION", 8),
		MaxSearchesPerSession: getEnvAsInt("MAX_SEARCHES_PER_SESSION", 3),
		GeminiAPIKeys:         getEnvAsSlice("GEMINI_API_KEYS", []string{}),
		SerpAPIKeys:           getEnvAsSlice("SERP_API_KEYS", []string{}),
		GeminiModel:           getEnv("GEMINI_MODEL", "gemini-flash-latest"),
		GeminiFallbackModel:   getEnv("GEMINI_FALLBACK_MODEL", "gemini-flash-lite-latest"),
		GeminiTemperature:     float32(getEnvAsFloat("GEMINI_TEMPERATURE", 0.7)),
		GeminiMaxOutputTokens: getEnvAsInt("GEMINI_MAX_OUTPUT_TOKENS", 8192),
		GeminiUseGrounding:    getEnvAsBool("GEMINI_USE_GROUNDING", true),

		// Grounding Strategy
		GeminiGroundingMode:     getEnv("GEMINI_GROUNDING_MODE", "balanced"),
		GeminiGroundingMinWords: getEnvAsInt("GEMINI_GROUNDING_MIN_WORDS", 2),

		// Grounding Strategy Thresholds
		GeminiBrandQueryConfidence:     getEnvAsFloat("GEMINI_BRAND_QUERY_CONFIDENCE", 0.95),
		GeminiGroundingWeightFreshInfo: getEnvAsFloat("GEMINI_GROUNDING_WEIGHT_FRESH_INFO", 0.3),
		GeminiGroundingWeightSpecific:  getEnvAsFloat("GEMINI_GROUNDING_WEIGHT_SPECIFIC_PRODUCT", 0.35),
		GeminiGroundingWeightDrift:     getEnvAsFloat("GEMINI_GROUNDING_WEIGHT_DIALOGUE_DRIFT", 0.2),
		GeminiGroundingWeightElectron:  getEnvAsFloat("GEMINI_GROUNDING_WEIGHT_ELECTRONICS", 0.15),
		GeminiGroundingDecisionThresh:  getEnvAsFloat("GEMINI_GROUNDING_DECISION_THRESHOLD", 0.5),
		GeminiBrandQueryMaxWords:       getEnvAsInt("GEMINI_BRAND_QUERY_MAX_WORDS", 3),
		GeminiBrandSimilarityThresh:    getEnvAsFloat("GEMINI_BRAND_SIMILARITY_THRESHOLD", 0.65),
		GeminiDialogueHistoryWindow:    getEnvAsInt("GEMINI_DIALOGUE_HISTORY_WINDOW", 4),
		GeminiDialogueDriftThresh:      getEnvAsFloat("GEMINI_DIALOGUE_DRIFT_THRESHOLD", 0.4),
		GeminiDriftScoreBonus:          getEnvAsFloat("GEMINI_DRIFT_SCORE_BONUS", 0.8),
		GeminiElectronicsThreshHigh:    getEnvAsFloat("GEMINI_ELECTRONICS_THRESHOLD_HIGH", 0.7),
		GeminiElectronicsScoreHigh:     getEnvAsFloat("GEMINI_ELECTRONICS_SCORE_HIGH", 0.9),
		GeminiCategorySimilarityThresh: getEnvAsFloat("GEMINI_CATEGORY_SIMILARITY_THRESHOLD", 0.6),
		GeminiCategoryScore:            getEnvAsFloat("GEMINI_CATEGORY_SCORE", 0.5),

		// Gemini Translation Settings
		GeminiTranslationTemperature: float32(getEnvAsFloat("GEMINI_TRANSLATION_TEMPERATURE", 0.3)),
		GeminiTranslationMaxTokens:   getEnvAsInt("GEMINI_TRANSLATION_MAX_TOKENS", 100),

		// Embedding Settings
		GeminiEmbeddingModel:             getEnv("GEMINI_EMBEDDING_MODEL", "text-embedding-004"),
		EmbeddingCategoryDetectionThresh: getEnvAsFloat("EMBEDDING_CATEGORY_DETECTION_THRESHOLD", 0.6),
		CacheQueryEmbeddingTTL:           getEnvAsInt("CACHE_QUERY_EMBEDDING_TTL", 86400),

		// SERP Relevance Thresholds
		SerpThresholdExact:      getEnvAsFloat("SERP_THRESHOLD_EXACT", 0.4),
		SerpThresholdParameters: getEnvAsFloat("SERP_THRESHOLD_PARAMETERS", 0.2),
		SerpThresholdCategory:   getEnvAsFloat("SERP_THRESHOLD_CATEGORY", 0.1),
		SerpLogTopResultsCount:  getEnvAsInt("SERP_LOG_TOP_RESULTS_COUNT", 5),
		SerpFallbackMinResults:  getEnvAsInt("SERP_FALLBACK_MIN_RESULTS", 3),

		// SERP Scoring Weights
		SerpScorePhraseMatch:     getEnvAsFloat("SERP_SCORE_PHRASE_MATCH", 1.0),
		SerpScoreAllWords:        getEnvAsFloat("SERP_SCORE_ALL_WORDS", 0.6),
		SerpScorePartialWords:    getEnvAsFloat("SERP_SCORE_PARTIAL_WORDS", 0.5),
		SerpScoreWordOrderWeight: getEnvAsFloat("SERP_SCORE_WORD_ORDER_WEIGHT", 0.2),
		SerpScoreBrandMatch:      getEnvAsFloat("SERP_SCORE_BRAND_MATCH", 0.3),
		SerpScoreModelMatch:      getEnvAsFloat("SERP_SCORE_MODEL_MATCH", 0.3),
		SerpMinWordLength:        getEnvAsInt("SERP_MIN_WORD_LENGTH", 2),
		SerpModelNumberMinLength: getEnvAsInt("SERP_MODEL_NUMBER_MIN_LENGTH", 2),

		// SERP Max Products
		SerpMaxProductsExact:      getEnvAsInt("SERP_MAX_PRODUCTS_EXACT", 3),
		SerpMaxProductsParameters: getEnvAsInt("SERP_MAX_PRODUCTS_PARAMETERS", 6),
		SerpMaxProductsCategory:   getEnvAsInt("SERP_MAX_PRODUCTS_CATEGORY", 8),
		SerpMaxProductsDefault:    getEnvAsInt("SERP_MAX_PRODUCTS_DEFAULT", 6),

		// Default Values
		DefaultCountry:  getEnv("DEFAULT_COUNTRY", "CH"),
		DefaultLanguage: getEnv("DEFAULT_LANGUAGE", "en"),
		DefaultCurrency: getEnv("DEFAULT_CURRENCY", "CHF"),

		CacheGeminiTTL:    getEnvAsInt("CACHE_GEMINI_TTL", 3600),
		CacheSerpTTL:      getEnvAsInt("CACHE_SERP_TTL", 86400),
		CacheImmersiveTTL: getEnvAsInt("CACHE_IMMERSIVE_TTL", 43200),
		RateLimitRequests: getEnvAsInt("RATE_LIMIT_REQUESTS", 100),
		RateLimitWindow:   getEnvAsInt("RATE_LIMIT_WINDOW", 60),
		CORSOrigins:       getEnvAsSlice("CORS_ORIGINS", []string{"http://localhost:3000"}),

		// Notifications
		DiscordWebhookURL: os.Getenv("DISCORD_WEBHOOK_URL"),
		LogLevel:          getEnv("LOG_LEVEL", "info"),
		LogFormat:         getEnv("LOG_FORMAT", "json"),
		LokiEnabled:       getEnvAsBool("LOKI_ENABLED", false),
		LokiURL:           getEnv("LOKI_URL", "http://localhost:3100/loki/api/v1/push"),
		LokiServiceName:   getEnv("LOKI_SERVICE_NAME", "mylittleprice-backend"),
	}

	if err := config.validate(); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) validate() error {
	if len(c.GeminiAPIKeys) == 0 {
		return fmt.Errorf("at least one GEMINI_API_KEY is required")
	}

	if len(c.SerpAPIKeys) == 0 {
		return fmt.Errorf("at least one SERP_API_KEY is required")
	}

	// Validate Google OAuth config (required for authentication)
	if c.GoogleClientID == "" {
		return fmt.Errorf("GOOGLE_CLIENT_ID is required")
	}
	if c.GoogleClientSecret == "" {
		return fmt.Errorf("GOOGLE_CLIENT_SECRET is required")
	}

	// Validate grounding mode
	validModes := []string{"conservative", "balanced", "aggressive"}
	validMode := false
	for _, mode := range validModes {
		if c.GeminiGroundingMode == mode {
			validMode = true
			break
		}
	}
	if !validMode {
		return fmt.Errorf("GEMINI_GROUNDING_MODE must be one of: %v", validModes)
	}

	// Validate min words
	if c.GeminiGroundingMinWords < 1 || c.GeminiGroundingMinWords > 10 {
		return fmt.Errorf("GEMINI_GROUNDING_MIN_WORDS must be between 1 and 10")
	}

	// Validate max searches
	if c.MaxSearchesPerSession < 1 || c.MaxSearchesPerSession > 10 {
		return fmt.Errorf("MAX_SEARCHES_PER_SESSION must be between 1 and 10")
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsFloat(key string, defaultValue float64) float64 {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseFloat(valueStr, 64); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsSlice(key string, defaultValue []string) []string {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}

	values := strings.Split(valueStr, ",")
	result := make([]string, 0, len(values))
	for _, v := range values {
		trimmed := strings.TrimSpace(v)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	if len(result) == 0 {
		return defaultValue
	}

	return result
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}

	valueStr = strings.ToLower(valueStr)
	if valueStr == "true" || valueStr == "1" || valueStr == "yes" {
		return true
	}
	if valueStr == "false" || valueStr == "0" || valueStr == "no" {
		return false
	}

	return defaultValue
}
