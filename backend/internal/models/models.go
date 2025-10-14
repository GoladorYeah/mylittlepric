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
	Status          SearchStatus `json:"status"`
	Category        string       `json:"category"`
	ProductType     string       `json:"product_type,omitempty"`
	Brand           string       `json:"brand,omitempty"`
	CollectedParams []string     `json:"collected_params,omitempty"`
	LastSearchTime  time.Time    `json:"last_search_time,omitempty"`
	SearchCount     int          `json:"search_count"`
	LastProduct     *ProductInfo `json:"last_product,omitempty"`
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
	Category     string                 `json:"category,omitempty" db:"category"`
	Metadata     map[string]interface{} `json:"metadata,omitempty" db:"metadata"`
	CreatedAt    time.Time              `json:"created_at" db:"created_at"`
}

type ChatRequest struct {
	SessionID string `json:"session_id"`
	Message   string `json:"message"`
	Country   string `json:"country"`
	Language  string `json:"language"`
	NewSearch bool   `json:"new_search"`
}

type ChatResponse struct {
	Type         string           `json:"type"`
	Output       string           `json:"output,omitempty"`
	QuickReplies []string         `json:"quick_replies,omitempty"`
	Products     []ProductCard    `json:"products,omitempty"`
	SearchType   string           `json:"search_type,omitempty"`
	SessionID    string           `json:"session_id"`
	MessageCount int              `json:"message_count"`
	SearchState  *SearchStateInfo `json:"search_state,omitempty"`
}

type SearchStateInfo struct {
	Status      SearchStatus `json:"status"`
	Category    string       `json:"category,omitempty"`
	CanContinue bool         `json:"can_continue"`
	SearchCount int          `json:"search_count"`
	MaxSearches int          `json:"max_searches"`
	Message     string       `json:"message,omitempty"`
}

type GeminiResponse struct {
	ResponseType string   `json:"response_type"`
	Output       string   `json:"output,omitempty"`
	QuickReplies []string `json:"quick_replies,omitempty"`
	SearchPhrase string   `json:"search_phrase,omitempty"`
	SearchType   string   `json:"search_type,omitempty"`
	Category     string   `json:"category,omitempty"`
	PriceFilter  string   `json:"price_filter,omitempty"`
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

type Product struct {
	ID             uuid.UUID              `json:"id" db:"id"`
	PageToken      string                 `json:"page_token" db:"page_token"`
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

type ProductDetailsRequest struct {
	PageToken string `json:"page_token"`
	Country   string `json:"country"`
}

type ProductDetailsResponse struct {
	Type            string               `json:"type"`
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

type ProductSpec struct {
	Title string `json:"title"`
	Value string `json:"value"`
}

type ProductVariant struct {
	Title string        `json:"title"`
	Items []VariantItem `json:"items"`
}

type VariantItem struct {
	Name      string `json:"name"`
	Selected  bool   `json:"selected"`
	Available bool   `json:"available"`
	PageToken string `json:"page_token"`
}

type ProductOffer struct {
	Merchant     string  `json:"merchant"`
	Price        string  `json:"price"`
	Currency     string  `json:"currency"`
	Link         string  `json:"link"`
	Availability string  `json:"availability,omitempty"`
	Shipping     string  `json:"shipping,omitempty"`
	Rating       float32 `json:"rating,omitempty"`
}

type ProductVideo struct {
	Title     string `json:"title"`
	Link      string `json:"link"`
	Source    string `json:"source"`
	Channel   string `json:"channel,omitempty"`
	Duration  string `json:"duration,omitempty"`
	Thumbnail string `json:"thumbnail"`
}

type AlternativeProduct struct {
	Title     string  `json:"title"`
	Thumbnail string  `json:"thumbnail"`
	Price     string  `json:"price"`
	Rating    float32 `json:"rating,omitempty"`
	Reviews   int     `json:"reviews,omitempty"`
	PageToken string  `json:"page_token"`
}

type RatingBreakdown struct {
	Stars  int `json:"stars"`
	Amount int `json:"amount"`
}

type SearchQuery struct {
	ID             uuid.UUID `json:"id" db:"id"`
	SessionID      uuid.UUID `json:"session_id,omitempty" db:"session_id"`
	OriginalQuery  string    `json:"original_query" db:"original_query"`
	OptimizedQuery string    `json:"optimized_query" db:"optimized_query"`
	SearchType     string    `json:"search_type" db:"search_type"`
	CountryCode    string    `json:"country_code" db:"country_code"`
	ResultCount    int       `json:"result_count" db:"result_count"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

type APIUsage struct {
	ID             uuid.UUID `json:"id" db:"id"`
	APIName        string    `json:"api_name" db:"api_name"`
	KeyIndex       int       `json:"key_index" db:"key_index"`
	RequestType    string    `json:"request_type,omitempty" db:"request_type"`
	ResponseTimeMs int       `json:"response_time_ms,omitempty" db:"response_time_ms"`
	Success        bool      `json:"success" db:"success"`
	ErrorMessage   string    `json:"error_message,omitempty" db:"error_message"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

type HealthResponse struct {
	Status    string                 `json:"status"`
	Timestamp time.Time              `json:"timestamp"`
	Services  map[string]ServiceInfo `json:"services"`
}

type ServiceInfo struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    string `json:"code,omitempty"`
}

func (s SearchStatus) IsIdle() bool {
	return s == SearchStatusIdle
}

func (s SearchStatus) IsInProgress() bool {
	return s == SearchStatusInProgress
}

func (s SearchStatus) IsCompleted() bool {
	return s == SearchStatusCompleted
}

func (s SearchStatus) String() string {
	return string(s)
}

func (ss *SearchState) HasCategory() bool {
	return ss.Category != ""
}

func (ss *SearchState) HasBrand() bool {
	return ss.Brand != ""
}

func (ss *SearchState) HasProductType() bool {
	return ss.ProductType != ""
}

func (ss *SearchState) HasLastProduct() bool {
	return ss.LastProduct != nil
}

func (cr *ChatResponse) HasProducts() bool {
	return len(cr.Products) > 0
}

func (gr *GeminiResponse) IsDialogue() bool {
	return gr.ResponseType == "dialogue"
}

func (gr *GeminiResponse) IsSearch() bool {
	return gr.ResponseType == "search"
}
