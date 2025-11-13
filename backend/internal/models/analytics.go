// backend/internal/models/analytics.go
package models

// UserRecommendations contains personalized recommendations for a user
type UserRecommendations struct {
	PreferredCategories  []string               `json:"preferred_categories"`
	PreferredBrands      []string               `json:"preferred_brands"`
	CommunicationStyle   string                 `json:"communication_style"`
	TotalSessions        int                    `json:"total_sessions"`
	SuccessRate          float64                `json:"success_rate"`
	CommonKeywords       []string               `json:"common_keywords"`
	InteractionStats     map[string]interface{} `json:"interaction_stats"`
}

// SessionInsights contains insights about a specific chat session
type SessionInsights struct {
	SessionID            string   `json:"session_id"`
	MessageCount         int      `json:"message_count"`
	SearchCount          int      `json:"search_count"`
	ProductsShown        int      `json:"products_shown"`
	ProductsClicked      int      `json:"products_clicked"`
	PrimaryIntent        string   `json:"primary_intent"`
	Sentiment            string   `json:"sentiment"`
	SentimentScore       float64  `json:"sentiment_score"`
	KeyTopics            []string `json:"key_topics"`
	CategoriesExplored   []string `json:"categories_explored"`
	FlowQualityScore     float64  `json:"flow_quality_score"`
	FoundProduct         bool     `json:"found_product"`
	UserSatisfied        *bool    `json:"user_satisfied"`
	InteractionsCount    int      `json:"interactions_count"`
	Summary              *string  `json:"summary"`
}

// UserBehaviorSummary contains a summary of user behavior profile
type UserBehaviorSummary struct {
	UserID                  string                 `json:"user_id"`
	CommunicationStyle      string                 `json:"communication_style"`
	TotalSessions           int                    `json:"total_sessions"`
	AvgSessionDuration      float64                `json:"avg_session_duration_mins"`
	AvgMessagesPerSession   float64                `json:"avg_messages_per_session"`
	SuccessRate             float64                `json:"success_rate"`
	TotalProductsViewed     int                    `json:"total_products_viewed"`
	TotalProductsClicked    int                    `json:"total_products_clicked"`
	TopCategories           []string               `json:"top_categories"`
	TopBrands               []string               `json:"top_brands"`
	CommonKeywords          []string               `json:"common_keywords"`
	CategoryPreferences     map[string]float64     `json:"category_preferences"`
	PriceRanges             map[string]interface{} `json:"price_ranges"`
	LastUpdated             string                 `json:"last_updated"`
}

// ConversationAnalyticsSummary contains summary analytics for multiple sessions
type ConversationAnalyticsSummary struct {
	TotalSessions        int                    `json:"total_sessions"`
	AvgMessageCount      float64                `json:"avg_message_count"`
	AvgSearchCount       float64                `json:"avg_search_count"`
	AvgProductsShown     float64                `json:"avg_products_shown"`
	TotalProductsClicked int                    `json:"total_products_clicked"`
	SuccessRate          float64                `json:"success_rate"`
	SentimentDistribution map[string]int        `json:"sentiment_distribution"`
	IntentDistribution   map[string]int         `json:"intent_distribution"`
	TopCategories        []string               `json:"top_categories"`
	AvgFlowQuality       float64                `json:"avg_flow_quality"`
	RecentSessions       []*SessionInsights     `json:"recent_sessions"`
}

// ProductInteractionSummary contains summary of product interactions
type ProductInteractionSummary struct {
	TotalInteractions int                    `json:"total_interactions"`
	ViewCount         int                    `json:"view_count"`
	ClickCount        int                    `json:"click_count"`
	ComparisonCount   int                    `json:"comparison_count"`
	DismissalCount    int                    `json:"dismissal_count"`
	AvgImplicitScore  float64                `json:"avg_implicit_score"`
	EngagementRate    float64                `json:"engagement_rate"`
	TopProducts       []string               `json:"top_products"`
	TopCategories     []string               `json:"top_categories"`
	TopBrands         []string               `json:"top_brands"`
}
