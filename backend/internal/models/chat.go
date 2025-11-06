package models

// ═══════════════════════════════════════════════════════════
// CHAT REQUEST/RESPONSE MODELS
// ═══════════════════════════════════════════════════════════

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

// ═══════════════════════════════════════════════════════════
// AI RESPONSE MODELS
// ═══════════════════════════════════════════════════════════

type GeminiResponse struct {
	ResponseType  string   `json:"response_type"` // "dialogue", "search", or "api_request"
	Output        string   `json:"output"`
	QuickReplies  []string `json:"quick_replies"`
	SearchPhrase  string   `json:"search_phrase"` // For response_type="search"
	SearchType    string   `json:"search_type"`   // "exact", "parameters", or "category"
	Category      string   `json:"category"`
	PriceFilter   string   `json:"price_filter,omitempty"` // "cheaper" or "expensive"
	MinPrice      *float64 `json:"min_price,omitempty"`    // Minimum price in user's currency
	MaxPrice      *float64 `json:"max_price,omitempty"`    // Maximum price in user's currency
	ProductType   string   `json:"product_type"`
	Brand         string   `json:"brand"`
	Confidence    float32  `json:"confidence"`
	RequiresInput bool     `json:"requires_input"`
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

// ═══════════════════════════════════════════════════════════
// ERROR MODELS
// ═══════════════════════════════════════════════════════════

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
