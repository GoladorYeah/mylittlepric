// backend/internal/services/conversation_analytics.go
package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"mylittleprice/ent"
	"mylittleprice/ent/chatsession"
	"mylittleprice/ent/conversationanalytics"
	"mylittleprice/ent/message"
)

// ConversationAnalyticsService tracks and analyzes chat sessions
type ConversationAnalyticsService struct {
	db        *ent.Client
	embedding *EmbeddingService
	ctx       context.Context
}

// NewConversationAnalyticsService creates a new conversation analytics service
func NewConversationAnalyticsService(db *ent.Client, embedding *EmbeddingService) *ConversationAnalyticsService {
	return &ConversationAnalyticsService{
		db:        db,
		embedding: embedding,
		ctx:       context.Background(),
	}
}

// StartSession initializes analytics tracking for a new session
func (s *ConversationAnalyticsService) StartSession(sessionID uuid.UUID, userID *uuid.UUID) (*ent.ConversationAnalytics, error) {
	builder := s.db.ConversationAnalytics.Create().
		SetSessionID(sessionID).
		SetSessionStartedAt(time.Now())

	if userID != nil {
		builder.SetUserID(*userID)
	}

	analytics, err := builder.Save(s.ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create conversation analytics: %w", err)
	}

	fmt.Printf("📈 Started analytics tracking for session %s\n", sessionID)
	return analytics, nil
}

// UpdateSessionMetrics updates metrics during the session
func (s *ConversationAnalyticsService) UpdateSessionMetrics(sessionID uuid.UUID, updates map[string]interface{}) error {
	analytics, err := s.db.ConversationAnalytics.
		Query().
		Where(conversationanalytics.SessionID(sessionID)).
		Only(s.ctx)

	if err != nil {
		return fmt.Errorf("failed to get analytics: %w", err)
	}

	updater := s.db.ConversationAnalytics.UpdateOne(analytics)

	// Apply updates
	if val, ok := updates["message_count"].(int); ok {
		updater.SetMessageCount(val)
	}
	if val, ok := updates["user_message_count"].(int); ok {
		updater.SetUserMessageCount(val)
	}
	if val, ok := updates["assistant_message_count"].(int); ok {
		updater.SetAssistantMessageCount(val)
	}
	if val, ok := updates["search_count"].(int); ok {
		updater.SetSearchCount(val)
	}
	if val, ok := updates["products_shown"].(int); ok {
		updater.SetProductsShown(val)
	}
	if val, ok := updates["products_clicked"].(int); ok {
		updater.SetProductsClicked(val)
	}
	if val, ok := updates["found_product"].(bool); ok {
		updater.SetFoundProduct(val)
	}
	if val, ok := updates["clarification_count"].(int); ok {
		updater.SetClarificationCount(val)
	}
	if val, ok := updates["search_refinement_count"].(int); ok {
		updater.SetSearchRefinementCount(val)
	}

	_, err = updater.Save(s.ctx)
	return err
}

// IncrementMetric increments a counter metric
func (s *ConversationAnalyticsService) IncrementMetric(sessionID uuid.UUID, metric string) error {
	analytics, err := s.db.ConversationAnalytics.
		Query().
		Where(conversationanalytics.SessionID(sessionID)).
		Only(s.ctx)

	if err != nil {
		return fmt.Errorf("failed to get analytics: %w", err)
	}

	updater := s.db.ConversationAnalytics.UpdateOne(analytics)

	switch metric {
	case "message_count":
		updater.SetMessageCount(analytics.MessageCount + 1)
	case "user_message_count":
		updater.SetUserMessageCount(analytics.UserMessageCount + 1)
	case "assistant_message_count":
		updater.SetAssistantMessageCount(analytics.AssistantMessageCount + 1)
	case "search_count":
		updater.SetSearchCount(analytics.SearchCount + 1)
	case "products_shown":
		updater.SetProductsShown(analytics.ProductsShown + 1)
	case "products_clicked":
		updater.SetProductsClicked(analytics.ProductsClicked + 1)
	case "clarification_count":
		updater.SetClarificationCount(analytics.ClarificationCount + 1)
	case "search_refinement_count":
		updater.SetSearchRefinementCount(analytics.SearchRefinementCount + 1)
	default:
		return fmt.Errorf("unknown metric: %s", metric)
	}

	_, err = updater.Save(s.ctx)
	return err
}

