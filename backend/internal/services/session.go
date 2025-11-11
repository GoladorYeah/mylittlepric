package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"mylittleprice/ent"
	"mylittleprice/ent/chatsession"
	"mylittleprice/ent/user"
	"mylittleprice/internal/constants"
	"mylittleprice/internal/models"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type SessionService struct {
	redis              *redis.Client
	client             *ent.Client
	authService        *AuthService
	ctx                context.Context
	ttl                time.Duration
	maxMsgs            int
	maxSearches        int
	universalPromptMgr *UniversalPromptManager
}

func NewSessionService(redisClient *redis.Client, client *ent.Client, sessionTTL int, maxMessages int) *SessionService {
	return &SessionService{
		redis:              redisClient,
		client:             client,
		authService:        nil, // Will be set later via SetAuthService
		ctx:                context.Background(),
		ttl:                time.Duration(sessionTTL) * time.Second,
		maxMsgs:            maxMessages,
		maxSearches:        999999,
		universalPromptMgr: NewUniversalPromptManager(),
	}
}

// SetAuthService sets the AuthService (used to avoid circular dependency)
func (s *SessionService) SetAuthService(authService *AuthService) {
	s.authService = authService
}

func (s *SessionService) CreateSession(sessionID, country, language, currency string) (*models.ChatSession, error) {
	return s.CreateSessionWithUser(sessionID, country, language, currency, nil)
}

