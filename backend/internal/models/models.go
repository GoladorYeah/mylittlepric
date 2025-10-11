package models

import (
	"time"

	"github.com/google/uuid"
)

// ═══════════════════════════════════════════════════════════════
// SESSION MODELS (WITH SEARCH STATE)
// ═══════════════════════════════════════════════════════════════

// ChatSession represents a user chat session with search state management
type ChatSession struct {
	ID           uuid.UUID   `json:"id" db:"id"`
	SessionID    string      `json:"session_id" db:"session_id"`
	CountryCode  string      `json:"country_code" db:"country_code"`
	LanguageCode string      `json:"language_code" db:"language_code"`
	Currency     string      `json:"currency" db:"currency"`
	MessageCount int         `json:"message_count" db:"message_count"`
	SearchState  SearchState `json:"search_state" db:"search_state"` // 🆕 NEW
	CreatedAt    time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at" db:"updated_at"`
	ExpiresAt    time.Time   `json:"expires_at" db:"expires_at"`
}

// SearchState tracks the current search context within a session
type SearchState struct {
	Status          SearchStatus `json:"status"`                     // Current search status
	Category        string       `json:"category"`                   // Product category (electronics, clothing, etc.) - 🎯 SAVED!
	ProductType     string       `json:"product_type,omitempty"`     // Specific product type (phone, laptop, etc.)
	Brand           string       `json:"brand,omitempty"`            // Brand name (Apple, Samsung, etc.)
	CollectedParams []string     `json:"collected_params,omitempty"` // User-provided parameters
	LastSearchTime  time.Time    `json:"last_search_time,omitempty"` // When last search completed
	SearchCount     int          `json:"search_count"`               // Number of searches in this session
}

// SearchStatus represents the current state of a product search
type SearchStatus string

const (
	SearchStatusIdle       SearchStatus = "idle"        // No active search
	SearchStatusInProgress SearchStatus = "in_progress" // Collecting information
	SearchStatusCompleted  SearchStatus = "completed"   // Results shown, chat blocked
	SearchStatusBlocked    SearchStatus = "blocked"     // Max searches reached
)

// ═══════════════════════════════════════════════════════════════
// MESSAGE MODELS
// ═══════════════════════════════════════════════════════════════

// Message represents a single chat message
type Message struct {
	ID           uuid.UUID              `json:"id" db:"id"`
	SessionID    uuid.UUID              `json:"session_id" db:"session_id"`
	Role         string                 `json:"role" db:"role"` // "user" or "assistant"
	Content      string                 `json:"content" db:"content"`
	ResponseType string                 `json:"response_type,omitempty" db:"response_type"` // "dialogue", "search"
	QuickReplies []string               `json:"quick_replies,omitempty" db:"quick_replies"`
	Category     string                 `json:"category,omitempty" db:"category"` // 🆕 NEW - preserved category context
	Metadata     map[string]interface{} `json:"metadata,omitempty" db:"metadata"`
	CreatedAt    time.Time              `json:"created_at" db:"created_at"`
}

// ═══════════════════════════════════════════════════════════════
// REQUEST/RESPONSE MODELS
// ═══════════════════════════════════════════════════════════════

// ChatRequest represents incoming chat request from frontend
type ChatRequest struct {
	SessionID string `json:"session_id"` // Optional - will be created if empty
	Message   string `json:"message"`    // User message
	Country   string `json:"country"`    // Country code (CH, DE, etc.)
	Language  string `json:"language"`   // Language code (de, en, etc.)
	NewSearch bool   `json:"new_search"` // 🆕 NEW - Set to true to start new search
}

// ChatResponse represents chat response to frontend
type ChatResponse struct {
	Type         string           `json:"type"` // "text", "product_card", "search_blocked"
	Output       string           `json:"output,omitempty"`
	QuickReplies []string         `json:"quick_replies,omitempty"`
	Products     []ProductCard    `json:"products,omitempty"`
	SearchType   string           `json:"search_type,omitempty"`
	SessionID    string           `json:"session_id"`
	MessageCount int              `json:"message_count"`
	SearchState  *SearchStateInfo `json:"search_state,omitempty"` // 🆕 NEW - Search state info for frontend
}

// SearchStateInfo provides frontend information about search state
type SearchStateInfo struct {
	Status      SearchStatus `json:"status"`             // Current status
	Category    string       `json:"category,omitempty"` // Current category
	CanContinue bool         `json:"can_continue"`       // Can start new search?
	SearchCount int          `json:"search_count"`       // Current search count
	MaxSearches int          `json:"max_searches"`       // Maximum allowed searches
	Message     string       `json:"message,omitempty"`  // Hint message for user
}

