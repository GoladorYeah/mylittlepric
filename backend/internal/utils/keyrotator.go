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
	keys       []string
	serviceName string
	redis      *redis.Client
	mu         sync.Mutex
	ctx        context.Context
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

// GetNextKey returns the next API key in rotation
func (kr *KeyRotator) GetNextKey() (string, int, error) {
	kr.mu.Lock()
	defer kr.mu.Unlock()

	if len(kr.keys) == 0 {
		return "", -1, fmt.Errorf("no API keys available for %s", kr.serviceName)
	}

	// Single key - no rotation needed
	if len(kr.keys) == 1 {
		return kr.keys[0], 0, nil
	}

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
	
	return kr.keys[index], index, nil
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