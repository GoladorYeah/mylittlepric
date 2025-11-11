package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"mylittleprice/ent"
	"mylittleprice/ent/userpreference"
	"mylittleprice/internal/models"
)

var (
	ErrPreferencesNotFound = errors.New("user preferences not found")
)

type PreferencesService struct {
	client      *ent.Client
	authService *AuthService
	ctx         context.Context
}

func NewPreferencesService(client *ent.Client, authService *AuthService) *PreferencesService {
	return &PreferencesService{
		client:      client,
		authService: authService,
		ctx:         context.Background(),
	}
}

// GetUserPreferences retrieves user preferences from PostgreSQL
// Returns nil if preferences not found (first time user)
func (s *PreferencesService) GetUserPreferences(userID uuid.UUID) (*models.UserPreferences, error) {
	prefs, err := s.client.UserPreference.Query().
		Where(userpreference.UserIDEQ(userID)).
		Only(s.ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			// Not found is OK - return nil (user hasn't set preferences yet)
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user preferences: %w", err)
	}

	// Convert Ent entity to models
	return convertEntPreferencesToModel(prefs), nil
}

// UpsertUserPreferences creates or updates user preferences
// Only updates fields that are provided (non-nil in update struct)
func (s *PreferencesService) UpsertUserPreferences(userID uuid.UUID, update *models.UserPreferencesUpdate) (*models.UserPreferences, error) {
	// First, try to get existing preferences
	existing, err := s.GetUserPreferences(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing preferences: %w", err)
	}

	// If no existing preferences, insert new row
	if existing == nil {
		return s.createPreferences(userID, update)
	}

	// Otherwise, update existing row
	return s.updatePreferences(userID, update)
}

// createPreferences inserts new preferences row
func (s *PreferencesService) createPreferences(userID uuid.UUID, update *models.UserPreferencesUpdate) (*models.UserPreferences, error) {
	builder := s.client.UserPreference.Create().
		SetUserID(userID)

	// Set optional fields
	if update.Country != nil {
		builder.SetCountry(*update.Country)
	}
	if update.Currency != nil {
		builder.SetCurrency(*update.Currency)
	}
	if update.Language != nil {
		builder.SetLanguage(*update.Language)
	}
	if update.Theme != nil {
		builder.SetTheme(*update.Theme)
	}
	if update.SidebarOpen != nil {
		builder.SetSidebarOpen(*update.SidebarOpen)
	}
	if update.LastActiveSessionID != nil {
		builder.SetLastActiveSessionID(*update.LastActiveSessionID)
	}
	if update.SavedSearch != nil {
		// Convert SavedSearch to map[string]interface{}
		savedSearchMap := savedSearchToMap(update.SavedSearch)
		builder.SetSavedSearch(savedSearchMap)
	}

	prefs, err := builder.Save(s.ctx)
	if err != nil {
		// Check if it's a foreign key constraint error (user doesn't exist in PostgreSQL)
		if isForeignKeyError(err) {
			fmt.Printf("⚠️ User %s not found in PostgreSQL, attempting to sync from Redis...\n", userID.String())

			// Ensure user exists in PostgreSQL by fetching from Redis and saving
			if syncErr := s.ensureUserExistsInPostgres(userID); syncErr != nil {
				return nil, fmt.Errorf("failed to sync user to PostgreSQL: %w", syncErr)
			}

			// Retry the insert
			prefs, err = builder.Save(s.ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to create user preferences after sync: %w", err)
			}

			fmt.Printf("✅ Created preferences for user %s (after sync)\n", userID.String())
			return convertEntPreferencesToModel(prefs), nil
		}

		return nil, fmt.Errorf("failed to create user preferences: %w", err)
	}

	fmt.Printf("✅ Created preferences for user %s\n", userID.String())
	return convertEntPreferencesToModel(prefs), nil
}