// ═══════════════════════════════════════════════════════════════
// AI RESPONSE MODELS
// ═══════════════════════════════════════════════════════════════

// GeminiResponse represents AI response structure from Gemini
type GeminiResponse struct {
	ResponseType string   `json:"response_type"`           // "dialogue" or "search"
	Output       string   `json:"output,omitempty"`        // Text output
	QuickReplies []string `json:"quick_replies,omitempty"` // Quick reply options
	SearchPhrase string   `json:"search_phrase,omitempty"` // Search query
	SearchType   string   `json:"search_type,omitempty"`   // "exact", "parameters", "category"
	Category     string   `json:"category,omitempty"`      // 🆕 NEW - Product category (set once)
}

// ═══════════════════════════════════════════════════════════════
// PRODUCT MODELS
// ═══════════════════════════════════════════════════════════════

// ProductCard represents a product card for frontend display
type ProductCard struct {
	Name        string `json:"name"`
	Price       string `json:"price"`
	OldPrice    string `json:"old_price,omitempty"`
	Link        string `json:"link"`
	Image       string `json:"image"`
	Description string `json:"description,omitempty"`
	Badge       string `json:"badge,omitempty"`
	PageToken   string `json:"page_token"` // For Google Immersive Product API
}

// Product represents a cached product in database
type Product struct {
	ID             uuid.UUID              `json:"id" db:"id"`
	PageToken      string                 `json:"page_token" db:"page_token"` // Google Immersive Product token
	ProductName    string                 `json:"product_name" db:"product_name"`
	CountryCode    string                 `json:"country_code" db:"country_code"`
	PriceInfo      map[string]interface{} `json:"price_info" db:"price_info"`
	MerchantInfo   map[string]interface{} `json:"merchant_info,omitempty" db:"merchant_info"`
	ProductDetails map[string]interface{} `json:"product_details,omitempty" db:"product_details"`
	ImageURL       string                 `json:"image_url,omitempty" db:"image_url"`
	Link           string                 `json:"link,omitempty" db:"link"`
	Rating         float32                `json:"rating,omitempty" db:"rating"`
	FetchCount     int                    `json:"fetch_count" db:"fetch_count"`
	LastFetchedAt  time.Time              `json:"last_fetched_at" db:"last_fetched_at"`
	CreatedAt      time.Time              `json:"created_at" db:"created_at"`
}

// ═══════════════════════════════════════════════════════════════
// PRODUCT DETAILS MODELS
// ═══════════════════════════════════════════════════════════════

// ProductDetailsRequest represents product details request
type ProductDetailsRequest struct {
	PageToken string `json:"page_token"` // Google Immersive Product page_token
	Country   string `json:"country"`
}

// ProductDetailsResponse represents detailed product information
type ProductDetailsResponse struct {
	Type            string               `json:"type"` // "product_details"
	Title           string               `json:"title"`
	Price           string               `json:"price"`
	Rating          float32              `json:"rating,omitempty"`
	Reviews         int                  `json:"reviews,omitempty"`
	Description     string               `json:"description,omitempty"`
	Images          []string             `json:"images,omitempty"`
	Specifications  []ProductSpec        `json:"specifications,omitempty"`
	Variants        []ProductVariant     `json:"variants,omitempty"`
	Offers          []ProductOffer       `json:"offers"`
	Videos          []ProductVideo       `json:"videos,omitempty"`
	MoreOptions     []AlternativeProduct `json:"more_options,omitempty"`
	RatingBreakdown []RatingBreakdown    `json:"rating_breakdown,omitempty"`
}

// ProductSpec represents product specification
type ProductSpec struct {
	Title string `json:"title"`
	Value string `json:"value"`
}

// ProductVariant represents variant option (color, storage, etc.)
type ProductVariant struct {
	Title string        `json:"title"` // "Storage Capacity", "Colour"
	Items []VariantItem `json:"items"`
}

// VariantItem represents a single variant option
type VariantItem struct {
	Name      string `json:"name"`
	Selected  bool   `json:"selected"`
	Available bool   `json:"available"`
	PageToken string `json:"page_token"` // Token to load this variant
}

// ProductOffer represents a single offer from a merchant
type ProductOffer struct {
	Merchant     string  `json:"merchant"`
	Price        string  `json:"price"`
	Currency     string  `json:"currency"`
	Link         string  `json:"link"`
	Availability string  `json:"availability,omitempty"`
	Shipping     string  `json:"shipping,omitempty"`
	Rating       float32 `json:"rating,omitempty"`
}