func (s *SessionService) CreateSessionWithUser(sessionID, country, language, currency string, userID *uuid.UUID) (*models.ChatSession, error) {
	session := &models.ChatSession{
		ID:           uuid.New(),
		SessionID:    sessionID,
		UserID:       userID,
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
	// Validate input
	if err := validateSessionID(sessionID); err != nil {
		return nil, fmt.Errorf("invalid session ID: %w", err)
	}

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

// getSessionFromDB retrieves session from PostgreSQL using Ent
func (s *SessionService) getSessionFromDB(sessionID string) (*models.ChatSession, error) {
	entSession, err := s.client.ChatSession.Query().
		Where(chatsession.SessionIDEQ(sessionID)).
		Only(s.ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("session not found in database")
		}
		return nil, fmt.Errorf("failed to get session from database: %w", err)
	}

	return convertEntSessionToModel(entSession)
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

// saveSessionToDB saves or updates session in PostgreSQL using Ent
func (s *SessionService) saveSessionToDB(session *models.ChatSession) error {
	// If session has a user_id, ensure user exists in PostgreSQL first
	if session.UserID != nil {
		if err := s.ensureUserExistsInPostgres(*session.UserID); err != nil {
			return fmt.Errorf("failed to ensure user exists: %w", err)
		}
	}

	// Convert SearchState to map
	searchStateMap, err := structToMap(session.SearchState)
	if err != nil {
		return fmt.Errorf("failed to convert search_state: %w", err)
	}

	// Convert CycleState to map
	cycleStateMap, err := structToMap(session.CycleState)
	if err != nil {
		return fmt.Errorf("failed to convert cycle_state: %w", err)
	}

	// Convert ConversationContext to map (if present)
	var conversationContextMap map[string]interface{}
	if session.ConversationContext != nil {
		conversationContextMap, err = structToMap(session.ConversationContext)
		if err != nil {
			return fmt.Errorf("failed to convert conversation_context: %w", err)
		}
	}

	// Check if session exists
	exists, err := s.client.ChatSession.Query().
		Where(chatsession.SessionIDEQ(session.SessionID)).
		Exist(s.ctx)

	if err != nil {
		return fmt.Errorf("failed to check session existence: %w", err)
	}

	if exists {
		// Update existing session
		updateBuilder := s.client.ChatSession.Update().
			Where(chatsession.SessionIDEQ(session.SessionID)).
			SetCountryCode(session.CountryCode).
			SetLanguageCode(session.LanguageCode).
			SetCurrency(session.Currency).
			SetMessageCount(session.MessageCount).
			SetSearchState(searchStateMap).
			SetCycleState(cycleStateMap).
			SetUpdatedAt(session.UpdatedAt).
			SetExpiresAt(session.ExpiresAt)

		// Set optional user_id
		if session.UserID != nil {
			updateBuilder.SetUserID(*session.UserID)
		} else {
			updateBuilder.ClearUserID()
		}

		// Set optional conversation_context
		if conversationContextMap != nil {
			updateBuilder.SetConversationContext(conversationContextMap)
		} else {
			updateBuilder.ClearConversationContext()
		}

		_, err = updateBuilder.Save(s.ctx)
		if err != nil {
			return fmt.Errorf("failed to update session: %w", err)
		}
	} else {
		// Create new session
		createBuilder := s.client.ChatSession.Create().
			SetID(session.ID).
			SetSessionID(session.SessionID).
			SetCountryCode(session.CountryCode).
			SetLanguageCode(session.LanguageCode).
			SetCurrency(session.Currency).
			SetMessageCount(session.MessageCount).
			SetSearchState(searchStateMap).
			SetCycleState(cycleStateMap).
			SetCreatedAt(session.CreatedAt).
			SetUpdatedAt(session.UpdatedAt).
			SetExpiresAt(session.ExpiresAt)

		// Set optional user_id
		if session.UserID != nil {
			createBuilder.SetUserID(*session.UserID)
		}

		// Set optional conversation_context
		if conversationContextMap != nil {
			createBuilder.SetConversationContext(conversationContextMap)
		}

		_, err = createBuilder.Save(s.ctx)
		if err != nil {
			return fmt.Errorf("failed to create session: %w", err)
		}
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
			fmt.Printf("⚠️ Failed to unmarshal message in session %s: %v\n", sessionID, err)
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
			fmt.Printf("⚠️ Failed to unmarshal recent message in session %s: %v\n", sessionID, err)
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

// GetActiveSessionForUser returns the most recent active session for a user
// Returns nil if no active session found (not an error - user can start a new session)
func (s *SessionService) GetActiveSessionForUser(userID uuid.UUID) (*models.ChatSession, error) {
	// Get most recent session for this user that hasn't expired
	// Order by updated_at DESC to get the latest active session
	entSession, err := s.client.ChatSession.Query().
		Where(
			chatsession.And(
				chatsession.UserIDEQ(userID),
				chatsession.ExpiresAtGT(time.Now()),
			),
		).
		Order(ent.Desc(chatsession.FieldUpdatedAt)).
		First(s.ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			// No active session found - this is OK, user can start fresh
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get active session for user: %w", err)
	}

	session, err := convertEntSessionToModel(entSession)
	if err != nil {
		return nil, err
	}

	// Update Redis cache for faster subsequent access
	if err := s.saveSessionToRedis(session); err != nil {
		fmt.Printf("⚠️ Failed to cache active session to Redis: %v\n", err)
	}

	return session, nil
}

// GetSessionWithOngoingSearch returns a specific session if it has an ongoing search
// This is used for cross-device search continuity
// Returns nil if session not found or search is not in progress
func (s *SessionService) GetSessionWithOngoingSearch(sessionID string) (*models.ChatSession, error) {
	session, err := s.GetSession(sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	if session == nil {
		return nil, nil
	}

	// Check if this session has an ongoing search
	if session.SearchState.Status == models.SearchStatusInProgress {
		return session, nil
	}

	// Session exists but no ongoing search
	return nil, nil
}

// LinkSessionToUser links an existing session to a user (for when anonymous user logs in)
// If the session doesn't exist, creates a new one with default preferences
func (s *SessionService) LinkSessionToUser(sessionID string, userID uuid.UUID) error {
	session, err := s.GetSession(sessionID)
	if err != nil {
		// Session doesn't exist on backend - create it with default settings
		fmt.Printf("📝 Session %s not found, creating new session for user %s\n", sessionID, userID.String())

		// Use default locale settings (can be updated later via preferences)
		_, err := s.CreateSessionWithUser(sessionID, "US", "en", "USD", &userID)
		if err != nil {
			return fmt.Errorf("failed to create session for user: %w", err)
		}
		return nil
	}

	// Session exists - link it to the user
	session.UserID = &userID
	return s.UpdateSession(session)
}

// ensureUserExistsInPostgres checks if user exists in PostgreSQL and syncs from Redis if needed
// This prevents foreign key constraint violations when creating sessions for authenticated users
func (s *SessionService) ensureUserExistsInPostgres(userID uuid.UUID) error {
	if s.authService == nil {
		return fmt.Errorf("authService not set in SessionService")
	}

	// Check if user exists in PostgreSQL using Ent
	exists, err := s.client.User.Query().
		Where(user.IDEQ(userID)).
		Exist(s.ctx)

	if err != nil {
		return fmt.Errorf("failed to check if user exists: %w", err)
	}

	if exists {
		// User already exists in PostgreSQL
		return nil
	}

	// User doesn't exist in PostgreSQL - sync from Redis
	fmt.Printf("⚠️ User %s not found in PostgreSQL, syncing from Redis...\n", userID.String())

	userModel, err := s.authService.GetUserByID(userID)
	if err != nil {
		return fmt.Errorf("failed to get user from Redis: %w", err)
	}

	if err := s.authService.SaveUserToPostgres(userModel); err != nil {
		return fmt.Errorf("failed to save user to PostgreSQL: %w", err)
	}

	fmt.Printf("✅ Synced user %s to PostgreSQL\n", userID.String())
	return nil
}

// ==================== Helper Functions ====================

// convertEntSessionToModel converts Ent ChatSession to models.ChatSession
func convertEntSessionToModel(entSession *ent.ChatSession) (*models.ChatSession, error) {
	if entSession == nil {
		return nil, nil
	}

	// Convert search_state from map to SearchState
	var searchState models.SearchState
	if entSession.SearchState != nil {
		if err := mapToStruct(entSession.SearchState, &searchState); err != nil {
			return nil, fmt.Errorf("failed to convert search_state: %w", err)
		}
	}

	// Convert cycle_state from map to CycleState
	var cycleState models.CycleState
	if entSession.CycleState != nil {
		if err := mapToStruct(entSession.CycleState, &cycleState); err != nil {
			return nil, fmt.Errorf("failed to convert cycle_state: %w", err)
		}
	}

	// Convert conversation_context from map to ConversationContext (if present)
	var conversationContext *models.ConversationContext
	if entSession.ConversationContext != nil {
		conversationContext = &models.ConversationContext{}
		if err := mapToStruct(entSession.ConversationContext, conversationContext); err != nil {
			return nil, fmt.Errorf("failed to convert conversation_context: %w", err)
		}
	}

	// Convert user_id
	var userID *uuid.UUID
	if entSession.UserID != uuid.Nil {
		userID = &entSession.UserID
	}

	return &models.ChatSession{
		ID:                  entSession.ID,
		SessionID:           entSession.SessionID,
		UserID:              userID,
		CountryCode:         entSession.CountryCode,
		LanguageCode:        entSession.LanguageCode,
		Currency:            entSession.Currency,
		MessageCount:        entSession.MessageCount,
		SearchState:         searchState,
		CycleState:          cycleState,
		ConversationContext: conversationContext,
		CreatedAt:           entSession.CreatedAt,
		UpdatedAt:           entSession.UpdatedAt,
		ExpiresAt:           entSession.ExpiresAt,
	}, nil
}

// structToMap converts a struct to map[string]interface{} via JSON
func structToMap(v interface{}) (map[string]interface{}, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// mapToStruct converts map[string]interface{} to a struct via JSON
func mapToStruct(m map[string]interface{}, v interface{}) error {
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, v)
}
