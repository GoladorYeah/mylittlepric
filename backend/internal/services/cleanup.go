package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"mylittleprice/ent"
	"mylittleprice/ent/chatsession"
	"mylittleprice/ent/message"
)

// CleanupService handles periodic cleanup of expired data
type CleanupService struct {
	client *ent.Client
	ctx    context.Context
}

// NewCleanupService creates a new CleanupService
func NewCleanupService(client *ent.Client) *CleanupService {
	return &CleanupService{
		client: client,
		ctx:    context.Background(),
	}
}

// CleanupExpiredSessions removes expired chat sessions from the database
// Returns the number of sessions deleted
func (s *CleanupService) CleanupExpiredSessions() (int, error) {
	// Delete sessions that have expired
	deleted, err := s.client.ChatSession.Delete().
		Where(chatsession.ExpiresAtLT(time.Now())).
		Exec(s.ctx)

	if err != nil {
		return 0, fmt.Errorf("failed to cleanup expired sessions: %w", err)
	}

	if deleted > 0 {
		log.Printf("🧹 Cleaned up %d expired sessions", deleted)
	}

	return deleted, nil
}

// CleanupOrphanedMessages removes messages that belong to deleted sessions
// Returns the number of messages deleted
func (s *CleanupService) CleanupOrphanedMessages() (int, error) {
	// Get all session IDs
	sessions, err := s.client.ChatSession.Query().
		Select(chatsession.FieldID).
		All(s.ctx)

	if err != nil {
		return 0, fmt.Errorf("failed to query sessions: %w", err)
	}

	// Create set of valid session IDs
	validSessionIDs := make(map[string]bool)
	for _, session := range sessions {
		validSessionIDs[session.ID.String()] = true
	}

	// Delete messages that don't have a valid session
	// Note: This query might be slow for large datasets
	// Consider using a more efficient approach for production (e.g., JOIN query)
	allMessages, err := s.client.Message.Query().
		Select(message.FieldID, message.FieldSessionID).
		All(s.ctx)

	if err != nil {
		return 0, fmt.Errorf("failed to query messages: %w", err)
	}

	orphanedCount := 0
	for _, msg := range allMessages {
		if !validSessionIDs[msg.SessionID.String()] {
			// This message is orphaned
			if err := s.client.Message.DeleteOneID(msg.ID).Exec(s.ctx); err != nil {
				log.Printf("⚠️ Failed to delete orphaned message %s: %v", msg.ID.String(), err)
			} else {
				orphanedCount++
			}
		}
	}

	if orphanedCount > 0 {
		log.Printf("🧹 Cleaned up %d orphaned messages", orphanedCount)
	}

	return orphanedCount, nil
}

// CleanupOldMessages removes messages older than a specified duration
// This helps manage database size for high-volume applications
func (s *CleanupService) CleanupOldMessages(olderThan time.Duration) (int, error) {
	cutoffTime := time.Now().Add(-olderThan)

	deleted, err := s.client.Message.Delete().
		Where(message.CreatedAtLT(cutoffTime)).
		Exec(s.ctx)

	if err != nil {
		return 0, fmt.Errorf("failed to cleanup old messages: %w", err)
	}

	if deleted > 0 {
		log.Printf("🧹 Cleaned up %d old messages (older than %v)", deleted, olderThan)
	}

	return deleted, nil
}

// RunFullCleanup runs all cleanup operations
// This should be called periodically (e.g., daily via cron job)
func (s *CleanupService) RunFullCleanup() error {
	log.Println("🧹 Starting full cleanup...")

	// 1. Cleanup expired sessions
	sessionsDeleted, err := s.CleanupExpiredSessions()
	if err != nil {
		log.Printf("⚠️ Error during session cleanup: %v", err)
	}

	// 2. Cleanup orphaned messages (messages without sessions)
	messagesDeleted, err := s.CleanupOrphanedMessages()
	if err != nil {
		log.Printf("⚠️ Error during orphaned message cleanup: %v", err)
	}

	// 3. Cleanup very old messages (older than 90 days)
	// This is optional and can be configured
	oldMessagesDeleted, err := s.CleanupOldMessages(90 * 24 * time.Hour)
	if err != nil {
		log.Printf("⚠️ Error during old message cleanup: %v", err)
	}

	log.Printf("🧹 Cleanup completed: %d sessions, %d orphaned messages, %d old messages",
		sessionsDeleted, messagesDeleted, oldMessagesDeleted)

	return nil
}

// StartPeriodicCleanup starts a background goroutine that runs cleanup periodically
// interval: how often to run cleanup (e.g., 24 * time.Hour for daily cleanup)
func (s *CleanupService) StartPeriodicCleanup(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		// Run cleanup immediately on start
		if err := s.RunFullCleanup(); err != nil {
			log.Printf("❌ Initial cleanup failed: %v", err)
		}

		// Then run periodically
		for range ticker.C {
			if err := s.RunFullCleanup(); err != nil {
				log.Printf("❌ Periodic cleanup failed: %v", err)
			}
		}
	}()

	log.Printf("🕐 Periodic cleanup started (interval: %v)", interval)
}
