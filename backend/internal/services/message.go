package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"mylittleprice/ent"
	"mylittleprice/ent/chatsession"
	"mylittleprice/ent/message"
	"mylittleprice/internal/constants"
	"mylittleprice/internal/models"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// MessageService handles message-related operations
// Separated from SessionService for better SRP (Single Responsibility Principle)
type MessageService struct {
	redis  *redis.Client
	client *ent.Client
	ctx    context.Context
	ttl    time.Duration
}

// NewMessageService creates a new MessageService instance
func NewMessageService(redisClient *redis.Client, entClient *ent.Client, sessionTTL int) *MessageService {
	return &MessageService{
		redis:  redisClient,
		client: entClient,
		ctx:    context.Background(),
		ttl:    time.Duration(sessionTTL) * time.Second,
	}
}

// IncrementMessageCountInMemory increments message count in an in-memory session (avoids N+1)
func (s *MessageService) IncrementMessageCountInMemory(session *models.ChatSession) {
	session.MessageCount++
}

// AddMessage adds a message to a session's message list (by sessionID)
// Saves to both PostgreSQL (persistent) and Redis (cache)
func (s *MessageService) AddMessage(sessionID string, msg *models.Message) error {
	// 1. Save to PostgreSQL first (persistent storage)
	if err := s.saveMessageToDB(msg); err != nil {
		return fmt.Errorf("failed to save message to database: %w", err)
	}

	// 2. Save to Redis (cache) - non-critical, log but don't fail
	if err := s.saveMessageToRedis(sessionID, msg); err != nil {
		fmt.Printf("⚠️ Failed to cache message to Redis (non-critical): %v\n", err)
	}

	return nil
}

// saveMessageToDB saves a message to PostgreSQL using Ent
func (s *MessageService) saveMessageToDB(msg *models.Message) error {
	// msg.SessionID already contains the correct UUID (foreign key to chat_sessions.id)
	// No need to look it up again

	// Convert products to proper format
	var productsJSON []map[string]interface{}
	if msg.Products != nil && len(msg.Products) > 0 {
		for _, product := range msg.Products {
			productMap := map[string]interface{}{
				"name":       product.Name,
				"price":      product.Price,
				"link":       product.Link,
				"image":      product.Image,
				"page_token": product.PageToken,
			}
			if product.OldPrice != "" {
				productMap["old_price"] = product.OldPrice
			}
			if product.Description != "" {
				productMap["description"] = product.Description
			}
			if product.Badge != "" {
				productMap["badge"] = product.Badge
			}
			productsJSON = append(productsJSON, productMap)
		}
	}

	// Create message in PostgreSQL
	createBuilder := s.client.Message.Create().
		SetID(msg.ID).
		SetSessionID(msg.SessionID).
		SetRole(msg.Role).
		SetContent(msg.Content).
		SetCreatedAt(msg.CreatedAt)

	// Set optional fields
	if msg.ResponseType != "" {
		createBuilder.SetResponseType(msg.ResponseType)
	}
	if len(msg.QuickReplies) > 0 {
		createBuilder.SetQuickReplies(msg.QuickReplies)
	}
	if len(productsJSON) > 0 {
		createBuilder.SetProducts(productsJSON)
	}
	if msg.SearchInfo != nil {
		createBuilder.SetSearchInfo(msg.SearchInfo)
	}

	_, err = createBuilder.Save(s.ctx)
	if err != nil {
		return fmt.Errorf("failed to create message in database: %w", err)
	}

	return nil
}

// saveMessageToRedis saves a message to Redis cache
func (s *MessageService) saveMessageToRedis(sessionID string, msg *models.Message) error {
	key := fmt.Sprintf(constants.CachePrefixMessages, sessionID)

	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	pipe := s.redis.Pipeline()
	pipe.RPush(s.ctx, key, data)
	pipe.Expire(s.ctx, key, s.ttl)

	_, err = pipe.Exec(s.ctx)
	return err
}