// ProductVideo represents product video
type ProductVideo struct {
	Title     string `json:"title"`
	Link      string `json:"link"`
	Source    string `json:"source"`
	Channel   string `json:"channel,omitempty"`
	Duration  string `json:"duration,omitempty"`
	Thumbnail string `json:"thumbnail"`
}

// AlternativeProduct represents alternative/similar product
type AlternativeProduct struct {
	Title     string  `json:"title"`
	Thumbnail string  `json:"thumbnail"`
	Price     string  `json:"price"`
	Rating    float32 `json:"rating,omitempty"`
	Reviews   int     `json:"reviews,omitempty"`
	PageToken string  `json:"page_token"`
}

// RatingBreakdown represents rating distribution
type RatingBreakdown struct {
	Stars  int `json:"stars"`
	Amount int `json:"amount"`
}

// ═══════════════════════════════════════════════════════════════
// ANALYTICS MODELS
// ═══════════════════════════════════════════════════════════════

// SearchQuery represents a product search query for analytics
type SearchQuery struct {
	ID             uuid.UUID `json:"id" db:"id"`
	SessionID      uuid.UUID `json:"session_id,omitempty" db:"session_id"`
	OriginalQuery  string    `json:"original_query" db:"original_query"`
	OptimizedQuery string    `json:"optimized_query" db:"optimized_query"`
	SearchType     string    `json:"search_type" db:"search_type"` // "exact", "parameters", "category"
	CountryCode    string    `json:"country_code" db:"country_code"`
	ResultCount    int       `json:"result_count" db:"result_count"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

// APIUsage represents API usage tracking for monitoring
type APIUsage struct {
	ID             uuid.UUID `json:"id" db:"id"`
	APIName        string    `json:"api_name" db:"api_name"` // "gemini" or "serp"
	KeyIndex       int       `json:"key_index" db:"key_index"`
	RequestType    string    `json:"request_type,omitempty" db:"request_type"`
	ResponseTimeMs int       `json:"response_time_ms,omitempty" db:"response_time_ms"`
	Success        bool      `json:"success" db:"success"`
	ErrorMessage   string    `json:"error_message,omitempty" db:"error_message"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

// ═══════════════════════════════════════════════════════════════
// HEALTH & ERROR MODELS
// ═══════════════════════════════════════════════════════════════

// HealthResponse represents health check response
type HealthResponse struct {
	Status    string                 `json:"status"`
	Timestamp time.Time              `json:"timestamp"`
	Services  map[string]ServiceInfo `json:"services"`
}

// ServiceInfo represents status of a service
type ServiceInfo struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

// ErrorResponse represents error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    string `json:"code,omitempty"`
}

// ═══════════════════════════════════════════════════════════════
// HELPER METHODS
// ═══════════════════════════════════════════════════════════════

// IsIdle checks if search is idle
func (s SearchStatus) IsIdle() bool {
	return s == SearchStatusIdle
}

// IsInProgress checks if search is in progress
func (s SearchStatus) IsInProgress() bool {
	return s == SearchStatusInProgress
}

// IsCompleted checks if search is completed
func (s SearchStatus) IsCompleted() bool {
	return s == SearchStatusCompleted
}

// IsBlocked checks if search is blocked
func (s SearchStatus) IsBlocked() bool {
	return s == SearchStatusBlocked
}

// String returns string representation of SearchStatus
func (s SearchStatus) String() string {
	return string(s)
}

// HasCategory checks if search state has a category set
func (ss *SearchState) HasCategory() bool {
	return ss.Category != ""
}

// HasBrand checks if search state has a brand set
func (ss *SearchState) HasBrand() bool {
	return ss.Brand != ""
}

// HasProductType checks if search state has a product type set
func (ss *SearchState) HasProductType() bool {
	return ss.ProductType != ""
}

// IsSearchBlocked checks if chat response indicates blocked search
func (cr *ChatResponse) IsSearchBlocked() bool {
	return cr.Type == "search_blocked"
}

// HasProducts checks if chat response has products
func (cr *ChatResponse) HasProducts() bool {
	return len(cr.Products) > 0
}

// IsDialogue checks if Gemini response is dialogue
func (gr *GeminiResponse) IsDialogue() bool {
	return gr.ResponseType == "dialogue"
}

// IsSearch checks if Gemini response is search
func (gr *GeminiResponse) IsSearch() bool {
	return gr.ResponseType == "search"
}
