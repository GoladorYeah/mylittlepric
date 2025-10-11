package constants

// ═══════════════════════════════════════════════════════════
// APPLICATION CONSTANTS
// ═══════════════════════════════════════════════════════════

const (
	// Application Info
	AppName    = "MyLittlePrice"
	AppVersion = "2.0.0"
)

// ═══════════════════════════════════════════════════════════
// CONVERSATION LIMITS
// ═══════════════════════════════════════════════════════════

const (
	// MaxQuestionsBeforeForceSearch forces search after this many dialogue messages
	MaxQuestionsBeforeForceSearch = 6

	// MaxHistoryMessages limits conversation history sent to AI
	MaxHistoryMessages = 4

	// MaxRecentMessagesForContext number of recent messages to analyze
	MaxRecentMessagesForContext = 4
)

// ═══════════════════════════════════════════════════════════
// SEARCH RELEVANCE THRESHOLDS
// ═══════════════════════════════════════════════════════════

const (
	// MinRelevanceThresholdExact for exact product search
	MinRelevanceThresholdExact = 0.7

	// MinRelevanceThresholdParameters for parametric search
	MinRelevanceThresholdParameters = 0.5

	// MinRelevanceThresholdCategory for category search
	MinRelevanceThresholdCategory = 0.3
)

// ═══════════════════════════════════════════════════════════
// PRODUCT DISPLAY LIMITS
// ═══════════════════════════════════════════════════════════

const (
	// MaxProductsExactSearch for exact product matches
	MaxProductsExactSearch = 3

	// MaxProductsParameterSearch for parametric searches
	MaxProductsParameterSearch = 6

	// MaxProductsCategorySearch for category searches
	MaxProductsCategorySearch = 8
)

// ═══════════════════════════════════════════════════════════
// GROUNDING CONFIGURATION
// ═══════════════════════════════════════════════════════════

const (
	// MinMessagesForAdvancedDialogue (conservative mode)
	MinMessagesForAdvancedDialogueConservative = 6

	// MinMessagesForAdvancedDialogue (balanced mode)
	MinMessagesForAdvancedDialogueBalanced = 4

	// MinMessagesForAdvancedDialogue (aggressive mode)
	MinMessagesForAdvancedDialogueAggressive = 3
)

// ═══════════════════════════════════════════════════════════
// CACHE KEY PREFIXES
// ═══════════════════════════════════════════════════════════

const (
	CachePrefixSession    = "session:"
	CachePrefixMessages   = "session:%s:messages"
	CachePrefixProduct    = "product:"
	CachePrefixSearch     = "search:"
	CachePrefixGemini     = "gemini:"
	CachePrefixKeyRotator = "keyrotator:"
)

// ═══════════════════════════════════════════════════════════
// ERROR CODES
// ═══════════════════════════════════════════════════════════

const (
	ErrCodeSessionNotFound    = "SESSION_NOT_FOUND"
	ErrCodeSessionExpired     = "SESSION_EXPIRED"
	ErrCodeMaxSearchesReached = "MAX_SEARCHES_REACHED"
	ErrCodeMaxMessagesReached = "MAX_MESSAGES_REACHED"
	ErrCodeInvalidRequest     = "INVALID_REQUEST"
	ErrCodeValidationError    = "VALIDATION_ERROR"
	ErrCodeAIError            = "AI_ERROR"
	ErrCodeSearchError        = "SEARCH_ERROR"
	ErrCodeCacheError         = "CACHE_ERROR"
	ErrCodeInternalError      = "INTERNAL_ERROR"
)

// ═══════════════════════════════════════════════════════════
// DEFAULT VALUES
// ═══════════════════════════════════════════════════════════

const (
	DefaultCountry  = "CH"
	DefaultLanguage = "de"
	DefaultCurrency = "CHF"
)

// ═══════════════════════════════════════════════════════════
// VALIDATION LIMITS
// ═══════════════════════════════════════════════════════════

const (
	MinQueryLength = 2
	MaxQueryLength = 200

	MinMeaningfulParamLength = 3
)

// ═══════════════════════════════════════════════════════════
// HTTP STATUS MESSAGES
// ═══════════════════════════════════════════════════════════

const (
	MsgSearchBlocked          = "This search is complete. To search for another product, please click 'New Search' button."
	MsgMaxSearchesReached     = "Maximum searches per session reached. Please start a new session."
	MsgNoProductsFound        = "Sorry, I couldn't find any products matching your criteria. Could you try describing what you're looking for differently?"
	MsgProductDetailsNotFound = "Product details not found."
)

// ═══════════════════════════════════════════════════════════
// RESPONSE TYPES
// ═══════════════════════════════════════════════════════════

const (
	ResponseTypeText           = "text"
	ResponseTypeProductCard    = "product_card"
	ResponseTypeSearchBlocked  = "search_blocked"
	ResponseTypeProductDetails = "product_details"
)

// ═══════════════════════════════════════════════════════════
// SEARCH TYPES
// ═══════════════════════════════════════════════════════════

const (
	SearchTypeExact      = "exact"
	SearchTypeParameters = "parameters"
	SearchTypeCategory   = "category"
)
