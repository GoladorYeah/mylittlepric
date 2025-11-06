package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"mylittleprice/internal/constants"
	"mylittleprice/internal/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type SessionService struct {
	redis              *redis.Client
	db                 *sqlx.DB
	ctx                context.Context
	ttl                time.Duration
	maxMsgs            int
	maxSearches        int
	universalPromptMgr *UniversalPromptManager
}

func NewSessionService(redisClient *redis.Client, db *sqlx.DB, sessionTTL int, maxMessages int) *SessionService {
	return &SessionService{
		redis:              redisClient,
		db:                 db,
		ctx:                context.Background(),
		ttl:                time.Duration(sessionTTL) * time.Second,
		maxMsgs:            maxMessages,
		maxSearches:        999999,
		universalPromptMgr: NewUniversalPromptManager(),
	}
}

func (s *SessionService) CreateSession(sessionID, country, language, currency string) (*models.ChatSession, error) {
	session := &models.ChatSession{
		ID:           uuid.New(),
		SessionID:    sessionID,
		CountryCode:  country,
		LanguageCode: language,
		Currency:     currency,
		MessageCount: 0,
		SearchState: models.SearchState{
			Status:      models.SearchStatusIdle,
			Category:    "",
			SearchCount: 0,
			LastProduct: nil,
		},
		CycleState: s.universalPromptMgr.InitializeCycleState(),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		ExpiresAt:  time.Now().Add(s.ttl),
	}

	if err := s.saveSession(session); err != nil {
		return nil, err
	}

	return session, nil
}

func (s *SessionService) GetSession(sessionID string) (*models.ChatSession, error) {
	// Try Redis first (fast cache)
	key := fmt.Sprintf(constants.CachePrefixSession+"%s", sessionID)

	data, err := s.redis.Get(s.ctx, key).Bytes()
	if err == nil {
		// Found in Redis - unmarshal and return
		var session models.ChatSession
		err = json.Unmarshal(data, &session)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal session from Redis: %w", err)
		}
		return &session, nil
	}

	// Not in Redis - try PostgreSQL (persistent storage)
	if err != redis.Nil {
		// Redis error (not just "not found") - log but continue to DB
		fmt.Printf("⚠️ Redis error when getting session: %v\n", err)
	}

	// Try to load from PostgreSQL
	session, err := s.getSessionFromDB(sessionID)
	if err != nil {
		return nil, fmt.Errorf("session not found in Redis or PostgreSQL")
	}

	// Found in DB - restore to Redis for future requests
	fmt.Printf("📦 Session %s restored from PostgreSQL to Redis\n", sessionID)
	if err := s.saveSessionToRedis(session); err != nil {
		fmt.Printf("⚠️ Failed to restore session to Redis: %v\n", err)
	}

	return session, nil
}

// getSessionFromDB retrieves session from PostgreSQL
func (s *SessionService) getSessionFromDB(sessionID string) (*models.ChatSession, error) {
	var session models.ChatSession

	query := `SELECT id, session_id, country_code, language_code, currency, message_count,
	          search_state, cycle_state, conversation_context, created_at, updated_at, expires_at
	          FROM chat_sessions WHERE session_id = $1`

	err := s.db.Get(&session, query, sessionID)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("session not found in database")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get session from database: %w", err)
	}

	return &session, nil
}

// saveSessionToRedis saves session to Redis only
func (s *SessionService) saveSessionToRedis(session *models.ChatSession) error {
	key := fmt.Sprintf(constants.CachePrefixSession+"%s", session.SessionID)

	data, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	err = s.redis.Set(s.ctx, key, data, s.ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to save session to Redis: %w", err)
	}

	return nil
}

