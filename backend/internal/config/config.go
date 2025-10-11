package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	// Server
	Port string
	Env  string

	// Redis (PostgreSQL removed - not used)
	RedisURL      string
	RedisPassword string
	RedisDB       int

	// Session
	SessionTTL            int
	MaxMessagesPerSession int
	MaxSearchesPerSession int // Added from .env

	// API Keys
	GeminiAPIKeys []string
	SerpAPIKeys   []string

	// Gemini Configuration
	GeminiModel           string
	GeminiTemperature     float32
	GeminiMaxOutputTokens int
	GeminiUseGrounding    bool

	// Grounding Strategy Settings
	GeminiGroundingMode     string // "conservative", "balanced", "aggressive"
	GeminiGroundingMinWords int

	// Cache TTL (seconds)
	CacheGeminiTTL    int
	CacheSerpTTL      int
	CacheImmersiveTTL int

	// Rate Limiting
	RateLimitRequests int
	RateLimitWindow   int

	// CORS
	CORSOrigins []string

	// Logging
	LogLevel string
}

func Load() (*Config, error) {
	// Load .env file (ignore error if not exists)
	_ = godotenv.Load()

	config := &Config{
		Port:                  getEnv("PORT", "8080"),
		Env:                   getEnv("ENV", "development"),
		RedisURL:              getEnv("REDIS_URL", "localhost:6379"),
		RedisPassword:         getEnv("REDIS_PASSWORD", ""),
		RedisDB:               getEnvAsInt("REDIS_DB", 0),
		SessionTTL:            getEnvAsInt("SESSION_TTL", 86400),
		MaxMessagesPerSession: getEnvAsInt("MAX_MESSAGES_PER_SESSION", 8),
		MaxSearchesPerSession: getEnvAsInt("MAX_SEARCHES_PER_SESSION", 3),
		GeminiAPIKeys:         getEnvAsSlice("GEMINI_API_KEYS", []string{}),
		SerpAPIKeys:           getEnvAsSlice("SERP_API_KEYS", []string{}),
		GeminiModel:           getEnv("GEMINI_MODEL", "gemini-2.5-flash-preview-09-2025"),
		GeminiTemperature:     float32(getEnvAsFloat("GEMINI_TEMPERATURE", 0.7)),
		GeminiMaxOutputTokens: getEnvAsInt("GEMINI_MAX_OUTPUT_TOKENS", 1100),
		GeminiUseGrounding:    getEnvAsBool("GEMINI_USE_GROUNDING", true),

		// Grounding Strategy
		GeminiGroundingMode:     getEnv("GEMINI_GROUNDING_MODE", "balanced"),
		GeminiGroundingMinWords: getEnvAsInt("GEMINI_GROUNDING_MIN_WORDS", 2),

		CacheGeminiTTL:    getEnvAsInt("CACHE_GEMINI_TTL", 3600),
		CacheSerpTTL:      getEnvAsInt("CACHE_SERP_TTL", 86400),
		CacheImmersiveTTL: getEnvAsInt("CACHE_IMMERSIVE_TTL", 43200),
		RateLimitRequests: getEnvAsInt("RATE_LIMIT_REQUESTS", 100),
		RateLimitWindow:   getEnvAsInt("RATE_LIMIT_WINDOW", 60),
		CORSOrigins:       getEnvAsSlice("CORS_ORIGINS", []string{"http://localhost:3000"}),
		LogLevel:          getEnv("LOG_LEVEL", "info"),
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
