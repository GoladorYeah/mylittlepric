package models

import (
	"time"

	"github.com/google/uuid"
)

// ═══════════════════════════════════════════════════════════
// SESSION MODELS
// ═══════════════════════════════════════════════════════════

type ChatSession struct {
	ID                  uuid.UUID            `json:"id" db:"id"`
	SessionID           string               `json:"session_id" db:"session_id"`
	CountryCode         string               `json:"country_code" db:"country_code"`
	LanguageCode        string               `json:"language_code" db:"language_code"`
	Currency            string               `json:"currency" db:"currency"`
	MessageCount        int                  `json:"message_count" db:"message_count"`
	SearchState         SearchState          `json:"search_state" db:"search_state"`
	CycleState          CycleState           `json:"cycle_state" db:"cycle_state"`
	ConversationContext *ConversationContext `json:"conversation_context,omitempty" db:"conversation_context"`
	CreatedAt           time.Time            `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time            `json:"updated_at" db:"updated_at"`
	ExpiresAt           time.Time            `json:"expires_at" db:"expires_at"`
}

type SearchState struct {
	Status         SearchStatus `json:"status"`
	Category       string       `json:"category"`
	LastSearchTime time.Time    `json:"last_search_time,omitempty"`
	SearchCount    int          `json:"search_count"`
	LastProduct    *ProductInfo `json:"last_product,omitempty"`
}

type SearchStatus string

const (
	SearchStatusIdle       SearchStatus = "idle"
	SearchStatusInProgress SearchStatus = "in_progress"
	SearchStatusCompleted  SearchStatus = "completed"
)

// CycleState tracks the Universal Prompt Cycle system state
type CycleState struct {
	CycleID          int               `json:"cycle_id"`           // Current cycle number
	Iteration        int               `json:"iteration"`          // Current iteration within cycle (1-6)
	CycleHistory     []CycleMessage    `json:"cycle_history"`      // Messages in current cycle
	LastCycleContext *LastCycleContext `json:"last_cycle_context"` // Context from previous cycle
	LastDefined      []string          `json:"last_defined"`       // Last confirmed product names
	PromptID         string            `json:"prompt_id"`          // Prompt version identifier
	PromptHash       string            `json:"prompt_hash"`        // SHA-256 hash for drift detection
}

// CycleMessage represents a single message in a cycle
type CycleMessage struct {
	Role      string    `json:"role"`      // "user" or "assistant"
	Content   string    `json:"content"`   // Message content
	Timestamp time.Time `json:"timestamp"` // When message was sent
}

// LastCycleContext contains context from the previous cycle
type LastCycleContext struct {
	Groups      []string      `json:"groups"`       // Product groups discussed
	Subgroups   []string      `json:"subgroups"`    // Product subgroups discussed
	Products    []ProductInfo `json:"products"`     // Products identified
	LastRequest string        `json:"last_request"` // The final user request from last cycle
}

// ═══════════════════════════════════════════════════════════
// CONTEXT MANAGEMENT
// ═══════════════════════════════════════════════════════════

// ConversationContext stores optimized conversation context
type ConversationContext struct {
	Summary     string          `json:"summary"`               // AI-generated compact summary
	Preferences UserPreferences `json:"preferences"`           // Structured user preferences
	LastSearch  *SearchContext  `json:"last_search,omitempty"` // Most recent search context
	Exclusions  []string        `json:"exclusions,omitempty"`  // User exclusions
	UpdatedAt   time.Time       `json:"updated_at"`            // Last update timestamp
}

// UserPreferences stores structured user preferences
type UserPreferences struct {
	PriceRange   *PriceRange `json:"price_range,omitempty"`  // Price range preference
	Brands       []string    `json:"brands,omitempty"`       // Preferred brands
	Features     []string    `json:"features,omitempty"`     // Required features
	Requirements []string    `json:"requirements,omitempty"` // Special requirements
}

// PriceRange represents a price range with currency
type PriceRange struct {
	Min      *float64 `json:"min,omitempty"` // Minimum price
	Max      *float64 `json:"max,omitempty"` // Maximum price
	Currency string   `json:"currency"`      // Currency code
}

// SearchContext stores the most recent search
type SearchContext struct {
	Query         string        `json:"query"`                    // Search query used
	Category      string        `json:"category"`                 // Category searched
	ProductsShown []ProductInfo `json:"products_shown,omitempty"` // Products shown
	UserFeedback  string        `json:"user_feedback,omitempty"`  // User feedback
	Timestamp     time.Time     `json:"timestamp"`                // When search occurred
}
