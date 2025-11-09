package utils

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// KeyRotator manages API key rotation using Redis
type KeyRotator struct {
	keys        []string
	serviceName string
	redis       *redis.Client
	mu          sync.Mutex
	ctx         context.Context
}

// NewKeyRotator creates a new key rotator instance
func NewKeyRotator(ctx context.Context, serviceName string, keys []string, redisClient *redis.Client) *KeyRotator {
	return &KeyRotator{
		keys:        keys,
		serviceName: serviceName,
		redis:       redisClient,
		ctx:         ctx,
	}
}

// GetNextKey returns the next API key in rotation, skipping exhausted keys
func (kr *KeyRotator) GetNextKey() (string, int, error) {
	kr.mu.Lock()
	defer kr.mu.Unlock()

	if len(kr.keys) == 0 {
		return "", -1, fmt.Errorf("no API keys available for %s", kr.serviceName)
	}

	// Try to find an available key
	maxAttempts := len(kr.keys)
	for attempt := 0; attempt < maxAttempts; attempt++ {
		// Get current index from Redis
		counterKey := fmt.Sprintf("keyrotator:%s:counter", kr.serviceName)

		// Increment and get the counter (atomic operation)
		counter, err := kr.redis.Incr(kr.ctx, counterKey).Result()
		if err != nil {
			// Fallback to first key if Redis fails
			return kr.keys[0], 0, fmt.Errorf("redis error, using first key: %w", err)
		}

		// Calculate index using modulo
		index := int(counter-1) % len(kr.keys)

		// Check if this key is exhausted
		if !kr.isKeyExhausted(index) {
			return kr.keys[index], index, nil
		}

		// Key is exhausted, try next one
		fmt.Printf("   ⏭️  Key %d is exhausted, trying next key...\n", index)
	}

	// All keys are exhausted
	return "", -1, fmt.Errorf("all API keys are exhausted for %s", kr.serviceName)
}

// isKeyExhausted checks if a key has been marked as exhausted (quota exceeded)
func (kr *KeyRotator) isKeyExhausted(keyIndex int) bool {
	exhaustedKey := fmt.Sprintf("keyrotator:%s:exhausted:%d", kr.serviceName, keyIndex)
	exists, err := kr.redis.Exists(kr.ctx, exhaustedKey).Result()
	if err != nil {
		// If Redis fails, assume key is available
		return false
	}
	return exists > 0
}

// MarkKeyAsExhausted marks a key as exhausted (quota exceeded) until end of day
func (kr *KeyRotator) MarkKeyAsExhausted(keyIndex int) error {
	exhaustedKey := fmt.Sprintf("keyrotator:%s:exhausted:%d", kr.serviceName, keyIndex)

	// Calculate TTL until end of day (UTC)
	now := time.Now().UTC()
	endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.UTC)
	ttl := endOfDay.Sub(now)

	// If less than 1 minute left, set for 24 hours
	if ttl < time.Minute {
		ttl = 24 * time.Hour
	}

	err := kr.redis.Set(kr.ctx, exhaustedKey, "1", ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to mark key as exhausted: %w", err)
	}

	fmt.Printf("   🚫 Key %d marked as exhausted (will reset in %v)\n", keyIndex, ttl.Round(time.Minute))
	return nil
}

// GetKeyByIndex returns a specific key by index
func (kr *KeyRotator) GetKeyByIndex(index int) (string, error) {
	if index < 0 || index >= len(kr.keys) {
		return "", fmt.Errorf("invalid key index: %d", index)
	}
	return kr.keys[index], nil
}

// RecordUsage records API key usage for analytics
func (kr *KeyRotator) RecordUsage(keyIndex int, success bool, responseTime time.Duration) error {
	usageKey := fmt.Sprintf("keyrotator:%s:usage:%d", kr.serviceName, keyIndex)

	// Increment usage counter
	pipe := kr.redis.Pipeline()
	pipe.Incr(kr.ctx, usageKey)

	// Record success/failure
	if success {
		pipe.Incr(kr.ctx, fmt.Sprintf("%s:success", usageKey))
	} else {
		pipe.Incr(kr.ctx, fmt.Sprintf("%s:failures", usageKey))
	}

	// Record response time (milliseconds)
	pipe.HIncrBy(kr.ctx, fmt.Sprintf("%s:response_times", usageKey), "total", responseTime.Milliseconds())
	pipe.HIncrBy(kr.ctx, fmt.Sprintf("%s:response_times", usageKey), "count", 1)

	_, err := pipe.Exec(kr.ctx)
	return err
}

// GetKeyStats returns usage statistics for a specific key
func (kr *KeyRotator) GetKeyStats(keyIndex int) (map[string]interface{}, error) {
	usageKey := fmt.Sprintf("keyrotator:%s:usage:%d", kr.serviceName, keyIndex)

	// Get all stats
	pipe := kr.redis.Pipeline()
	totalUsage := pipe.Get(kr.ctx, usageKey)
	successCount := pipe.Get(kr.ctx, fmt.Sprintf("%s:success", usageKey))
	failureCount := pipe.Get(kr.ctx, fmt.Sprintf("%s:failures", usageKey))
	responseTimes := pipe.HGetAll(kr.ctx, fmt.Sprintf("%s:response_times", usageKey))

	_, err := pipe.Exec(kr.ctx)
	if err != nil && err != redis.Nil {
		return nil, err
	}

	stats := map[string]interface{}{
		"key_index":     keyIndex,
		"total_usage":   totalUsage.Val(),
		"success_count": successCount.Val(),
		"failure_count": failureCount.Val(),
	}

	// Calculate average response time
	rtMap, _ := responseTimes.Result()
	if total, ok := rtMap["total"]; ok {
		if count, ok := rtMap["count"]; ok {
			if countInt, _ := time.ParseDuration(count + "ms"); countInt > 0 {
				totalInt, _ := time.ParseDuration(total + "ms")
				avg := totalInt / countInt
				stats["avg_response_time_ms"] = avg.Milliseconds()
			}
		}
	}

	return stats, nil
}

// GetAllStats returns statistics for all keys
func (kr *KeyRotator) GetAllStats() ([]map[string]interface{}, error) {
	stats := make([]map[string]interface{}, len(kr.keys))

	for i := range kr.keys {
		keyStats, err := kr.GetKeyStats(i)
		if err != nil {
			return nil, err
		}
		stats[i] = keyStats
	}

	return stats, nil
}

// ResetCounter resets the rotation counter (useful for testing)
func (kr *KeyRotator) ResetCounter() error {
	counterKey := fmt.Sprintf("keyrotator:%s:counter", kr.serviceName)
	return kr.redis.Del(kr.ctx, counterKey).Err()
}

// GetTotalKeys returns the number of available keys
func (kr *KeyRotator) GetTotalKeys() int {
	return len(kr.keys)
}
