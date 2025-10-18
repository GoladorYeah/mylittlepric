package models

import (
	"time"

	"github.com/google/uuid"
)

type ChatSession struct {
	ID           uuid.UUID   `json:"id" db:"id"`
	SessionID    string      `json:"session_id" db:"session_id"`
	CountryCode  string      `json:"country_code" db:"country_code"`
	LanguageCode string      `json:"language_code" db:"language_code"`
	Currency     string      `json:"currency" db:"currency"`
	MessageCount int         `json:"message_count" db:"message_count"`
	SearchState  SearchState `json:"search_state" db:"search_state"`
	CreatedAt    time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at" db:"updated_at"`
	ExpiresAt    time.Time   `json:"expires_at" db:"expires_at"`
}

type SearchState struct {
	Status         SearchStatus `json:"status"`
	Category       string       `json:"category"`
	LastSearchTime time.Time    `json:"last_search_time,omitempty"`
	SearchCount    int          `json:"search_count"`
	LastProduct    *ProductInfo `json:"last_product,omitempty"`
}

type ProductInfo struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type SearchStatus string

const (
	SearchStatusIdle       SearchStatus = "idle"
	SearchStatusInProgress SearchStatus = "in_progress"
	SearchStatusCompleted  SearchStatus = "completed"
)

type Message struct {
	ID           uuid.UUID              `json:"id" db:"id"`
	SessionID    uuid.UUID              `json:"session_id" db:"session_id"`
	Role         string                 `json:"role" db:"role"`
	Content      string                 `json:"content" db:"content"`
	ResponseType string                 `json:"response_type,omitempty" db:"response_type"`
	QuickReplies []string               `json:"quick_replies,omitempty" db:"quick_replies"`
	Products     []ProductCard          `json:"products,omitempty" db:"products"`
	SearchInfo   map[string]interface{} `json:"search_info,omitempty" db:"search_info"`
	CreatedAt    time.Time              `json:"created_at" db:"created_at"`
}

type ProductCard struct {
	Name        string `json:"name"`
	Price       string `json:"price"`
	OldPrice    string `json:"old_price,omitempty"`
	Link        string `json:"link"`
	Image       string `json:"image"`
	Description string `json:"description,omitempty"`
	Badge       string `json:"badge,omitempty"`
	PageToken   string `json:"page_token"`
}

type ChatRequest struct {
	SessionID       string `json:"session_id"`
	Message         string `json:"message"`
	Country         string `json:"country"`
	Language        string `json:"language"`
	Currency        string `json:"currency"`
	NewSearch       bool   `json:"new_search"`
	CurrentCategory string `json:"current_category"`
}

type ChatResponse struct {
	Type         string               `json:"type"`
	Output       string               `json:"output,omitempty"`
	QuickReplies []string             `json:"quick_replies,omitempty"`
	Products     []ProductCard        `json:"products,omitempty"`
	SearchType   string               `json:"search_type,omitempty"`
	SessionID    string               `json:"session_id"`
	MessageCount int                  `json:"message_count"`
	SearchState  *SearchStateResponse `json:"search_state,omitempty"`
}

type SearchStateResponse struct {
	Status      string `json:"status"`
	Category    string `json:"category,omitempty"`
	CanContinue bool   `json:"can_continue"`
	SearchCount int    `json:"search_count"`
	MaxSearches int    `json:"max_searches"`
	Message     string `json:"message,omitempty"`
}

type GeminiResponse struct {
	ResponseType  string   `json:"response_type"`
	Output        string   `json:"output"`
	QuickReplies  []string `json:"quick_replies"`
	SearchPhrase  string   `json:"search_phrase"`
	SearchType    string   `json:"search_type"`
	Category      string   `json:"category"`
	ProductType   string   `json:"product_type"`
	Brand         string   `json:"brand"`
	Confidence    float32  `json:"confidence"`
	RequiresInput bool     `json:"requires_input"`
}

type SerpConfig struct {
	Type       string `json:"type"`
	Query      string `json:"query"`
	Country    string `json:"country"`
	Language   string `json:"language"`
	MaxResults int    `json:"max_results"`
}

type ProductDetailsRequest struct {
	PageToken string `json:"page_token"`
	Country   string `json:"country"`
}

