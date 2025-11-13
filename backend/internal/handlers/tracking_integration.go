// backend/internal/handlers/tracking_integration.go
package handlers

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"

	"mylittleprice/internal/container"
	"mylittleprice/internal/models"
	"mylittleprice/internal/utils"
)

// TrackingMiddleware provides tracking integration for chat processing
type TrackingMiddleware struct {
	container *container.Container
}

// NewTrackingMiddleware creates a new tracking middleware
func NewTrackingMiddleware(c *container.Container) *TrackingMiddleware {
	return &TrackingMiddleware{
		container: c,
	}
}

// TrackMessageAnalysis analyzes and tracks a user message
func (t *TrackingMiddleware) TrackMessageAnalysis(
	ctx context.Context,
	sessionID uuid.UUID,
	userID *uuid.UUID,
	message string,
) {
	// Analyze intent
	intent := t.container.MessageAnalysisService.DetectIntent(message)

	utils.LogInfo(ctx, "message intent detected",
		slog.String("session_id", sessionID.String()),
		slog.String("intent", intent.Intent),
		slog.Float64("confidence", intent.Confidence),
	)

	// Increment message count in analytics
	if err := t.container.ConversationAnalyticsService.IncrementMetric(sessionID, "user_message_count"); err != nil {
		utils.LogWarn(ctx, "failed to increment user message count", slog.Any("error", err))
	}

	// Extract and track categories from intent context
	if category, ok := intent.Context["category"].(string); ok && category != "" {
		if err := t.container.ConversationAnalyticsService.AddCategory(sessionID, category); err != nil {
			utils.LogWarn(ctx, "failed to add category", slog.Any("error", err))
		}
	}

	// Extract and track brands
	if brands, ok := intent.Context["brands"].([]string); ok {
		for _, brand := range brands {
			if err := t.container.ConversationAnalyticsService.AddBrandMention(sessionID, brand); err != nil {
				utils.LogWarn(ctx, "failed to add brand mention", slog.Any("error", err))
			}
		}
	}

	// Extract and track price mentions
	if prices, ok := intent.Context["prices"].([]float64); ok {
		for _, price := range prices {
			if err := t.container.ConversationAnalyticsService.AddPriceMention(sessionID, price); err != nil {
				utils.LogWarn(ctx, "failed to add price mention", slog.Any("error", err))
			}
		}
	}
}

// TrackProductViews tracks when products are shown to user
func (t *TrackingMiddleware) TrackProductViews(
	ctx context.Context,
	sessionID string,
	userID *uuid.UUID,
	products []models.ProductCard,
	searchQuery string,
	searchType string,
	messagePosition int,
) {
	sessionUUID, err := uuid.Parse(sessionID)
	if err != nil {
		utils.LogWarn(ctx, "invalid session ID for product tracking", slog.String("session_id", sessionID))
		return
	}

	// Track each product view
	for i, product := range products {
		productInfo := &models.ProductInfo{
			ID:       product.ID,
			Name:     product.Name,
			Price:    product.Price,
			Currency: product.Currency,
			Rating:   product.Rating,
			Source:   product.Source,
			URL:      product.URL,
		}

		if err := t.container.ProductInteractionService.TrackProductView(
			userID,
			sessionID,
			productInfo,
			messagePosition,
			i+1, // Position in results (1-indexed)
			searchQuery,
			searchType,
		); err != nil {
			utils.LogWarn(ctx, "failed to track product view",
				slog.String("product_id", product.ID),
				slog.Any("error", err),
			)
		}
	}

	// Update products shown count in analytics
	if err := t.container.ConversationAnalyticsService.UpdateSessionMetrics(sessionUUID, map[string]interface{}{
		"products_shown": len(products),
	}); err != nil {
		utils.LogWarn(ctx, "failed to update products shown", slog.Any("error", err))
	}

	// Mark that user found products
	if len(products) > 0 {
		if err := t.container.ConversationAnalyticsService.UpdateSessionMetrics(sessionUUID, map[string]interface{}{
			"found_product": true,
		}); err != nil {
			utils.LogWarn(ctx, "failed to mark product found", slog.Any("error", err))
		}
	}

	utils.LogInfo(ctx, "tracked product views",
		slog.Int("count", len(products)),
		slog.String("session_id", sessionID),
	)
}

