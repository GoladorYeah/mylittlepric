package models

import (
	"time"

	"github.com/google/uuid"
)

// ═══════════════════════════════════════════════════════════
// SEARCH HISTORY MODELS
// ═══════════════════════════════════════════════════════════

type SearchHistory struct {
	ID               uuid.UUID     `json:"id" db:"id"`
	UserID           *uuid.UUID    `json:"user_id,omitempty" db:"user_id"`
	SessionID        *string       `json:"session_id,omitempty" db:"session_id"`
	SearchQuery      string        `json:"search_query" db:"search_query"`
	OptimizedQuery   *string       `json:"optimized_query,omitempty" db:"optimized_query"`
	SearchType       string        `json:"search_type" db:"search_type"`
	Category         *string       `json:"category,omitempty" db:"category"`
	CountryCode      string        `json:"country_code" db:"country_code"`
	LanguageCode     string        `json:"language_code" db:"language_code"`
	Currency         string        `json:"currency" db:"currency"`
	ResultCount      int           `json:"result_count" db:"result_count"`
	ProductsFound    []ProductCard `json:"products_found,omitempty" db:"products_found"`
	ClickedProductID *string       `json:"clicked_product_id,omitempty" db:"clicked_product_id"`
	CreatedAt        time.Time     `json:"created_at" db:"created_at"`
	ExpiresAt        *time.Time    `json:"expires_at,omitempty" db:"expires_at"`
}

type SearchHistoryListRequest struct {
	UserID *uuid.UUID `json:"user_id,omitempty"`
	Limit  int        `json:"limit,omitempty"`
	Offset int        `json:"offset,omitempty"`
}

type SearchHistoryListResponse struct {
	Items   []SearchHistory `json:"items"`
	Total   int             `json:"total"`
	Limit   int             `json:"limit"`
	Offset  int             `json:"offset"`
	HasMore bool            `json:"has_more"`
}

type SearchHistoryDeleteRequest struct {
	ID     uuid.UUID  `json:"id" validate:"required"`
	UserID *uuid.UUID `json:"user_id,omitempty"`
}
