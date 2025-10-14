package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"mylittleprice/internal/constants"
	"mylittleprice/internal/models"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type SessionService struct {
	redis       *redis.Client
	ctx         context.Context
	ttl         time.Duration
	maxMsgs     int
	maxSearches int
}

func NewSessionService(redisClient *redis.Client, sessionTTL int, maxMessages int) *SessionService {
	return &SessionService{
		redis:       redisClient,
		ctx:         context.Background(),
		ttl:         time.Duration(sessionTTL) * time.Second,
		maxMsgs:     maxMessages,
		maxSearches: 999999,
	}
}

func (s *SessionService) CreateSession(sessionID, country, language string) (*models.ChatSession, error) {
	session := &models.ChatSession{
		ID:           uuid.New(),
		SessionID:    sessionID,
		CountryCode:  country,
		LanguageCode: language,
		Currency:     "",
		MessageCount: 0,
		SearchState: models.SearchState{
			Status:          models.SearchStatusIdle,
			Category:        "",
			ProductType:     "",
			Brand:           "",
			CollectedParams: []string{},
			SearchCount:     0,
			LastProduct:     nil,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		ExpiresAt: time.Now().Add(s.ttl),
	}

	if err := s.saveSession(session); err != nil {
		return nil, err
	}

	return session, nil
}

func (s *SessionService) GetSession(sessionID string) (*models.ChatSession, error) {
	key := fmt.Sprintf(constants.CachePrefixSession+"%s", sessionID)

	data, err := s.redis.Get(s.ctx, key).Bytes()
	if err == redis.Nil {
		return nil, fmt.Errorf("session not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	var session models.ChatSession
	err = json.Unmarshal(data, &session)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal session: %w", err)
	}

	return &session, nil
}

func (s *SessionService) UpdateSession(session *models.ChatSession) error {
	session.UpdatedAt = time.Now()
	return s.saveSession(session)
}

func (s *SessionService) saveSession(session *models.ChatSession) error {
	key := fmt.Sprintf(constants.CachePrefixSession+"%s", session.SessionID)

	data, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	err = s.redis.Set(s.ctx, key, data, s.ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to save session: %w", err)
	}

	return nil
}

func (s *SessionService) StartNewSearch(sessionID string) error {
	session, err := s.GetSession(sessionID)
	if err != nil {
		return err
	}

	session.SearchState = models.SearchState{
		Status:          models.SearchStatusIdle,
		Category:        "",
		ProductType:     "",
		Brand:           "",
		CollectedParams: []string{},
		SearchCount:     0,
		LastProduct:     nil,
	}

	return s.UpdateSession(session)
}

func (s *SessionService) SetCategory(sessionID, category string) error {
	session, err := s.GetSession(sessionID)
	if err != nil {
		return err
	}

	session.SearchState.Category = category
	return s.UpdateSession(session)
}

func (s *SessionService) IsSearchCompleted(sessionID string) bool {
	session, err := s.GetSession(sessionID)
	if err != nil {
		return false
	}
	return session.SearchState.Status == models.SearchStatusCompleted
}

func (s *SessionService) ResetSearchStatus(sessionID string) error {
	session, err := s.GetSession(sessionID)
	if err != nil {
		return err
	}

	session.SearchState.Status = models.SearchStatusIdle
	return s.UpdateSession(session)
}

func (s *SessionService) GetSessionInfo(sessionID string) map[string]interface{} {
	session, err := s.GetSession(sessionID)
	if err != nil {
		return map[string]interface{}{
			"error": err.Error(),
		}
	}

	info := map[string]interface{}{
		"session_id":    session.SessionID,
		"country":       session.CountryCode,
		"language":      session.LanguageCode,
		"currency":      session.Currency,
		"message_count": session.MessageCount,
		"search_state":  session.SearchState,
		"created_at":    session.CreatedAt,
		"updated_at":    session.UpdatedAt,
		"expires_at":    session.ExpiresAt,
		"ttl_seconds":   int(time.Until(session.ExpiresAt).Seconds()),
	}

	return info
}

func (s *SessionService) IncrementMessageCount(sessionID string) error {
	session, err := s.GetSession(sessionID)
	if err != nil {
		return err
	}

	session.MessageCount++
	return s.UpdateSession(session)
}

func (s *SessionService) CanSendMessage(sessionID string) (bool, error) {
	return true, nil
}

func (s *SessionService) AddMessage(sessionID string, message *models.Message) error {
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

func (s *SessionService) GetMessages(sessionID string) ([]*models.Message, error) {
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
			continue
		}
		messages = append(messages, &msg)
	}

	return messages, nil
}

func (s *SessionService) GetConversationHistory(sessionID string) ([]map[string]string, error) {
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

func (s *SessionService) GetRecentMessages(sessionID string, count int) ([]*models.Message, error) {
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
			continue
		}
		messages = append(messages, &msg)
	}

	return messages, nil
}

func (s *SessionService) DeleteSession(sessionID string) error {
	sessionKey := fmt.Sprintf(constants.CachePrefixSession+"%s", sessionID)
	messagesKey := fmt.Sprintf(constants.CachePrefixMessages, sessionID)

	pipe := s.redis.Pipeline()
	pipe.Del(s.ctx, sessionKey)
	pipe.Del(s.ctx, messagesKey)

	_, err := pipe.Exec(s.ctx)
	return err
}

func (s *SessionService) SetMaxSearches(max int) {
	s.maxSearches = max
}

func (s *SessionService) GetMaxSearches() int {
	return s.maxSearches
}

func (s *SessionService) GetSessionStats(sessionID string) (map[string]interface{}, error) {
	session, err := s.GetSession(sessionID)
	if err != nil {
		return nil, err
	}

	messages, _ := s.GetMessages(sessionID)

	stats := map[string]interface{}{
		"session_id":     session.SessionID,
		"country":        session.CountryCode,
		"language":       session.LanguageCode,
		"currency":       session.Currency,
		"message_count":  session.MessageCount,
		"search_count":   session.SearchState.SearchCount,
		"search_status":  session.SearchState.Status,
		"category":       session.SearchState.Category,
		"created_at":     session.CreatedAt,
		"updated_at":     session.UpdatedAt,
		"expires_at":     session.ExpiresAt,
		"ttl_seconds":    int(time.Until(session.ExpiresAt).Seconds()),
		"total_messages": len(messages),
	}

	return stats, nil
}