// getSessionUUIDBySessionID gets the session UUID from session_id string
func (s *MessageService) getSessionUUIDBySessionID(sessionID string) (uuid.UUID, error) {
	session, err := s.client.ChatSession.Query().
		Where(chatsession.SessionIDEQ(sessionID)).
		Only(s.ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return uuid.Nil, fmt.Errorf("session not found: %s", sessionID)
		}
		return uuid.Nil, fmt.Errorf("failed to query session: %w", err)
	}

	return session.ID, nil
}

// AddMessageInMemory adds a message using session object
// Saves to both PostgreSQL (persistent) and Redis (cache)
func (s *MessageService) AddMessageInMemory(session *models.ChatSession, msg *models.Message) error {
	// 1. Save to PostgreSQL first (persistent storage)
	if err := s.saveMessageToDB(msg); err != nil {
		return fmt.Errorf("failed to save message to database: %w", err)
	}

	// 2. Save to Redis (cache) - non-critical, log but don't fail
	if err := s.saveMessageToRedis(session.SessionID, msg); err != nil {
		fmt.Printf("⚠️ Failed to cache message to Redis (non-critical): %v\n", err)
	}

	return nil
}

// GetMessages retrieves all messages for a session
// Tries Redis first (cache), falls back to PostgreSQL (persistent storage)
func (s *MessageService) GetMessages(sessionID string) ([]*models.Message, error) {
	// Try Redis first (fast cache)
	messages, err := s.getMessagesFromRedis(sessionID)
	if err == nil && len(messages) > 0 {
		return messages, nil
	}

	// Redis miss or error - try PostgreSQL
	if err != nil {
		fmt.Printf("⚠️ Redis error when getting messages: %v, trying PostgreSQL\n", err)
	}

	// Get from PostgreSQL
	messages, err = s.getMessagesFromDB(sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages from database: %w", err)
	}

	// Restore to Redis for future requests
	if len(messages) > 0 {
		fmt.Printf("📦 Messages for session %s restored from PostgreSQL to Redis\n", sessionID)

		// Clear existing cache first to ensure consistency
		key := fmt.Sprintf(constants.CachePrefixMessages, sessionID)
		if err := s.redis.Del(s.ctx, key).Err(); err != nil {
			fmt.Printf("⚠️ Failed to clear old Redis cache: %v\n", err)
		}

		// Restore all messages in correct order
		for _, msg := range messages {
			if err := s.saveMessageToRedis(sessionID, msg); err != nil {
				fmt.Printf("⚠️ Failed to restore message to Redis: %v\n", err)
			}
		}
	}

	return messages, nil
}

