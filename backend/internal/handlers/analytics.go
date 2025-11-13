// backend/internal/handlers/analytics.go
package handlers

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"mylittleprice/internal/container"
	"mylittleprice/internal/models"
	"mylittleprice/internal/utils"
)

// AnalyticsHandler provides endpoints for analytics and insights
type AnalyticsHandler struct {
	container *container.Container
	tracking  *TrackingMiddleware
}

// NewAnalyticsHandler creates a new analytics handler
func NewAnalyticsHandler(c *container.Container) *AnalyticsHandler {
	return &AnalyticsHandler{
		container: c,
		tracking:  NewTrackingMiddleware(c),
	}
}

// GetUserBehaviorProfile returns the user's behavior profile
// GET /api/analytics/profile
func (h *AnalyticsHandler) GetUserBehaviorProfile(c *fiber.Ctx) error {
	ctx := c.Context()

	// Get user ID from JWT
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		utils.LogWarn(ctx, "unauthorized access to behavior profile")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Get behavior profile summary
	summary, err := h.container.UserBehaviorService.GetProfileSummary(userID)
	if err != nil {
		utils.LogError(ctx, "failed to get behavior profile", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get behavior profile",
		})
	}

	return c.JSON(summary)
}

// GetUserRecommendations returns personalized recommendations
// GET /api/analytics/recommendations
func (h *AnalyticsHandler) GetUserRecommendations(c *fiber.Ctx) error {
	ctx := c.Context()

	// Get user ID from JWT
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		utils.LogWarn(ctx, "unauthorized access to recommendations")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Get recommendations
	recommendations, err := h.tracking.GetUserRecommendations(ctx, userID)
	if err != nil {
		utils.LogError(ctx, "failed to get recommendations", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get recommendations",
		})
	}

	return c.JSON(recommendations)
}

// GetSessionInsights returns insights about a specific session
// GET /api/analytics/session/:sessionId
func (h *AnalyticsHandler) GetSessionInsights(c *fiber.Ctx) error {
	ctx := c.Context()

	// Get session ID from URL parameter
	sessionIDStr := c.Params("sessionId")
	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid session ID",
		})
	}

	// Get insights
	insights, err := h.tracking.GetSessionInsights(ctx, sessionID)
	if err != nil {
		utils.LogError(ctx, "failed to get session insights", err,
			slog.String("session_id", sessionIDStr),
		)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get session insights",
		})
	}

	return c.JSON(insights)
}