// AddCategory adds a category to the explored categories list
func (s *ConversationAnalyticsService) AddCategory(sessionID uuid.UUID, category string) error {
	analytics, err := s.db.ConversationAnalytics.
		Query().
		Where(conversationanalytics.SessionID(sessionID)).
		Only(s.ctx)

	if err != nil {
		return fmt.Errorf("failed to get analytics: %w", err)
	}

	// Add category if not already present
	categories := analytics.CategoriesExplored
	found := false
	for _, cat := range categories {
		if cat == category {
			found = true
			break
		}
	}

	if !found {
		categories = append(categories, category)
		_, err = s.db.ConversationAnalytics.UpdateOne(analytics).
			SetCategoriesExplored(categories).
			Save(s.ctx)
	}

	return err
}

// AddBrandMention adds a brand to the brand mentions list
func (s *ConversationAnalyticsService) AddBrandMention(sessionID uuid.UUID, brand string) error {
	analytics, err := s.db.ConversationAnalytics.
		Query().
		Where(conversationanalytics.SessionID(sessionID)).
		Only(s.ctx)

	if err != nil {
		return fmt.Errorf("failed to get analytics: %w", err)
	}

	// Add brand if not already present
	brands := analytics.BrandMentions
	found := false
	for _, b := range brands {
		if b == brand {
			found = true
			break
		}
	}

	if !found {
		brands = append(brands, brand)
		_, err = s.db.ConversationAnalytics.UpdateOne(analytics).
			SetBrandMentions(brands).
			Save(s.ctx)
	}

	return err
}

// AddPriceMention adds a price to the price mentions list
func (s *ConversationAnalyticsService) AddPriceMention(sessionID uuid.UUID, price float64) error {
	analytics, err := s.db.ConversationAnalytics.
		Query().
		Where(conversationanalytics.SessionID(sessionID)).
		Only(s.ctx)

	if err != nil {
		return fmt.Errorf("failed to get analytics: %w", err)
	}

	prices := analytics.PriceMentions
	prices = append(prices, price)

	_, err = s.db.ConversationAnalytics.UpdateOne(analytics).
		SetPriceMentions(prices).
		Save(s.ctx)

	return err
}

// EndSession finalizes analytics for a session
func (s *ConversationAnalyticsService) EndSession(sessionID uuid.UUID) error {
	// Get session analytics
	analytics, err := s.db.ConversationAnalytics.
		Query().
		Where(conversationanalytics.SessionID(sessionID)).
		Only(s.ctx)

	if err != nil {
		return fmt.Errorf("failed to get analytics: %w", err)
	}

	// Calculate session duration
	duration := int(time.Since(analytics.SessionStartedAt).Seconds())

	// Get all messages from the session
	messages, err := s.db.Message.
		Query().
		Where(message.SessionID(sessionID)).
		Order(ent.Asc(message.FieldCreatedAt)).
		All(s.ctx)

	if err != nil {
		return fmt.Errorf("failed to get session messages: %w", err)
	}

	// Extract key topics using embeddings
	keyTopics := s.extractKeyTopics(messages)

	// Analyze sentiment
	sentiment, sentimentScore := s.analyzeSentiment(messages)

	// Detect primary intent
	intent := s.detectIntent(messages)

	// Calculate flow quality
	flowQuality := s.calculateFlowQuality(analytics, messages)

	// Extract preferences from conversation
	preferences := s.extractPreferences(messages)

	// Generate summary
	summary := s.generateSummary(messages, analytics)

	// Update analytics
	updater := s.db.ConversationAnalytics.UpdateOne(analytics).
		SetSessionDuration(duration).
		SetSessionEndedAt(time.Now()).
		SetSessionCompleted(true).
		SetKeyTopics(keyTopics).
		SetOverallSentiment(sentiment).
		SetSentimentScore(sentimentScore).
		SetPrimaryIntent(intent).
		SetFlowQualityScore(flowQuality).
		SetExtractedPreferences(preferences).
		SetSummary(summary)

	// Infer user satisfaction based on metrics
	if analytics.FoundProduct && analytics.ProductsClicked > 0 {
		updater.SetUserSatisfied(true)
	} else if analytics.SearchRefinementCount > 5 {
		updater.SetUserSatisfied(false)
	}

	_, err = updater.Save(s.ctx)
	if err != nil {
		return fmt.Errorf("failed to finalize analytics: %w", err)
	}

	fmt.Printf("📊 Finalized analytics for session %s: intent=%s, sentiment=%s, quality=%.2f\n",
		sessionID, intent, sentiment, flowQuality)

	return nil
}