// saveSessionToDB saves or updates session in PostgreSQL
func (s *SessionService) saveSessionToDB(session *models.ChatSession) error {
	query := `INSERT INTO chat_sessions
	          (id, session_id, country_code, language_code, currency, message_count,
	           search_state, cycle_state, conversation_context, created_at, updated_at, expires_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	          ON CONFLICT (session_id) DO UPDATE SET
	            country_code = EXCLUDED.country_code,
	            language_code = EXCLUDED.language_code,
	            currency = EXCLUDED.currency,
	            message_count = EXCLUDED.message_count,
	            search_state = EXCLUDED.search_state,
	            cycle_state = EXCLUDED.cycle_state,
	            conversation_context = EXCLUDED.conversation_context,
	            updated_at = EXCLUDED.updated_at,
	            expires_at = EXCLUDED.expires_at`

	// Convert structs to JSONB
	searchStateJSON, err := json.Marshal(session.SearchState)
	if err != nil {
		return fmt.Errorf("failed to marshal search_state: %w", err)
	}

	cycleStateJSON, err := json.Marshal(session.CycleState)
	if err != nil {
		return fmt.Errorf("failed to marshal cycle_state: %w", err)
	}

	// Handle conversation_context as nullable JSONB
	var conversationContextParam interface{}
	if session.ConversationContext != nil {
		conversationContextJSON, err := json.Marshal(session.ConversationContext)
		if err != nil {
			return fmt.Errorf("failed to marshal conversation_context: %w", err)
		}
		conversationContextParam = conversationContextJSON
	} else {
		conversationContextParam = nil
	}

	_, err = s.db.Exec(query,
		session.ID,
		session.SessionID,
		session.CountryCode,
		session.LanguageCode,
		session.Currency,
		session.MessageCount,
		searchStateJSON,
		cycleStateJSON,
		conversationContextParam,
		session.CreatedAt,
		session.UpdatedAt,
		session.ExpiresAt,
	)

	if err != nil {
		return fmt.Errorf("failed to save session to database: %w", err)
	}

	return nil
}

func (s *SessionService) UpdateSession(session *models.ChatSession) error {
	session.UpdatedAt = time.Now()
	return s.saveSession(session)
}

func (s *SessionService) saveSession(session *models.ChatSession) error {
	// Save to both Redis (cache) and PostgreSQL (persistent storage)

	// Save to PostgreSQL first (persistent)
	if err := s.saveSessionToDB(session); err != nil {
		return fmt.Errorf("failed to save session to database: %w", err)
	}

	// Save to Redis (cache) - non-critical, log errors but don't fail
	if err := s.saveSessionToRedis(session); err != nil {
		fmt.Printf("⚠️ Failed to save session to Redis (non-critical): %v\n", err)
	}

	return nil
}

func (s *SessionService) StartNewSearch(sessionID string) error {
	session, err := s.GetSession(sessionID)
	if err != nil {
		return err
	}

	session.SearchState = models.SearchState{
		Status:      models.SearchStatusIdle,
		Category:    "",
		SearchCount: 0,
		LastProduct: nil,
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
		"cycle_id":       session.CycleState.CycleID,
		"iteration":      session.CycleState.Iteration,
		"prompt_id":      session.CycleState.PromptID,
		"created_at":     session.CreatedAt,
		"updated_at":     session.UpdatedAt,
		"expires_at":     session.ExpiresAt,
		"ttl_seconds":    int(time.Until(session.ExpiresAt).Seconds()),
		"total_messages": len(messages),
	}

	return stats, nil
}

// GetUniversalPromptManager returns the universal prompt manager
func (s *SessionService) GetUniversalPromptManager() *UniversalPromptManager {
	return s.universalPromptMgr
}

// IncrementCycleIteration increments the iteration and checks if we need a new cycle
// Returns true if we should start a new cycle
func (s *SessionService) IncrementCycleIteration(sessionID string) (bool, error) {
	session, err := s.GetSession(sessionID)
	if err != nil {
		return false, err
	}

	shouldStartNew := !s.universalPromptMgr.IncrementIteration(&session.CycleState)

	if err := s.UpdateSession(session); err != nil {
		return false, err
	}

	return shouldStartNew, nil
}

// StartNewCycle starts a new cycle with context carryover
func (s *SessionService) StartNewCycle(sessionID, lastRequest string, products []models.ProductInfo) error {
	session, err := s.GetSession(sessionID)
	if err != nil {
		return err
	}

	s.universalPromptMgr.StartNewCycle(&session.CycleState, lastRequest, products)

	return s.UpdateSession(session)
}

// AddToCycleHistory adds a message to the current cycle history
func (s *SessionService) AddToCycleHistory(sessionID, role, content string) error {
	session, err := s.GetSession(sessionID)
	if err != nil {
		return err
	}

	s.universalPromptMgr.AddToCycleHistory(&session.CycleState, role, content)

	return s.UpdateSession(session)
}