// GetUserAnalytics returns analytics summary for the user
// GET /api/analytics/summary
func (h *AnalyticsHandler) GetUserAnalytics(c *fiber.Ctx) error {
	ctx := c.Context()

	// Get user ID from JWT
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		utils.LogWarn(ctx, "unauthorized access to analytics")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Get user analytics (last 20 sessions)
	analyticsData, err := h.container.ConversationAnalyticsService.GetUserAnalytics(userID, 20)
	if err != nil {
		utils.LogError(ctx, "failed to get user analytics", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get analytics",
		})
	}

	// Calculate summary statistics
	summary := &models.ConversationAnalyticsSummary{
		TotalSessions:         len(analyticsData),
		SentimentDistribution: make(map[string]int),
		IntentDistribution:    make(map[string]int),
		TopCategories:         make([]string, 0),
		RecentSessions:        make([]*models.SessionInsights, 0),
	}

	totalMessages := 0
	totalSearches := 0
	totalProducts := 0
	totalClicks := 0
	totalFlowQuality := 0.0
	successCount := 0

	categoryCount := make(map[string]int)

	for _, analytics := range analyticsData {
		totalMessages += analytics.MessageCount
		totalSearches += analytics.SearchCount
		totalProducts += analytics.ProductsShown
		totalClicks += analytics.ProductsClicked
		totalFlowQuality += analytics.FlowQualityScore

		// Count sentiment
		summary.SentimentDistribution[analytics.OverallSentiment]++

		// Count intent
		if analytics.PrimaryIntent != "" {
			summary.IntentDistribution[analytics.PrimaryIntent]++
		}

		// Count categories
		for _, cat := range analytics.CategoriesExplored {
			categoryCount[cat]++
		}

		// Count success
		if analytics.FoundProduct {
			successCount++
		}

		// Add to recent sessions (only first 5)
		if len(summary.RecentSessions) < 5 {
			sessionInsights := &models.SessionInsights{
				SessionID:          analytics.SessionID.String(),
				MessageCount:       analytics.MessageCount,
				SearchCount:        analytics.SearchCount,
				ProductsShown:      analytics.ProductsShown,
				ProductsClicked:    analytics.ProductsClicked,
				PrimaryIntent:      analytics.PrimaryIntent,
				Sentiment:          analytics.OverallSentiment,
				SentimentScore:     analytics.SentimentScore,
				KeyTopics:          analytics.KeyTopics,
				CategoriesExplored: analytics.CategoriesExplored,
				FlowQualityScore:   analytics.FlowQualityScore,
				FoundProduct:       analytics.FoundProduct,
				UserSatisfied:      analytics.UserSatisfied,
				Summary:            analytics.Summary,
			}
			summary.RecentSessions = append(summary.RecentSessions, sessionInsights)
		}
	}

	// Calculate averages
	if len(analyticsData) > 0 {
		summary.AvgMessageCount = float64(totalMessages) / float64(len(analyticsData))
		summary.AvgSearchCount = float64(totalSearches) / float64(len(analyticsData))
		summary.AvgProductsShown = float64(totalProducts) / float64(len(analyticsData))
		summary.AvgFlowQuality = totalFlowQuality / float64(len(analyticsData))
		summary.SuccessRate = float64(successCount) / float64(len(analyticsData))
	}

	summary.TotalProductsClicked = totalClicks

	// Get top 5 categories
	type catCount struct {
		Name  string
		Count int
	}
	var categories []catCount
	for cat, count := range categoryCount {
		categories = append(categories, catCount{Name: cat, Count: count})
	}
	// Sort by count (simple bubble sort for small dataset)
	for i := 0; i < len(categories); i++ {
		for j := i + 1; j < len(categories); j++ {
			if categories[j].Count > categories[i].Count {
				categories[i], categories[j] = categories[j], categories[i]
			}
		}
	}
	for i := 0; i < 5 && i < len(categories); i++ {
		summary.TopCategories = append(summary.TopCategories, categories[i].Name)
	}

	return c.JSON(summary)
}

// GetProductInteractionStats returns product interaction statistics
// GET /api/analytics/interactions
func (h *AnalyticsHandler) GetProductInteractionStats(c *fiber.Ctx) error {
	ctx := c.Context()

	// Get user ID from JWT
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		utils.LogWarn(ctx, "unauthorized access to interaction stats")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Get interaction stats
	stats, err := h.container.ProductInteractionService.GetInteractionStats(userID)
	if err != nil {
		utils.LogError(ctx, "failed to get interaction stats", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get interaction stats",
		})
	}

	return c.JSON(stats)
}

// TrackProductClick tracks when user clicks on a product
// POST /api/analytics/track/click
func (h *AnalyticsHandler) TrackProductClick(c *fiber.Ctx) error {
	ctx := c.Context()

	var req struct {
		SessionID string `json:"session_id"`
		ProductID string `json:"product_id"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Get user ID (optional, may be anonymous)
	var userID *uuid.UUID
	if uid, err := utils.GetUserIDFromContext(c); err == nil {
		userID = &uid
	}

	// Track click
	if err := h.container.ProductInteractionService.TrackProductClick(userID, req.SessionID, req.ProductID); err != nil {
		utils.LogError(ctx, "failed to track product click", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to track click",
		})
	}

	// Also increment in analytics
	if sessionUUID, err := uuid.Parse(req.SessionID); err == nil {
		h.container.ConversationAnalyticsService.IncrementMetric(sessionUUID, "products_clicked")
	}

	return c.JSON(fiber.Map{
		"success": true,
	})
}

// FinalizeSession manually finalizes a session (triggers learning)
// POST /api/analytics/finalize/:sessionId
func (h *AnalyticsHandler) FinalizeSession(c *fiber.Ctx) error {
	ctx := c.Context()

	// Get session ID from URL parameter
	sessionIDStr := c.Params("sessionId")
	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid session ID",
		})
	}

	// Get user ID (optional)
	var userID *uuid.UUID
	if uid, err := utils.GetUserIDFromContext(c); err == nil {
		userID = &uid
	}

	// Finalize session
	if err := h.tracking.FinalizeSession(ctx, sessionID, userID); err != nil {
		utils.LogError(ctx, "failed to finalize session", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to finalize session",
		})
	}

	utils.LogInfo(ctx, "session finalized",
		slog.String("session_id", sessionIDStr),
	)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Session finalized and learning complete",
	})
}