type ProductDetailsResponse struct {
	Type            string                `json:"type"`
	Title           string                `json:"title"`
	Price           string                `json:"price"`
	Rating          float32               `json:"rating,omitempty"`
	Reviews         int                   `json:"reviews,omitempty"`
	Description     string                `json:"description,omitempty"`
	Images          []string              `json:"images,omitempty"`
	Specifications  []Specification       `json:"specifications,omitempty"`
	Variants        []Variant             `json:"variants,omitempty"`
	Offers          []Offer               `json:"offers"`
	Videos          []interface{}         `json:"videos,omitempty"`
	MoreOptions     []interface{}         `json:"more_options,omitempty"`
	RatingBreakdown []RatingBreakdownItem `json:"rating_breakdown,omitempty"`
}

type Specification struct {
	Title string `json:"title"`
	Value string `json:"value"`
}

type Variant struct {
	Title string        `json:"title"`
	Items []interface{} `json:"items"`
}

type Offer struct {
	Merchant     string  `json:"merchant"`
	Price        string  `json:"price"`
	Currency     string  `json:"currency"`
	Link         string  `json:"link"`
	Availability string  `json:"availability,omitempty"`
	Shipping     string  `json:"shipping,omitempty"`
	Rating       float32 `json:"rating,omitempty"`
}

type RatingBreakdownItem struct {
	Stars  int `json:"stars"`
	Amount int `json:"amount"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// ==================== Authentication Models ====================

type User struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Email       string    `json:"email" db:"email"`
	PasswordHash string   `json:"-" db:"password_hash"` // Never expose password
	FullName    string    `json:"full_name,omitempty" db:"full_name"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty" db:"last_login_at"`
}

type RefreshToken struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	UserID    uuid.UUID  `json:"user_id" db:"user_id"`
	TokenHash string     `json:"-" db:"token_hash"`
	ExpiresAt time.Time  `json:"expires_at" db:"expires_at"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	RevokedAt *time.Time `json:"revoked_at,omitempty" db:"revoked_at"`
}

// Auth Request/Response Models

type SignupRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	FullName string `json:"full_name,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	User         *UserInfo `json:"user"`
	ExpiresIn    int64     `json:"expires_in"` // seconds
}

type UserInfo struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type ClaimSessionsRequest struct {
	SessionIDs []string `json:"session_ids" validate:"required"`
}

// ==================== Search History Models ====================

type SearchHistory struct {
	ID               uuid.UUID       `json:"id" db:"id"`
	UserID           *uuid.UUID      `json:"user_id,omitempty" db:"user_id"`
	SessionID        *string         `json:"session_id,omitempty" db:"session_id"`
	SearchQuery      string          `json:"search_query" db:"search_query"`
	OptimizedQuery   *string         `json:"optimized_query,omitempty" db:"optimized_query"`
	SearchType       string          `json:"search_type" db:"search_type"`
	Category         *string         `json:"category,omitempty" db:"category"`
	CountryCode      string          `json:"country_code" db:"country_code"`
	LanguageCode     string          `json:"language_code" db:"language_code"`
	Currency         string          `json:"currency" db:"currency"`
	ResultCount      int             `json:"result_count" db:"result_count"`
	ProductsFound    []ProductCard   `json:"products_found,omitempty" db:"products_found"`
	ClickedProductID *string         `json:"clicked_product_id,omitempty" db:"clicked_product_id"`
	CreatedAt        time.Time       `json:"created_at" db:"created_at"`
	ExpiresAt        *time.Time      `json:"expires_at,omitempty" db:"expires_at"`
}

type SearchHistoryListRequest struct {
	UserID *uuid.UUID `json:"user_id,omitempty"`
	Limit  int        `json:"limit,omitempty"`
	Offset int        `json:"offset,omitempty"`
}

type SearchHistoryListResponse struct {
	Items      []SearchHistory `json:"items"`
	Total      int             `json:"total"`
	Limit      int             `json:"limit"`
	Offset     int             `json:"offset"`
	HasMore    bool            `json:"has_more"`
}

type SearchHistoryDeleteRequest struct {
	ID     uuid.UUID  `json:"id" validate:"required"`
	UserID *uuid.UUID `json:"user_id,omitempty"`
}