// extractKeyTopics extracts main topics from conversation using embeddings
func (s *ConversationAnalyticsService) extractKeyTopics(messages []*ent.Message) []string {
	topics := make([]string, 0)
	seenTopics := make(map[string]bool)

	for _, msg := range messages {
		if msg.Role == "user" && len(msg.Content) > 10 {
			// Detect category as a topic
			category := s.embedding.DetectCategory(msg.Content)
			if category != "" && !seenTopics[category] {
				topics = append(topics, category)
				seenTopics[category] = true
			}

			// Extract common product keywords
			words := strings.Fields(strings.ToLower(msg.Content))
			for _, word := range words {
				if len(word) > 5 && !seenTopics[word] && s.isProductKeyword(word) {
					topics = append(topics, word)
					seenTopics[word] = true
					if len(topics) >= 10 {
						return topics
					}
				}
			}
		}
	}

	return topics
}

// isProductKeyword checks if a word is likely a product-related keyword
func (s *ConversationAnalyticsService) isProductKeyword(word string) bool {
	productWords := []string{"laptop", "phone", "computer", "gaming", "wireless", "bluetooth",
		"camera", "screen", "display", "processor", "memory", "storage", "battery", "charger"}

	for _, pw := range productWords {
		if strings.Contains(word, pw) {
			return true
		}
	}
	return false
}

// analyzeSentiment analyzes overall conversation sentiment
func (s *ConversationAnalyticsService) analyzeSentiment(messages []*ent.Message) (string, float64) {
	positiveWords := []string{"great", "perfect", "excellent", "love", "thanks", "good", "awesome"}
	negativeWords := []string{"bad", "expensive", "poor", "slow", "hate", "terrible", "awful"}

	positiveCount := 0
	negativeCount := 0

	for _, msg := range messages {
		if msg.Role == "user" {
			content := strings.ToLower(msg.Content)
			for _, word := range positiveWords {
				if strings.Contains(content, word) {
					positiveCount++
				}
			}
			for _, word := range negativeWords {
				if strings.Contains(content, word) {
					negativeCount++
				}
			}
		}
	}

	total := positiveCount + negativeCount
	if total == 0 {
		return "neutral", 0.0
	}

	score := float64(positiveCount-negativeCount) / float64(total)

	if score > 0.3 {
		return "positive", score
	} else if score < -0.3 {
		return "negative", score
	} else {
		return "neutral", score
	}
}

// detectIntent detects the primary intent of the conversation
func (s *ConversationAnalyticsService) detectIntent(messages []*ent.Message) string {
	explorationWords := []string{"looking for", "need", "want", "searching", "find"}
	purchaseWords := []string{"buy", "purchase", "order", "price", "cost", "cheap"}
	comparisonWords := []string{"compare", "difference", "better", "vs", "versus", "which"}
	informationWords := []string{"how", "what", "why", "explain", "tell me", "information"}

	scores := map[string]int{
		"exploration": 0,
		"purchase":    0,
		"comparison":  0,
		"information": 0,
	}

	for _, msg := range messages {
		if msg.Role == "user" {
			content := strings.ToLower(msg.Content)

			for _, word := range explorationWords {
				if strings.Contains(content, word) {
					scores["exploration"]++
				}
			}
			for _, word := range purchaseWords {
				if strings.Contains(content, word) {
					scores["purchase"]++
				}
			}
			for _, word := range comparisonWords {
				if strings.Contains(content, word) {
					scores["comparison"]++
				}
			}
			for _, word := range informationWords {
				if strings.Contains(content, word) {
					scores["information"]++
				}
			}
		}
	}

	// Find intent with highest score
	maxIntent := "exploration"
	maxScore := scores[maxIntent]

	for intent, score := range scores {
		if score > maxScore {
			maxIntent = intent
			maxScore = score
		}
	}

	return maxIntent
}