// updatePreferences updates existing preferences (only provided fields)
func (s *PreferencesService) updatePreferences(userID uuid.UUID, update *models.UserPreferencesUpdate) (*models.UserPreferences, error) {
	builder := s.client.UserPreference.Update().
		Where(userpreference.UserIDEQ(userID))

	// Only update provided fields
	if update.Country != nil {
		builder.SetCountry(*update.Country)
	}
	if update.Currency != nil {
		builder.SetCurrency(*update.Currency)
	}
	if update.Language != nil {
		builder.SetLanguage(*update.Language)
	}
	if update.Theme != nil {
		builder.SetTheme(*update.Theme)
	}
	if update.SidebarOpen != nil {
		builder.SetSidebarOpen(*update.SidebarOpen)
	}
	if update.LastActiveSessionID != nil {
		builder.SetLastActiveSessionID(*update.LastActiveSessionID)
	}
	if update.SavedSearch != nil {
		// Convert SavedSearch to map[string]interface{}
		savedSearchMap := savedSearchToMap(update.SavedSearch)
		builder.SetSavedSearch(savedSearchMap)
	}

	_, err := builder.Save(s.ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to update user preferences: %w", err)
	}

	fmt.Printf("✅ Updated preferences for user %s\n", userID.String())

	// Return updated preferences
	return s.GetUserPreferences(userID)
}

// DeleteUserPreferences removes user preferences (e.g., on account deletion)
func (s *PreferencesService) DeleteUserPreferences(userID uuid.UUID) error {
	rowsAffected, err := s.client.UserPreference.Delete().
		Where(userpreference.UserIDEQ(userID)).
		Exec(s.ctx)

	if err != nil {
		return fmt.Errorf("failed to delete user preferences: %w", err)
	}

	if rowsAffected == 0 {
		return ErrPreferencesNotFound
	}

	fmt.Printf("✅ Deleted preferences for user %s\n", userID.String())
	return nil
}

// ==================== Search Synchronization ====================

// UpdateLastActiveSession updates the last active session ID for cross-device search continuity
func (s *PreferencesService) UpdateLastActiveSession(userID uuid.UUID, sessionID string) error {
	var sessionIDPtr *string
	if sessionID != "" {
		sessionIDPtr = &sessionID
	}

	update := &models.UserPreferencesUpdate{
		LastActiveSessionID: sessionIDPtr,
	}

	_, err := s.UpsertUserPreferences(userID, update)
	if err != nil {
		return fmt.Errorf("failed to update last active session: %w", err)
	}

	if sessionID == "" {
		fmt.Printf("✅ Cleared active session for user %s\n", userID.String())
	} else {
		fmt.Printf("✅ Updated active session to %s for user %s\n", sessionID, userID.String())
	}
	return nil
}

// GetLastActiveSession retrieves the session ID with an ongoing search
func (s *PreferencesService) GetLastActiveSession(userID uuid.UUID) (string, error) {
	prefs, err := s.GetUserPreferences(userID)
	if err != nil {
		return "", fmt.Errorf("failed to get user preferences: %w", err)
	}

	if prefs == nil || prefs.LastActiveSessionID == nil {
		return "", nil
	}

	return *prefs.LastActiveSessionID, nil
}

// ==================== Saved Search Synchronization ====================

// UpdateSavedSearch updates the saved search for cross-device synchronization
func (s *PreferencesService) UpdateSavedSearch(userID uuid.UUID, savedSearch *models.SavedSearch) error {
	update := &models.UserPreferencesUpdate{
		SavedSearch: savedSearch,
	}

	_, err := s.UpsertUserPreferences(userID, update)
	if err != nil {
		return fmt.Errorf("failed to update saved search: %w", err)
	}

	if savedSearch == nil {
		fmt.Printf("✅ Cleared saved search for user %s\n", userID.String())
	} else {
		fmt.Printf("✅ Updated saved search for user %s (session: %s, %d messages)\n",
			userID.String(), savedSearch.SessionID, len(savedSearch.Messages))
	}
	return nil
}

// GetSavedSearch retrieves the saved search
func (s *PreferencesService) GetSavedSearch(userID uuid.UUID) (*models.SavedSearch, error) {
	prefs, err := s.GetUserPreferences(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user preferences: %w", err)
	}

	if prefs == nil || prefs.SavedSearch == nil {
		return nil, nil
	}

	return prefs.SavedSearch, nil
}

// ==================== Helper Methods ====================