// TrackSearchQuery tracks search queries
func (t *TrackingMiddleware) TrackSearchQuery(
	ctx context.Context,
	sessionID uuid.UUID,
	searchPhrase string,
	searchType string,
	category string,
) {
	// Increment search count
	if err := t.container.ConversationAnalyticsService.IncrementMetric(sessionID, "search_count"); err != nil {
		utils.LogWarn(ctx, "failed to increment search count", slog.Any("error", err))
	}

	// Add category if present
	if category != "" {
		if err := t.container.ConversationAnalyticsService.AddCategory(sessionID, category); err != nil {
			utils.LogWarn(ctx, "failed to add search category", slog.Any("error", err))
		}
	}

	utils.LogInfo(ctx, "tracked search query",
		slog.String("session_id", sessionID.String()),
		slog.String("search_phrase", searchPhrase),
		slog.String("search_type", searchType),
	)
}

// InitializeSessionAnalytics initializes analytics tracking for a new session
func (t *TrackingMiddleware) InitializeSessionAnalytics(
	ctx context.Context,
	sessionID uuid.UUID,
	userID *uuid.UUID,
) error {
	// Check if analytics already exist
	_, err := t.container.ConversationAnalyticsService.GetSessionAnalytics(sessionID)
	if err == nil {
		// Analytics already exist
		return nil
	}

	// Create new analytics
	_, err = t.container.ConversationAnalyticsService.StartSession(sessionID, userID)
	if err != nil {
		return fmt.Errorf("failed to initialize session analytics: %w", err)
	}

	utils.LogInfo(ctx, "initialized session analytics",
		slog.String("session_id", sessionID.String()),
	)

	return nil
}

// FinalizeSession finalizes analytics and triggers learning
func (t *TrackingMiddleware) FinalizeSession(
	ctx context.Context,
	sessionID uuid.UUID,
	userID *uuid.UUID,
) error {
	// Finalize conversation analytics
	if err := t.container.ConversationAnalyticsService.EndSession(sessionID); err != nil {
		utils.LogWarn(ctx, "failed to end session analytics", slog.Any("error", err))
		return err
	}

	// If user is authenticated, trigger behavior learning
	if userID != nil {
		if err := t.container.UserBehaviorService.LearnFromSession(*userID, sessionID); err != nil {
			utils.LogWarn(ctx, "failed to learn from session", slog.Any("error", err))
			// Don't return error, learning is best-effort
		} else {
			utils.LogInfo(ctx, "learned from session",
				slog.String("user_id", userID.String()),
				slog.String("session_id", sessionID.String()),
			)
		}
	}

	return nil
}

// GetUserRecommendations gets personalized recommendations for a user
func (t *TrackingMiddleware) GetUserRecommendations(
	ctx context.Context,
	userID uuid.UUID,
) (*models.UserRecommendations, error) {
	// Get user behavior profile
	profile, err := t.container.UserBehaviorService.GetCachedProfile(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}

	// Get top categories
	topCategories, _ := t.container.UserBehaviorService.GetRecommendedCategories(userID, 5)

	// Get top brands
	topBrands, _ := t.container.UserBehaviorService.GetPreferredBrands(userID, 5)

	// Get interaction stats
	interactionStats, _ := t.container.ProductInteractionService.GetInteractionStats(userID)

	recommendations := &models.UserRecommendations{
		PreferredCategories:  topCategories,
		PreferredBrands:      topBrands,
		CommunicationStyle:   profile.CommunicationStyle,
		TotalSessions:        profile.TotalSessions,
		SuccessRate:          profile.SuccessRate,
		CommonKeywords:       profile.CommonKeywords,
		InteractionStats:     interactionStats,
	}

	return recommendations, nil
}

// GetSessionInsights gets insights about a specific session
func (t *TrackingMiddleware) GetSessionInsights(
	ctx context.Context,
	sessionID uuid.UUID,
) (*models.SessionInsights, error) {
	// Get session analytics
	analytics, err := t.container.ConversationAnalyticsService.GetSessionAnalytics(sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get session analytics: %w", err)
	}

	// Get session interactions
	interactions, err := t.container.ProductInteractionService.GetSessionInteractions(sessionID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get session interactions: %w", err)
	}

	insights := &models.SessionInsights{
		SessionID:            sessionID.String(),
		MessageCount:         analytics.MessageCount,
		SearchCount:          analytics.SearchCount,
		ProductsShown:        analytics.ProductsShown,
		ProductsClicked:      analytics.ProductsClicked,
		PrimaryIntent:        analytics.PrimaryIntent,
		Sentiment:            analytics.OverallSentiment,
		SentimentScore:       analytics.SentimentScore,
		KeyTopics:            analytics.KeyTopics,
		CategoriesExplored:   analytics.CategoriesExplored,
		FlowQualityScore:     analytics.FlowQualityScore,
		FoundProduct:         analytics.FoundProduct,
		UserSatisfied:        analytics.UserSatisfied,
		InteractionsCount:    len(interactions),
		Summary:              analytics.Summary,
	}

	return insights, nil
}
