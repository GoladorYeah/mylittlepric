package services

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"mylittleprice/internal/models"
)

var (
	ErrPreferencesNotFound = errors.New("user preferences not found")
)

type PreferencesService struct {
	db *sqlx.DB
}

func NewPreferencesService(db *sqlx.DB) *PreferencesService {
	return &PreferencesService{
		db: db,
	}
}

// GetUserPreferences retrieves user preferences from PostgreSQL
// Returns nil if preferences not found (first time user)
func (s *PreferencesService) GetUserPreferences(userID uuid.UUID) (*models.UserPreferences, error) {
	var prefs models.UserPreferences

	query := `
		SELECT user_id, country, currency, language, theme, sidebar_open, created_at, updated_at
		FROM user_preferences
		WHERE user_id = $1
	`

	err := s.db.Get(&prefs, query, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Not found is OK - return nil (user hasn't set preferences yet)
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user preferences: %w", err)
	}

	return &prefs, nil
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
	query := `
		INSERT INTO user_preferences (user_id, country, currency, language, theme, sidebar_open)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING user_id, country, currency, language, theme, sidebar_open, created_at, updated_at
	`

	var prefs models.UserPreferences
	err := s.db.QueryRowx(query,
		userID,
		update.Country,
		update.Currency,
		update.Language,
		update.Theme,
		update.SidebarOpen,
	).StructScan(&prefs)

	if err != nil {
		return nil, fmt.Errorf("failed to create user preferences: %w", err)
	}

	fmt.Printf("✅ Created preferences for user %s\n", userID.String())
	return &prefs, nil
}

// updatePreferences updates existing preferences (only provided fields)
func (s *PreferencesService) updatePreferences(userID uuid.UUID, update *models.UserPreferencesUpdate) (*models.UserPreferences, error) {
	// Build dynamic update query (only update provided fields)
	query := `UPDATE user_preferences SET `
	args := []interface{}{}
	argCounter := 1
	updates := []string{}

	if update.Country != nil {
		updates = append(updates, fmt.Sprintf("country = $%d", argCounter))
		args = append(args, update.Country)
		argCounter++
	}

	if update.Currency != nil {
		updates = append(updates, fmt.Sprintf("currency = $%d", argCounter))
		args = append(args, update.Currency)
		argCounter++
	}

	if update.Language != nil {
		updates = append(updates, fmt.Sprintf("language = $%d", argCounter))
		args = append(args, update.Language)
		argCounter++
	}

	if update.Theme != nil {
		updates = append(updates, fmt.Sprintf("theme = $%d", argCounter))
		args = append(args, update.Theme)
		argCounter++
	}

	if update.SidebarOpen != nil {
		updates = append(updates, fmt.Sprintf("sidebar_open = $%d", argCounter))
		args = append(args, update.SidebarOpen)
		argCounter++
	}

	// No fields to update
	if len(updates) == 0 {
		return s.GetUserPreferences(userID)
	}

	// Finish query
	for i, update := range updates {
		if i > 0 {
			query += ", "
		}
		query += update
	}
	query += fmt.Sprintf(" WHERE user_id = $%d RETURNING user_id, country, currency, language, theme, sidebar_open, created_at, updated_at", argCounter)
	args = append(args, userID)

	// Execute update
	var prefs models.UserPreferences
	err := s.db.QueryRowx(query, args...).StructScan(&prefs)
	if err != nil {
		return nil, fmt.Errorf("failed to update user preferences: %w", err)
	}

	fmt.Printf("✅ Updated preferences for user %s\n", userID.String())
	return &prefs, nil
}

// DeleteUserPreferences removes user preferences (e.g., on account deletion)
func (s *PreferencesService) DeleteUserPreferences(userID uuid.UUID) error {
	query := `DELETE FROM user_preferences WHERE user_id = $1`

	result, err := s.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user preferences: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrPreferencesNotFound
	}

	fmt.Printf("✅ Deleted preferences for user %s\n", userID.String())
	return nil
}
