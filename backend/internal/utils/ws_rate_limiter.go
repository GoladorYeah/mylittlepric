package utils

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

// WSRateLimiter implements per-connection and per-user rate limiting for WebSocket messages
type WSRateLimiter struct {
	// Connection-level limits (per clientID)
	connLimits     map[string]*rateLimitBucket
	connMu         sync.RWMutex
	connMaxMsg     int           // Max messages per window
	connWindow     time.Duration // Time window
	connBurstAllow int           // Allow burst above limit

	// User-level limits (per userID)
	userLimits     map[uuid.UUID]*rateLimitBucket
	userMu         sync.RWMutex
	userMaxMsg     int           // Max messages per window
	userWindow     time.Duration // Time window
	userBurstAllow int           // Allow burst above limit

	// Cleanup ticker
	cleanupTicker *time.Ticker
	stopCleanup   chan struct{}
}

type rateLimitBucket struct {
	messages  []time.Time // Timestamps of recent messages
	blocked   bool        // Whether this connection/user is currently blocked
	blockedAt time.Time   // When they were blocked
	blockUntil time.Time  // When the block expires
	mu        sync.Mutex
}

// WSRateLimitConfig defines rate limiting configuration
type WSRateLimitConfig struct {
	// Per-connection limits (prevents single connection spam)
	ConnMaxMessages int           // Max messages per connection per window
	ConnWindow      time.Duration // Connection rate limit window
	ConnBurst       int           // Allow burst messages

	// Per-user limits (prevents multi-connection spam)
	UserMaxMessages int           // Max messages per user per window
	UserWindow      time.Duration // User rate limit window
	UserBurst       int           // Allow burst messages

	// Block duration when limit exceeded
	BlockDuration time.Duration
}

// DefaultWSRateLimitConfig returns reasonable default configuration
func DefaultWSRateLimitConfig() *WSRateLimitConfig {
	return &WSRateLimitConfig{
		// Connection limits: 20 messages per minute, allow 5 burst
		ConnMaxMessages: 20,
		ConnWindow:      1 * time.Minute,
		ConnBurst:       5,

		// User limits: 50 messages per minute across all devices, allow 10 burst
		UserMaxMessages: 50,
		UserWindow:      1 * time.Minute,
		UserBurst:       10,

		// Block for 30 seconds when limit exceeded
		BlockDuration: 30 * time.Second,
	}
}

// NewWSRateLimiter creates a new WebSocket rate limiter
func NewWSRateLimiter(config *WSRateLimitConfig) *WSRateLimiter {
	rl := &WSRateLimiter{
		connLimits:     make(map[string]*rateLimitBucket),
		userLimits:     make(map[uuid.UUID]*rateLimitBucket),
		connMaxMsg:     config.ConnMaxMessages,
		connWindow:     config.ConnWindow,
		connBurstAllow: config.ConnBurst,
		userMaxMsg:     config.UserMaxMessages,
		userWindow:     config.UserWindow,
		userBurstAllow: config.UserBurst,
		stopCleanup:    make(chan struct{}),
	}

	// Start cleanup goroutine to remove old entries
	rl.cleanupTicker = time.NewTicker(5 * time.Minute)
	go rl.cleanupLoop()

	return rl
}

// CheckConnection checks if a connection is allowed to send a message
// Returns (allowed bool, reason string, retryAfter time.Duration)
func (rl *WSRateLimiter) CheckConnection(clientID string) (bool, string, time.Duration) {
	rl.connMu.Lock()
	bucket, exists := rl.connLimits[clientID]
	if !exists {
		bucket = &rateLimitBucket{
			messages: make([]time.Time, 0, rl.connMaxMsg+rl.connBurstAllow),
		}
		rl.connLimits[clientID] = bucket
	}
	rl.connMu.Unlock()

	return rl.checkBucket(bucket, rl.connMaxMsg, rl.connWindow, rl.connBurstAllow, "connection")
}

// CheckUser checks if a user is allowed to send a message
// Returns (allowed bool, reason string, retryAfter time.Duration)
func (rl *WSRateLimiter) CheckUser(userID uuid.UUID) (bool, string, time.Duration) {
	rl.userMu.Lock()
	bucket, exists := rl.userLimits[userID]
	if !exists {
		bucket = &rateLimitBucket{
			messages: make([]time.Time, 0, rl.userMaxMsg+rl.userBurstAllow),
		}
		rl.userLimits[userID] = bucket
	}
	rl.userMu.Unlock()

	return rl.checkBucket(bucket, rl.userMaxMsg, rl.userWindow, rl.userBurstAllow, "user")
}

