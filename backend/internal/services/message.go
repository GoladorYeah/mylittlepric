package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"mylittleprice/internal/constants"
	"mylittleprice/internal/models"

	"github.com/redis/go-redis/v9"
)

// MessageService handles message-related operations
// Separated from SessionService for better SRP (Single Responsibility Principle)
type MessageService struct {
	redis *redis.Client
	ctx   context.Context
	ttl   time.Duration
}

// NewMessageService creates a new MessageService instance
func NewMessageService(redisClient *redis.Client, sessionTTL int) *MessageService {
	return &MessageService{
		redis: redisClient,
		ctx:   context.Background(),
		ttl:   time.Duration(sessionTTL) * time.Second,
	}
}

// IncrementMessageCountInMemory increments message count in an in-memory session (avoids N+1)
func (s *MessageService) IncrementMessageCountInMemory(session *models.ChatSession) {
	session.MessageCount++
}

// AddMessage adds a message to a session's message list (by sessionID)
func (s *MessageService) AddMessage(sessionID string, message *models.Message) error {
	key := fmt.Sprintf(constants.CachePrefixMessages, sessionID)

	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	pipe := s.redis.Pipeline()
	pipe.RPush(s.ctx, key, data)
	pipe.Expire(s.ctx, key, s.ttl)

	_, err = pipe.Exec(s.ctx)
	return err
}

// AddMessageInMemory adds a message using session object (still requires Redis for message storage)
// This method doesn't reload the session, so it avoids N+1 when used with in-memory session
func (s *MessageService) AddMessageInMemory(session *models.ChatSession, message *models.Message) error {
	key := fmt.Sprintf(constants.CachePrefixMessages, session.SessionID)

	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	pipe := s.redis.Pipeline()
	pipe.RPush(s.ctx, key, data)
	pipe.Expire(s.ctx, key, s.ttl)

	_, err = pipe.Exec(s.ctx)
	return err
}

// GetMessages retrieves all messages for a session
func (s *MessageService) GetMessages(sessionID string) ([]*models.Message, error) {
	key := fmt.Sprintf(constants.CachePrefixMessages, sessionID)

	data, err := s.redis.LRange(s.ctx, key, 0, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}

	messages := make([]*models.Message, 0, len(data))
	for _, msgData := range data {
		var msg models.Message
		err = json.Unmarshal([]byte(msgData), &msg)
		if err != nil {
			fmt.Printf("⚠️ Failed to unmarshal message in session %s: %v\n", sessionID, err)
			continue
		}
		messages = append(messages, &msg)
	}

	return messages, nil
}

// GetConversationHistory retrieves conversation history as role/content pairs
func (s *MessageService) GetConversationHistory(sessionID string) ([]map[string]string, error) {
	messages, err := s.GetMessages(sessionID)
	if err != nil {
		return nil, err
	}

	history := make([]map[string]string, 0, len(messages))
	for _, msg := range messages {
		history = append(history, map[string]string{
			"role":    msg.Role,
			"content": msg.Content,
		})
	}

	return history, nil
}

// GetRecentMessages retrieves the last N messages for a session
func (s *MessageService) GetRecentMessages(sessionID string, count int) ([]*models.Message, error) {
	key := fmt.Sprintf(constants.CachePrefixMessages, sessionID)

	start := -int64(count)
	data, err := s.redis.LRange(s.ctx, key, start, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get recent messages: %w", err)
	}

	messages := make([]*models.Message, 0, len(data))
	for _, msgData := range data {
		var msg models.Message
		err = json.Unmarshal([]byte(msgData), &msg)
		if err != nil {
			fmt.Printf("⚠️ Failed to unmarshal recent message in session %s: %v\n", sessionID, err)
			continue
		}
		messages = append(messages, &msg)
	}

	return messages, nil
}
