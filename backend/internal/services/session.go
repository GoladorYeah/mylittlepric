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

// SessionService handles chat session management with optimized Redis operations
type SessionService struct {
	redis       *redis.Client
	ctx         context.Context
	ttl         time.Duration
	maxMsgs     int
	maxSearches int
}

// NewSessionService creates a new session service
func NewSessionService(redisClient *redis.Client, sessionTTL int, maxMessages int) *SessionService {
	return &SessionService{
		redis:       redisClient,
		ctx:         context.Background(),
		ttl:         time.Duration(sessionTTL) * time.Second,
		maxMsgs:     maxMessages,
		maxSearches: 3, // Will be configurable
	}
}

// CreateSession creates a new chat session
func (s *SessionService) CreateSession(sessionID, country, language string) (*models.ChatSession, error) {
	currency := getCurrencyForCountry(country)

	session := &models.ChatSession{
		ID:           uuid.New(),
		SessionID:    sessionID,
		CountryCode:  country,
		LanguageCode: language,
		Currency:     currency,
		MessageCount: 0,
		SearchState: models.SearchState{
			Status:          models.SearchStatusIdle,
			Category:        "",
			ProductType:     "",
			Brand:           "",
			CollectedParams: []string{},
			SearchCount:     0,
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

// GetSession retrieves a session from Redis
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

// UpdateSession updates session data in Redis (optimized with pipeline)
func (s *SessionService) UpdateSession(session *models.ChatSession) error {
	session.UpdatedAt = time.Now()
	return s.saveSession(session)
}

// saveSession saves session to Redis
func (s *SessionService) saveSession(session *models.ChatSession) error {
	key := fmt.Sprintf(constants.CachePrefixSession+"%s", session.SessionID)
	data, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	return s.redis.Set(s.ctx, key, data, s.ttl).Err()
}

// ═══════════════════════════════════════════════════════════
// SEARCH STATE MANAGEMENT (OPTIMIZED)
// ═══════════════════════════════════════════════════════════

// StartNewSearch resets search state for new product search
func (s *SessionService) StartNewSearch(sessionID string) error {
	session, err := s.GetSession(sessionID)
	if err != nil {
		return err
	}

	if session.SearchState.SearchCount >= s.maxSearches {
		return fmt.Errorf("maximum searches per session reached")
	}

	// Reset search state
	session.SearchState = models.SearchState{
		Status:          models.SearchStatusInProgress,
		Category:        "",
		ProductType:     "",
		Brand:           "",
		CollectedParams: []string{},
		SearchCount:     session.SearchState.SearchCount + 1,
		LastSearchTime:  time.Time{},
	}

	return s.UpdateSession(session)
}

// UpdateSearchState updates the current search state
func (s *SessionService) UpdateSearchState(sessionID string, update func(*models.SearchState)) error {
	session, err := s.GetSession(sessionID)
	if err != nil {
		return err
	}

	update(&session.SearchState)

	return s.UpdateSession(session)
}

// SetCategory sets category for current search
func (s *SessionService) SetCategory(sessionID, category string) error {
	return s.UpdateSearchState(sessionID, func(state *models.SearchState) {
		if state.Category == "" {
			state.Category = category
		}
	})
}

// SetProductType sets product type for current search
func (s *SessionService) SetProductType(sessionID, productType string) error {
	return s.UpdateSearchState(sessionID, func(state *models.SearchState) {
		if state.ProductType == "" {
			state.ProductType = productType
		}
	})
}

// SetBrand sets brand for current search
func (s *SessionService) SetBrand(sessionID, brand string) error {
	return s.UpdateSearchState(sessionID, func(state *models.SearchState) {
		if state.Brand == "" {
			state.Brand = brand
		}
	})
}

// AddCollectedParam adds a parameter to collected params
func (s *SessionService) AddCollectedParam(sessionID, param string) error {
	return s.UpdateSearchState(sessionID, func(state *models.SearchState) {
		// Avoid duplicates
		for _, p := range state.CollectedParams {
			if p == param {
				return
			}
		}
		state.CollectedParams = append(state.CollectedParams, param)
	})
}

// MarkSearchCompleted marks search as completed
func (s *SessionService) MarkSearchCompleted(sessionID string) error {
	return s.UpdateSearchState(sessionID, func(state *models.SearchState) {
		state.Status = models.SearchStatusCompleted
		state.LastSearchTime = time.Now()
	})
}

// CanStartNewSearch checks if user can start a new search
func (s *SessionService) CanStartNewSearch(sessionID string) (bool, string) {
	session, err := s.GetSession(sessionID)
	if err != nil {
		return false, "Session not found"
	}

	if session.SearchState.SearchCount >= s.maxSearches {
		return false, fmt.Sprintf("Maximum %d searches per session reached. Please start a new session.", s.maxSearches)
	}

	return true, ""
}

// IsSearchInProgress checks if there's an active search
func (s *SessionService) IsSearchInProgress(sessionID string) bool {
	session, err := s.GetSession(sessionID)
	if err != nil {
		return false
	}

	return session.SearchState.Status == models.SearchStatusInProgress
}

// IsSearchCompleted checks if search was completed
func (s *SessionService) IsSearchCompleted(sessionID string) bool {
	session, err := s.GetSession(sessionID)
	if err != nil {
		return false
	}

	return session.SearchState.Status == models.SearchStatusCompleted
}

// GetSearchStateInfo returns search state info for frontend
func (s *SessionService) GetSearchStateInfo(sessionID string) *models.SearchStateInfo {
	session, err := s.GetSession(sessionID)
	if err != nil {
		return nil
	}

	info := &models.SearchStateInfo{
		Status:      session.SearchState.Status,
		Category:    session.SearchState.Category,
		SearchCount: session.SearchState.SearchCount,
		MaxSearches: s.maxSearches,
		CanContinue: session.SearchState.SearchCount < s.maxSearches,
	}

	// Set appropriate message
	switch session.SearchState.Status {
	case models.SearchStatusCompleted:
		info.Message = "Search completed. Start a new search to find another product."
	case models.SearchStatusBlocked:
		info.Message = fmt.Sprintf("Maximum %d searches reached. Please create a new session.", s.maxSearches)
	case models.SearchStatusInProgress:
		info.Message = "Collecting product information..."
	case models.SearchStatusIdle:
		info.Message = "Ready to start searching!"
	}

	return info
}

// ═══════════════════════════════════════════════════════════
// MESSAGE MANAGEMENT (OPTIMIZED WITH PIPELINE)
// ═══════════════════════════════════════════════════════════

// IncrementMessageCount increments the message count
func (s *SessionService) IncrementMessageCount(sessionID string) error {
	session, err := s.GetSession(sessionID)
	if err != nil {
		return err
	}

	session.MessageCount++
	return s.UpdateSession(session)
}

// CanSendMessage checks if user can send more messages
func (s *SessionService) CanSendMessage(sessionID string) (bool, error) {
	session, err := s.GetSession(sessionID)
	if err != nil {
		return false, err
	}

	return session.MessageCount < s.maxMsgs, nil
}

// AddMessage adds a message to session history (optimized)
func (s *SessionService) AddMessage(sessionID string, message *models.Message) error {
	key := fmt.Sprintf(constants.CachePrefixMessages, sessionID)

	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	// Use pipeline for atomic operations
	pipe := s.redis.Pipeline()
	pipe.RPush(s.ctx, key, data)
	pipe.Expire(s.ctx, key, s.ttl)

	_, err = pipe.Exec(s.ctx)
	return err
}

// GetMessages retrieves all messages for a session
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

// GetConversationHistory returns formatted conversation for AI
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

// GetRecentMessages retrieves last N messages
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

// ═══════════════════════════════════════════════════════════
// SESSION CLEANUP (OPTIMIZED WITH PIPELINE)
// ═══════════════════════════════════════════════════════════

// DeleteSession deletes a session and its messages atomically
func (s *SessionService) DeleteSession(sessionID string) error {
	sessionKey := fmt.Sprintf(constants.CachePrefixSession+"%s", sessionID)
	messagesKey := fmt.Sprintf(constants.CachePrefixMessages, sessionID)

	pipe := s.redis.Pipeline()
	pipe.Del(s.ctx, sessionKey)
	pipe.Del(s.ctx, messagesKey)

	_, err := pipe.Exec(s.ctx)
	return err
}

// ═══════════════════════════════════════════════════════════
// HELPER METHODS
// ═══════════════════════════════════════════════════════════

// getCurrencyForCountry maps country code to currency
func getCurrencyForCountry(country string) string {
	currencyMap := map[string]string{
		"CH": "CHF", "DE": "EUR", "AT": "EUR", "FR": "EUR",
		"IT": "EUR", "ES": "EUR", "PT": "EUR", "NL": "EUR",
		"BE": "EUR", "PL": "PLN", "CZ": "CZK", "SE": "SEK",
		"NO": "NOK", "DK": "DKK", "FI": "EUR", "GB": "GBP",
		"US": "USD",
	}
	if currency, ok := currencyMap[country]; ok {
		return currency
	}
	return "EUR"
}

// GetLanguageForCountry returns default language for country
func (s *SessionService) GetLanguageForCountry(country string) string {
	languageMap := map[string]string{
		"CH": "de", "DE": "de", "AT": "de",
		"FR": "fr", "IT": "it", "ES": "es",
		"PT": "pt", "NL": "nl", "BE": "nl",
		"PL": "pl", "CZ": "cs", "SE": "sv",
		"NO": "no", "DK": "da", "FI": "fi",
		"GB": "en", "US": "en",
	}

	if language, ok := languageMap[country]; ok {
		return language
	}
	return "en"
}

// SetMaxSearches sets maximum searches per session
func (s *SessionService) SetMaxSearches(max int) {
	s.maxSearches = max
}

// GetMaxSearches returns current max searches setting
func (s *SessionService) GetMaxSearches() int {
	return s.maxSearches
}

// GetSessionStats returns statistics about a session
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
