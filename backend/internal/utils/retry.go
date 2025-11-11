package utils

import (
	"context"
	"fmt"
	"log"
	"math"
	"time"
)

// RetryConfig holds configuration for retry logic
type RetryConfig struct {
	MaxRetries      int
	InitialDelay    time.Duration
	MaxDelay        time.Duration
	BackoffFactor   float64
	RetriableErrors []error // Specific errors that should trigger retry (optional)
}

// DefaultRetryConfig returns a sensible default retry configuration
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxRetries:    3,
		InitialDelay:  100 * time.Millisecond,
		MaxDelay:      5 * time.Second,
		BackoffFactor: 2.0,
	}
}

// RetryWithBackoff executes a function with exponential backoff retry logic
func RetryWithBackoff(ctx context.Context, fn func() error, config RetryConfig) error {
	var lastErr error

	for attempt := 0; attempt <= config.MaxRetries; attempt++ {
		// Check context cancellation
		select {
		case <-ctx.Done():
			return fmt.Errorf("retry cancelled: %w", ctx.Err())
		default:
		}

		// Execute the function
		err := fn()
		if err == nil {
			if attempt > 0 {
				log.Printf("✅ Operation succeeded on retry attempt %d", attempt+1)
			}
			return nil
		}

		lastErr = err

		// If this is the last attempt, don't wait
		if attempt == config.MaxRetries {
			break
		}

		// Calculate backoff delay using exponential backoff
		delay := time.Duration(float64(config.InitialDelay) * math.Pow(config.BackoffFactor, float64(attempt)))
		if delay > config.MaxDelay {
			delay = config.MaxDelay
		}

		log.Printf("⚠️ Operation failed (attempt %d/%d): %v. Retrying in %v...",
			attempt+1, config.MaxRetries+1, err, delay)

		// Wait with context cancellation support
		select {
		case <-ctx.Done():
			return fmt.Errorf("retry cancelled during backoff: %w", ctx.Err())
		case <-time.After(delay):
		}
	}

	return fmt.Errorf("operation failed after %d retries: %w", config.MaxRetries+1, lastErr)
}

// IsRetriableError determines if an error should trigger a retry
// This can be extended to check for specific error types (network errors, timeouts, etc.)
func IsRetriableError(err error) bool {
	if err == nil {
		return false
	}

	// Add specific error type checks here
	// For now, we'll retry on any error (let the config decide)
	return true
}

// RetryWithBackoffSelective is like RetryWithBackoff but only retries on specific errors
func RetryWithBackoffSelective(ctx context.Context, fn func() error, config RetryConfig, isRetriable func(error) bool) error {
	var lastErr error

	for attempt := 0; attempt <= config.MaxRetries; attempt++ {
		select {
		case <-ctx.Done():
			return fmt.Errorf("retry cancelled: %w", ctx.Err())
		default:
		}

		err := fn()
		if err == nil {
			if attempt > 0 {
				log.Printf("✅ Operation succeeded on retry attempt %d", attempt+1)
			}
			return nil
		}

		lastErr = err

		// Check if error is retriable
		if !isRetriable(err) {
			log.Printf("❌ Non-retriable error encountered: %v", err)
			return err
		}

		if attempt == config.MaxRetries {
			break
		}

		delay := time.Duration(float64(config.InitialDelay) * math.Pow(config.BackoffFactor, float64(attempt)))
		if delay > config.MaxDelay {
			delay = config.MaxDelay
		}

		log.Printf("⚠️ Retriable error (attempt %d/%d): %v. Retrying in %v...",
			attempt+1, config.MaxRetries+1, err, delay)

		select {
		case <-ctx.Done():
			return fmt.Errorf("retry cancelled during backoff: %w", ctx.Err())
		case <-time.After(delay):
		}
	}

	return fmt.Errorf("operation failed after %d retries: %w", config.MaxRetries+1, lastErr)
}
