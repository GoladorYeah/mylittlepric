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
	CycleState   CycleState  `json:"cycle_state" db:"cycle_state"` // NEW: Cycle tracking
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

// CycleState tracks the Universal Prompt Cycle system state
type CycleState struct {
	CycleID          int                      `json:"cycle_id"`            // Current cycle number (increments on new cycle)
	Iteration        int                      `json:"iteration"`           // Current iteration within cycle (1-6)
	CycleHistory     []CycleMessage           `json:"cycle_history"`       // Messages in current cycle
	LastCycleContext *LastCycleContext        `json:"last_cycle_context"`  // Context from previous cycle
	LastDefined      []string                 `json:"last_defined"`        // Last confirmed product names or shortlist
	PromptID         string                   `json:"prompt_id"`           // Prompt version identifier
	PromptHash       string                   `json:"prompt_hash"`         // SHA-256 hash of prompt for drift detection
}

// CycleMessage represents a single message in a cycle
type CycleMessage struct {
	Role      string    `json:"role"`      // "user" or "assistant"
	Content   string    `json:"content"`   // Message content
	Timestamp time.Time `json:"timestamp"` // When message was sent
}

// LastCycleContext contains context from the previous cycle to carry over
type LastCycleContext struct {
	Groups      []string      `json:"groups"`       // Product groups discussed
	Subgroups   []string      `json:"subgroups"`    // Product subgroups discussed
	Products    []ProductInfo `json:"products"`     // Products identified
	LastRequest string        `json:"last_request"` // The final user request from last cycle
}

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
	ResponseType  string                 `json:"response_type"` // "dialogue", "search", or "api_request"
	Output        string                 `json:"output"`
	QuickReplies  []string               `json:"quick_replies"`
	SearchPhrase  string                 `json:"search_phrase"` // For response_type="search"
	SearchType    string                 `json:"search_type"`   // "exact", "parameters", or "category"
	Category      string                 `json:"category"`
	PriceFilter   string                 `json:"price_filter,omitempty"`   // "cheaper" or "expensive"
	MinPrice      *float64               `json:"min_price,omitempty"`      // Minimum price in user's currency
	MaxPrice      *float64               `json:"max_price,omitempty"`      // Maximum price in user's currency
	ProductType   string                 `json:"product_type"`
	Brand         string                 `json:"brand"`
	Confidence    float32                `json:"confidence"`
	RequiresInput bool                   `json:"requires_input"`
	// New fields for api_request response type
	API    string                 `json:"api,omitempty"`    // API name (e.g., "google_shopping")
	Params map[string]interface{} `json:"params,omitempty"` // API parameters
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
	Merchant             string   `json:"merchant"`
	Logo                 string   `json:"logo,omitempty"`
	Price                string   `json:"price"`
	ExtractedPrice       float64  `json:"extracted_price,omitempty"`
	Currency             string   `json:"currency,omitempty"`
	Link                 string   `json:"link"`
	Title                string   `json:"title,omitempty"`
	Availability         string   `json:"availability,omitempty"`
	Shipping             string   `json:"shipping,omitempty"`
	ShippingExtracted    float64  `json:"shipping_extracted,omitempty"`
	Total                string   `json:"total,omitempty"`
	ExtractedTotal       float64  `json:"extracted_total,omitempty"`
	Rating               float32  `json:"rating,omitempty"`
	Reviews              int      `json:"reviews,omitempty"`
	PaymentMethods       string   `json:"payment_methods,omitempty"`
	Tag                  string   `json:"tag,omitempty"`
	DetailsAndOffers     []string `json:"details_and_offers,omitempty"`
	MonthlyPaymentDur    int      `json:"monthly_payment_duration,omitempty"`
	DownPayment          string   `json:"down_payment,omitempty"`
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
	ID           uuid.UUID `json:"id" db:"id"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"` // Never expose password (optional for OAuth)
	FullName     string    `json:"full_name,omitempty" db:"full_name"`
	Picture      string    `json:"picture,omitempty" db:"picture"`        // Profile picture URL
	Provider     string    `json:"provider" db:"provider"`                // "google", "email"
	ProviderID   string    `json:"provider_id,omitempty" db:"provider_id"` // Google user ID
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty" db:"last_login_at"`
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
	Picture   string    `json:"picture,omitempty"`
	Provider  string    `json:"provider"`
	CreatedAt time.Time `json:"created_at"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type ClaimSessionsRequest struct {
	SessionIDs []string `json:"session_ids" validate:"required"`
}

// Google OAuth Request/Response Models

type GoogleAuthRequest struct {
	IDToken string `json:"id_token" validate:"required"`
}

type GoogleUserInfo struct {
	Sub           string `json:"sub"`            // Google user ID
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Locale        string `json:"locale"`
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
