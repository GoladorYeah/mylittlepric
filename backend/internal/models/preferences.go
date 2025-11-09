package models

import (
	"time"

	"github.com/google/uuid"
)

// UserPreferences represents user settings that should sync across devices
type UserPreferences struct {
	UserID uuid.UUID `json:"user_id" db:"user_id"`

	// Regional settings
	Country  *string `json:"country,omitempty" db:"country"`   // ISO 3166-1 alpha-2 (e.g., "US", "GB")
	Currency *string `json:"currency,omitempty" db:"currency"` // ISO 4217 (e.g., "USD", "EUR")
	Language *string `json:"language,omitempty" db:"language"` // ISO 639-1 (e.g., "en", "es")

	// UI settings
	Theme       *string `json:"theme,omitempty" db:"theme"`             // "light", "dark", "system"
	SidebarOpen *bool   `json:"sidebar_open,omitempty" db:"sidebar_open"` // UI state

	// Timestamps
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// UserPreferencesUpdate represents fields that can be updated
// All fields are pointers to distinguish between "not set" and "set to empty/default"
type UserPreferencesUpdate struct {
	Country     *string `json:"country,omitempty"`
	Currency    *string `json:"currency,omitempty"`
	Language    *string `json:"language,omitempty"`
	Theme       *string `json:"theme,omitempty"`
	SidebarOpen *bool   `json:"sidebar_open,omitempty"`
}