// checkBucket performs the actual rate limit check
func (rl *WSRateLimiter) checkBucket(bucket *rateLimitBucket, maxMsg int, window time.Duration, burstAllow int, limitType string) (bool, string, time.Duration) {
	bucket.mu.Lock()
	defer bucket.mu.Unlock()

	now := time.Now()

	// Check if currently blocked
	if bucket.blocked && now.Before(bucket.blockUntil) {
		remaining := bucket.blockUntil.Sub(now)
		return false, fmt.Sprintf("Rate limit exceeded (%s). Try again in %v", limitType, remaining.Round(time.Second)), remaining
	}

	// Unblock if block period has expired
	if bucket.blocked && now.After(bucket.blockUntil) {
		bucket.blocked = false
		bucket.messages = bucket.messages[:0] // Clear old messages
	}

	// Remove messages outside the current window
	cutoff := now.Add(-window)
	validMessages := make([]time.Time, 0, len(bucket.messages))
	for _, msgTime := range bucket.messages {
		if msgTime.After(cutoff) {
			validMessages = append(validMessages, msgTime)
		}
	}
	bucket.messages = validMessages

	// Check if limit exceeded
	currentCount := len(bucket.messages)
	limit := maxMsg + burstAllow

	if currentCount >= limit {
		// Exceeded limit - block this connection/user
		bucket.blocked = true
		bucket.blockedAt = now
		bucket.blockUntil = now.Add(30 * time.Second) // Block for 30 seconds

		return false, fmt.Sprintf("Rate limit exceeded (%s): %d messages in %v. Blocked for 30s", limitType, currentCount, window), 30 * time.Second
	}

	// Allow message - record timestamp
	bucket.messages = append(bucket.messages, now)

	return true, "", 0
}

// RecordMessage records a message for rate limiting (call after successful send)
// This is used when you want to check limits separately
func (rl *WSRateLimiter) RecordMessage(clientID string, userID *uuid.UUID) {
	// Record for connection
	rl.connMu.RLock()
	connBucket, exists := rl.connLimits[clientID]
	rl.connMu.RUnlock()

	if exists {
		connBucket.mu.Lock()
		// Message already recorded by CheckConnection
		connBucket.mu.Unlock()
	}

	// Record for user if authenticated
	if userID != nil {
		rl.userMu.RLock()
		userBucket, exists := rl.userLimits[*userID]
		rl.userMu.RUnlock()

		if exists {
			userBucket.mu.Lock()
			// Message already recorded by CheckUser
			userBucket.mu.Unlock()
		}
	}
}

// RemoveConnection removes rate limit data for a disconnected client
func (rl *WSRateLimiter) RemoveConnection(clientID string) {
	rl.connMu.Lock()
	delete(rl.connLimits, clientID)
	rl.connMu.Unlock()
}

// cleanupLoop periodically cleans up old entries
func (rl *WSRateLimiter) cleanupLoop() {
	for {
		select {
		case <-rl.cleanupTicker.C:
			rl.cleanup()
		case <-rl.stopCleanup:
			rl.cleanupTicker.Stop()
			return
		}
	}
}

// cleanup removes old/expired entries to prevent memory leaks
func (rl *WSRateLimiter) cleanup() {
	now := time.Now()

	// Cleanup connection limits
	rl.connMu.Lock()
	for clientID, bucket := range rl.connLimits {
		bucket.mu.Lock()
		// Remove if no messages in last 10 minutes and not blocked
		if len(bucket.messages) == 0 || (len(bucket.messages) > 0 && now.Sub(bucket.messages[len(bucket.messages)-1]) > 10*time.Minute) {
			if !bucket.blocked || now.After(bucket.blockUntil) {
				delete(rl.connLimits, clientID)
			}
		}
		bucket.mu.Unlock()
	}
	rl.connMu.Unlock()

	// Cleanup user limits
	rl.userMu.Lock()
	for userID, bucket := range rl.userLimits {
		bucket.mu.Lock()
		// Remove if no messages in last 10 minutes and not blocked
		if len(bucket.messages) == 0 || (len(bucket.messages) > 0 && now.Sub(bucket.messages[len(bucket.messages)-1]) > 10*time.Minute) {
			if !bucket.blocked || now.After(bucket.blockUntil) {
				delete(rl.userLimits, userID)
			}
		}
		bucket.mu.Unlock()
	}
	rl.userMu.Unlock()
}

// Stop stops the rate limiter and cleanup goroutine
func (rl *WSRateLimiter) Stop() {
	close(rl.stopCleanup)
}

// GetStats returns rate limiter statistics
func (rl *WSRateLimiter) GetStats() map[string]interface{} {
	rl.connMu.RLock()
	connCount := len(rl.connLimits)
	rl.connMu.RUnlock()

	rl.userMu.RLock()
	userCount := len(rl.userLimits)
	rl.userMu.RUnlock()

	return map[string]interface{}{
		"tracked_connections": connCount,
		"tracked_users":       userCount,
		"config": map[string]interface{}{
			"conn_max_messages": rl.connMaxMsg,
			"conn_window":       rl.connWindow.String(),
			"user_max_messages": rl.userMaxMsg,
			"user_window":       rl.userWindow.String(),
		},
	}
}