// calculateFlowQuality calculates conversation flow quality (0-1)
func (s *ConversationAnalyticsService) calculateFlowQuality(analytics *ent.ConversationAnalytics, messages []*ent.Message) float64 {
	quality := 1.0

	// Penalize for too many clarifications
	if analytics.ClarificationCount > 3 {
		quality -= 0.1 * float64(analytics.ClarificationCount-3)
	}

	// Penalize for too many search refinements
	if analytics.SearchRefinementCount > 3 {
		quality -= 0.1 * float64(analytics.SearchRefinementCount-3)
	}

	// Reward for product interactions
	if analytics.ProductsClicked > 0 {
		quality += 0.2
	}

	// Penalize for very short sessions
	if analytics.MessageCount < 3 {
		quality -= 0.3
	}

	// Clamp to 0-1 range
	if quality < 0 {
		quality = 0
	}
	if quality > 1 {
		quality = 1
	}

	return quality
}

// extractPreferences extracts user preferences from conversation
func (s *ConversationAnalyticsService) extractPreferences(messages []*ent.Message) map[string]interface{} {
	preferences := make(map[string]interface{})

	// This is a simplified version - in production, use NLP/LLM
	for _, msg := range messages {
		if msg.Role == "user" {
			content := strings.ToLower(msg.Content)

			// Extract size preferences
			if strings.Contains(content, "small") {
				preferences["size"] = "small"
			} else if strings.Contains(content, "large") {
				preferences["size"] = "large"
			}

			// Extract color preferences
			colors := []string{"black", "white", "red", "blue", "green", "silver", "gold"}
			for _, color := range colors {
				if strings.Contains(content, color) {
					preferences["color"] = color
					break
				}
			}

			// Extract budget preferences
			if strings.Contains(content, "budget") || strings.Contains(content, "cheap") {
				preferences["budget_conscious"] = true
			} else if strings.Contains(content, "premium") || strings.Contains(content, "best") {
				preferences["premium"] = true
			}
		}
	}

	return preferences
}

// generateSummary generates a brief summary of the conversation
func (s *ConversationAnalyticsService) generateSummary(messages []*ent.Message, analytics *ent.ConversationAnalytics) string {
	if len(messages) == 0 {
		return "Empty session"
	}

	// Build summary
	parts := []string{}

	if len(analytics.CategoriesExplored) > 0 {
		parts = append(parts, fmt.Sprintf("Explored: %s", strings.Join(analytics.CategoriesExplored, ", ")))
	}

	if analytics.ProductsShown > 0 {
		parts = append(parts, fmt.Sprintf("%d products shown", analytics.ProductsShown))
	}

	if analytics.ProductsClicked > 0 {
		parts = append(parts, fmt.Sprintf("%d clicked", analytics.ProductsClicked))
	}

	if len(parts) == 0 {
		return fmt.Sprintf("%d messages exchanged", analytics.MessageCount)
	}

	return strings.Join(parts, ". ")
}

// GetSessionAnalytics retrieves analytics for a session
func (s *ConversationAnalyticsService) GetSessionAnalytics(sessionID uuid.UUID) (*ent.ConversationAnalytics, error) {
	return s.db.ConversationAnalytics.
		Query().
		Where(conversationanalytics.SessionID(sessionID)).
		Only(s.ctx)
}

// GetUserAnalytics retrieves all analytics for a user
func (s *ConversationAnalyticsService) GetUserAnalytics(userID uuid.UUID, limit int) ([]*ent.ConversationAnalytics, error) {
	return s.db.ConversationAnalytics.
		Query().
		Where(conversationanalytics.UserID(userID)).
		Order(ent.Desc(conversationanalytics.FieldCreatedAt)).
		Limit(limit).
		All(s.ctx)
}
