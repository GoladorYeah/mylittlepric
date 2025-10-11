package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"mylittleprice/internal/config"
	"mylittleprice/internal/models"
)

// CacheService handles all caching operations
type CacheService struct {
	redis  *redis.Client
	config *config.Config
	ctx    context.Context
}

// NewCacheService creates a new cache service
func NewCacheService(cfg *config.Config) *CacheService {
	// Initialize Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisURL,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	return &CacheService{
		redis:  redisClient,
		config: cfg,
		ctx:    context.Background(),
	}
}

// NewCacheServiceWithClient creates cache service with existing Redis client
func NewCacheServiceWithClient(redisClient *redis.Client, cfg *config.Config) *CacheService {
	return &CacheService{
		redis:  redisClient,
		config: cfg,
		ctx:    context.Background(),
	}
}

// GetSearchResults retrieves cached search results
func (c *CacheService) GetSearchResults(cacheKey string) ([]models.ProductCard, error) {
	data, err := c.redis.Get(c.ctx, cacheKey).Bytes()
	if err == redis.Nil {
		return nil, fmt.Errorf("cache miss")
	}
	if err != nil {
		return nil, fmt.Errorf("redis error: %w", err)
	}

	var cards []models.ProductCard
	if err := json.Unmarshal(data, &cards); err != nil {
		return nil, fmt.Errorf("unmarshal error: %w", err)
	}

	return cards, nil
}

// SetSearchResults stores search results in cache
func (c *CacheService) SetSearchResults(cacheKey string, cards []models.ProductCard, ttl time.Duration) error {
	data, err := json.Marshal(cards)
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}

	return c.redis.Set(c.ctx, cacheKey, data, ttl).Err()
}

// GetProductByToken retrieves cached product details by page_token
func (c *CacheService) GetProductByToken(pageToken string) (map[string]interface{}, error) {
	cacheKey := fmt.Sprintf("product:%s", pageToken)

	data, err := c.redis.Get(c.ctx, cacheKey).Bytes()
	if err == redis.Nil {
		return nil, fmt.Errorf("cache miss")
	}
	if err != nil {
		return nil, fmt.Errorf("redis error: %w", err)
	}

	var product map[string]interface{}
	if err := json.Unmarshal(data, &product); err != nil {
		return nil, fmt.Errorf("unmarshal error: %w", err)
	}

	return product, nil
}

// SetProductByToken stores product details in cache
func (c *CacheService) SetProductByToken(pageToken string, product map[string]interface{}, ttl int) error {
	cacheKey := fmt.Sprintf("product:%s", pageToken)

	data, err := json.Marshal(product)
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}

	duration := time.Duration(ttl) * time.Second
	return c.redis.Set(c.ctx, cacheKey, data, duration).Err()
}

// GetGeminiResponse retrieves cached Gemini response
func (c *CacheService) GetGeminiResponse(cacheKey string) (*models.GeminiResponse, error) {
	data, err := c.redis.Get(c.ctx, cacheKey).Bytes()
	if err == redis.Nil {
		return nil, fmt.Errorf("cache miss")
	}
	if err != nil {
		return nil, fmt.Errorf("redis error: %w", err)
	}

	var response models.GeminiResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("unmarshal error: %w", err)
	}

	return &response, nil
}

// SetGeminiResponse stores Gemini response in cache
func (c *CacheService) SetGeminiResponse(cacheKey string, response *models.GeminiResponse) error {
	data, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}

	ttl := time.Duration(c.config.CacheGeminiTTL) * time.Second
	return c.redis.Set(c.ctx, cacheKey, data, ttl).Err()
}

// InvalidateProduct removes product from cache
func (c *CacheService) InvalidateProduct(pageToken string) error {
	cacheKey := fmt.Sprintf("product:%s", pageToken)
	return c.redis.Del(c.ctx, cacheKey).Err()
}

// InvalidateSearch removes search results from cache
func (c *CacheService) InvalidateSearch(cacheKey string) error {
	return c.redis.Del(c.ctx, cacheKey).Err()
}

// GetCacheStats returns cache statistics
func (c *CacheService) GetCacheStats() (map[string]interface{}, error) {
	// Get total keys
	var cursor uint64
	var keys []string

	for {
		var scanKeys []string
		var err error

		scanKeys, cursor, err = c.redis.Scan(c.ctx, cursor, "*", 100).Result()
		if err != nil {
			return nil, err
		}

		keys = append(keys, scanKeys...)

		if cursor == 0 {
			break
		}
	}

	// Count by type
	productCount := 0
	searchCount := 0
	sessionCount := 0
	otherCount := 0

	for _, key := range keys {
		if len(key) >= 8 {
			prefix := key[:8]
			switch prefix {
			case "product:":
				productCount++
			case "search::":
				searchCount++
			case "session:":
				sessionCount++
			default:
				otherCount++
			}
		}
	}

	// Get memory info
	info, err := c.redis.Info(c.ctx, "memory").Result()
	memoryUsed := "unknown"
	if err == nil {
		// Parse memory info
		lines := splitString(info, "\n")
		for _, line := range lines {
			if len(line) > 12 && line[:12] == "used_memory:" {
				memoryUsed = line[12:]
				break
			}
		}
	}

	stats := map[string]interface{}{
		"total_keys":    len(keys),
		"product_cache": productCount,
		"search_cache":  searchCount,
		"session_cache": sessionCount,
		"other_keys":    otherCount,
		"memory_used":   memoryUsed,
	}

	return stats, nil
}

// ClearAll clears all cache (use with caution!)
func (c *CacheService) ClearAll() error {
	return c.redis.FlushDB(c.ctx).Err()
}

// ClearExpired removes expired keys (Redis does this automatically, but can be triggered manually)
func (c *CacheService) ClearExpired() error {
	// Redis handles expiration automatically
	// This is just a placeholder if manual cleanup is needed
	return nil
}

// Ping checks if Redis is accessible
func (c *CacheService) Ping() error {
	return c.redis.Ping(c.ctx).Err()
}

// Close closes Redis connection
func (c *CacheService) Close() error {
	return c.redis.Close()
}

// Helper functions

func splitString(s, sep string) []string {
	if s == "" {
		return []string{}
	}

	var result []string
	start := 0

	for i := 0; i <= len(s)-len(sep); i++ {
		if s[i:i+len(sep)] == sep {
			result = append(result, s[start:i])
			start = i + len(sep)
			i += len(sep) - 1
		}
	}

	// Add last part
	if start < len(s) {
		result = append(result, s[start:])
	}

	return result
}

// BuildSearchCacheKey creates a consistent cache key for search results
func BuildSearchCacheKey(country, searchType, query string) string {
	return fmt.Sprintf("search:%s:%s:%s", country, searchType, query)
}

// BuildProductCacheKey creates a consistent cache key for product details
func BuildProductCacheKey(pageToken string) string {
	return fmt.Sprintf("product:%s", pageToken)
}

// BuildGeminiCacheKey creates a cache key for Gemini responses
func BuildGeminiCacheKey(sessionID, message string) string {
	// Hash the message to create shorter key
	return fmt.Sprintf("gemini:%s:%s", sessionID, hashString(message))
}

// hashString creates a simple hash of a string (for cache keys)
func hashString(s string) string {
	// Simple hash function - in production use crypto/sha256
	hash := uint32(0)
	for i := 0; i < len(s); i++ {
		hash = hash*31 + uint32(s[i])
	}
	return fmt.Sprintf("%x", hash)
}