// Convert Ent UserPreference to models.UserPreferences
func convertEntPreferencesToModel(prefs *ent.UserPreference) *models.UserPreferences {
	if prefs == nil {
		return nil
	}

	var country, currency, language, theme *string
	var sidebarOpen *bool
	var lastActiveSessionID *string

	if prefs.Country != nil {
		country = prefs.Country
	}
	if prefs.Currency != nil {
		currency = prefs.Currency
	}
	if prefs.Language != nil {
		language = prefs.Language
	}
	if prefs.Theme != nil {
		theme = prefs.Theme
	}
	if prefs.SidebarOpen != nil {
		sidebarOpen = prefs.SidebarOpen
	}
	if prefs.LastActiveSessionID != nil {
		lastActiveSessionID = prefs.LastActiveSessionID
	}

	// Convert saved_search from map to SavedSearch model
	var savedSearch *models.SavedSearch
	if prefs.SavedSearch != nil {
		savedSearch = mapToSavedSearch(prefs.SavedSearch)
	}

	return &models.UserPreferences{
		UserID:              prefs.UserID,
		Country:             country,
		Currency:            currency,
		Language:            language,
		Theme:               theme,
		SidebarOpen:         sidebarOpen,
		LastActiveSessionID: lastActiveSessionID,
		SavedSearch:         savedSearch,
		CreatedAt:           prefs.CreatedAt,
		UpdatedAt:           prefs.UpdatedAt,
	}
}

// Convert SavedSearch to map[string]interface{}
func savedSearchToMap(ss *models.SavedSearch) map[string]interface{} {
	if ss == nil {
		return nil
	}

	messages := make([]interface{}, len(ss.Messages))
	for i, msg := range ss.Messages {
		messages[i] = map[string]interface{}{
			"id":            msg.ID,
			"role":          msg.Role,
			"content":       msg.Content,
			"timestamp":     msg.Timestamp,
			"quick_replies": msg.QuickReplies,
			"products":      msg.Products,
			"search_type":   msg.SearchType,
		}
	}

	return map[string]interface{}{
		"session_id": ss.SessionID,
		"category":   ss.Category,
		"timestamp":  ss.Timestamp,
		"messages":   messages,
	}
}

// Convert map[string]interface{} to SavedSearch
func mapToSavedSearch(m map[string]interface{}) *models.SavedSearch {
	if m == nil {
		return nil
	}

	ss := &models.SavedSearch{}

	if v, ok := m["session_id"].(string); ok {
		ss.SessionID = v
	}
	if v, ok := m["category"].(string); ok {
		ss.Category = v
	}
	if v, ok := m["timestamp"].(float64); ok {
		ss.Timestamp = int64(v)
	}

	if messagesRaw, ok := m["messages"].([]interface{}); ok {
		ss.Messages = make([]models.SavedMessage, len(messagesRaw))
		for i, msgRaw := range messagesRaw {
			if msgMap, ok := msgRaw.(map[string]interface{}); ok {
				msg := models.SavedMessage{}
				if v, ok := msgMap["id"].(string); ok {
					msg.ID = v
				}
				if v, ok := msgMap["role"].(string); ok {
					msg.Role = v
				}
				if v, ok := msgMap["content"].(string); ok {
					msg.Content = v
				}
				if v, ok := msgMap["timestamp"].(float64); ok {
					msg.Timestamp = int64(v)
				}
				if v, ok := msgMap["search_type"].(string); ok {
					msg.SearchType = v
				}
				// TODO: Handle quick_replies and products if needed
				ss.Messages[i] = msg
			}
		}
	}

	return ss
}

// isForeignKeyError checks if an error is a PostgreSQL foreign key constraint violation
func isForeignKeyError(err error) bool {
	if err == nil {
		return false
	}
	errMsg := err.Error()
	return contains(errMsg, "violates foreign key constraint") || contains(errMsg, "23503")
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && stringContains(s, substr)))
}

func stringContains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// ensureUserExistsInPostgres fetches user from Redis and saves to PostgreSQL if missing
func (s *PreferencesService) ensureUserExistsInPostgres(userID uuid.UUID) error {
	// Get user from Redis (via AuthService)
	user, err := s.authService.GetUserByID(userID)
	if err != nil {
		return fmt.Errorf("failed to fetch user from Redis: %w", err)
	}

	// Save user to PostgreSQL
	if err := s.authService.SaveUserToPostgres(user); err != nil {
		return fmt.Errorf("failed to save user to PostgreSQL: %w", err)
	}

	fmt.Printf("✅ Synced user %s to PostgreSQL\n", userID.String())
	return nil
}