// getMessagesFromRedis retrieves messages from Redis cache
func (s *MessageService) getMessagesFromRedis(sessionID string) ([]*models.Message, error) {
	key := fmt.Sprintf(constants.CachePrefixMessages, sessionID)

	data, err := s.redis.LRange(s.ctx, key, 0, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get messages from Redis: %w", err)
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

// getMessagesFromDB retrieves messages from PostgreSQL
func (s *MessageService) getMessagesFromDB(sessionID string) ([]*models.Message, error) {
	// Get session UUID
	sessionUUID, err := s.getSessionUUIDBySessionID(sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get session UUID: %w", err)
	}

	// Query messages from database
	entMessages, err := s.client.Message.Query().
		Where(message.SessionIDEQ(sessionUUID)).
		Order(ent.Asc(message.FieldCreatedAt)).
		All(s.ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to query messages: %w", err)
	}

	// Convert to models
	messages := make([]*models.Message, 0, len(entMessages))
	for _, entMsg := range entMessages {
		msg, err := convertEntMessageToModel(entMsg)
		if err != nil {
			fmt.Printf("⚠️ Failed to convert message %s: %v\n", entMsg.ID.String(), err)
			continue
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

// convertEntMessageToModel converts Ent Message to models.Message
func convertEntMessageToModel(entMsg *ent.Message) (*models.Message, error) {
	if entMsg == nil {
		return nil, nil
	}

	// Convert products from []map[string]interface{} to []ProductCard
	var products []models.ProductCard
	if entMsg.Products != nil {
		for _, productMap := range entMsg.Products {
			product := models.ProductCard{}

			if name, ok := productMap["name"].(string); ok {
				product.Name = name
			}
			if price, ok := productMap["price"].(string); ok {
				product.Price = price
			}
			if link, ok := productMap["link"].(string); ok {
				product.Link = link
			}
			if image, ok := productMap["image"].(string); ok {
				product.Image = image
			}
			if pageToken, ok := productMap["page_token"].(string); ok {
				product.PageToken = pageToken
			}
			if oldPrice, ok := productMap["old_price"].(string); ok {
				product.OldPrice = oldPrice
			}
			if description, ok := productMap["description"].(string); ok {
				product.Description = description
			}
			if badge, ok := productMap["badge"].(string); ok {
				product.Badge = badge
			}

			products = append(products, product)
		}
	}

	return &models.Message{
		ID:           entMsg.ID,
		SessionID:    entMsg.SessionID,
		Role:         entMsg.Role,
		Content:      entMsg.Content,
		ResponseType: entMsg.ResponseType,
		QuickReplies: entMsg.QuickReplies,
		Products:     products,
		SearchInfo:   entMsg.SearchInfo,
		CreatedAt:    entMsg.CreatedAt,
	}, nil
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

// GetMessagesSince retrieves all messages for a session created after a specific time
// This is useful for reconnection scenarios where client wants to catch up on missed messages
func (s *MessageService) GetMessagesSince(sessionID string, since time.Time) ([]*models.Message, error) {
	// Get all messages from database (source of truth)
	sessionUUID, err := s.getSessionUUIDBySessionID(sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get session UUID: %w", err)
	}

	// Query messages created after 'since' timestamp
	entMessages, err := s.client.Message.Query().
		Where(
			message.And(
				message.SessionIDEQ(sessionUUID),
				message.CreatedAtGT(since),
			),
		).
		Order(ent.Asc(message.FieldCreatedAt)).
		All(s.ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to query messages since %v: %w", since, err)
	}

	// Convert to models
	messages := make([]*models.Message, 0, len(entMessages))
	for _, entMsg := range entMessages {
		msg, err := convertEntMessageToModel(entMsg)
		if err != nil {
			fmt.Printf("⚠️ Failed to convert message %s: %v\n", entMsg.ID.String(), err)
			continue
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

// GetMessagesAfterID retrieves all messages created after a specific message ID
// Useful for pagination and reconnection scenarios
func (s *MessageService) GetMessagesAfterID(sessionID string, afterID uuid.UUID) ([]*models.Message, error) {
	// First get the timestamp of the reference message
	refMsg, err := s.client.Message.Get(s.ctx, afterID)
	if err != nil {
		if ent.IsNotFound(err) {
			// Message not found, return all messages for session
			return s.GetMessages(sessionID)
		}
		return nil, fmt.Errorf("failed to get reference message: %w", err)
	}

	// Get messages created after this timestamp
	return s.GetMessagesSince(sessionID, refMsg.CreatedAt)
}

// InvalidateMessageCache invalidates the Redis cache for a specific session's messages
// This should be called when messages are modified directly in PostgreSQL
func (s *MessageService) InvalidateMessageCache(sessionID string) error {
	key := fmt.Sprintf(constants.CachePrefixMessages, sessionID)
	return s.redis.Del(s.ctx, key).Err()
}

// RefreshMessageCache refreshes the Redis cache from PostgreSQL
// This ensures cache consistency after direct database modifications
func (s *MessageService) RefreshMessageCache(sessionID string) error {
	// Get fresh data from PostgreSQL
	messages, err := s.getMessagesFromDB(sessionID)
	if err != nil {
		return fmt.Errorf("failed to get messages from database: %w", err)
	}

	// Clear old cache
	key := fmt.Sprintf(constants.CachePrefixMessages, sessionID)
	if err := s.redis.Del(s.ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to clear Redis cache: %w", err)
	}

	// Restore to Redis in correct order
	for _, msg := range messages {
		if err := s.saveMessageToRedis(sessionID, msg); err != nil {
			return fmt.Errorf("failed to restore message to Redis: %w", err)
		}
	}

	fmt.Printf("✅ Refreshed message cache for session %s (%d messages)\n", sessionID, len(messages))
	return nil
}
