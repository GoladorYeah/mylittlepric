// backend/internal/services/user_behavior.go
package services

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"mylittleprice/ent"
	"mylittleprice/ent/conversationanalytics"
	"mylittleprice/ent/productinteraction"
	"mylittleprice/ent/userbehaviorprofile"
)

// UserBehaviorService handles learning from user interactions to improve recommendations
type UserBehaviorService struct {
	db    *ent.Client
	redis *redis.Client
	ctx   context.Context
}

// NewUserBehaviorService creates a new user behavior service
func NewUserBehaviorService(db *ent.Client, redis *redis.Client) *UserBehaviorService {
	return &UserBehaviorService{
		db:    db,
		redis: redis,
		ctx:   context.Background(),
	}
}

// GetOrCreateProfile gets or creates a behavior profile for a user
func (s *UserBehaviorService) GetOrCreateProfile(userID uuid.UUID) (*ent.UserBehaviorProfile, error) {
	// Try to get existing profile
	profile, err := s.db.UserBehaviorProfile.
		Query().
		Where(userbehaviorprofile.UserID(userID)).
		Only(s.ctx)

	if err == nil {
		return profile, nil
	}

	// Profile doesn't exist, create it
	profile, err = s.db.UserBehaviorProfile.
		Create().
		SetUserID(userID).
		Save(s.ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to create user behavior profile: %w", err)
	}

	fmt.Printf("📊 Created new behavior profile for user %s\n", userID)
	return profile, nil
}

// LearnFromSession analyzes a completed session and updates the user's behavior profile
func (s *UserBehaviorService) LearnFromSession(userID uuid.UUID, sessionID uuid.UUID) error {
	// Get the profile
	profile, err := s.GetOrCreateProfile(userID)
	if err != nil {
		return err
	}

	// Get session analytics
	analytics, err := s.db.ConversationAnalytics.
		Query().
		Where(conversationanalytics.SessionID(sessionID)).
		Only(s.ctx)

	if err != nil {
		return fmt.Errorf("failed to get session analytics: %w", err)
	}

	// Get product interactions from this session
	interactions, err := s.db.ProductInteraction.
		Query().
		Where(productinteraction.UserID(userID)).
		Where(productinteraction.SessionID(sessionID.String())).
		All(s.ctx)

	if err != nil {
		return fmt.Errorf("failed to get product interactions: %w", err)
	}

	// Update profile based on learnings
	updater := s.db.UserBehaviorProfile.UpdateOne(profile)

	// Update category preferences based on categories explored
	categoryPrefs := profile.CategoryPreferences
	for _, category := range analytics.CategoriesExplored {
		// Increase weight for explored categories
		currentWeight := categoryPrefs[category]
		newWeight := math.Min(1.0, currentWeight+0.1) // Increment by 0.1, max 1.0
		categoryPrefs[category] = newWeight
	}
	updater.SetCategoryPreferences(categoryPrefs)

	// Update brand preferences based on product interactions
	brandPrefs := profile.BrandPreferences
	for _, interaction := range interactions {
		if interaction.ProductBrand != nil && *interaction.ProductBrand != "" {
			brand := *interaction.ProductBrand
			// Increase count for interacted brands
			weight := 1
			if interaction.InteractionType == "clicked" {
				weight = 3 // Clicks are stronger signals
			} else if interaction.InteractionType == "compared" {
				weight = 2
			}
			brandPrefs[brand] = brandPrefs[brand] + weight
		}
	}
	updater.SetBrandPreferences(brandPrefs)

	// Update price ranges based on viewed products
	priceRanges := profile.PriceRanges
	for _, interaction := range interactions {
		if interaction.ProductPrice != nil && interaction.ProductCategory != nil {
			category := *interaction.ProductCategory
			price := *interaction.ProductPrice

			// Get existing range or create new
			rangeData, exists := priceRanges[category]
			if !exists {
				priceRanges[category] = map[string]interface{}{
					"min": price,
					"max": price,
				}
			} else {
				rangeMap := rangeData.(map[string]interface{})
				minPrice := rangeMap["min"].(float64)
				maxPrice := rangeMap["max"].(float64)

				// Update range
				if price < minPrice {
					rangeMap["min"] = price
				}
				if price > maxPrice {
					rangeMap["max"] = price
				}
				priceRanges[category] = rangeMap
			}
		}
	}
	updater.SetPriceRanges(priceRanges)

	// Infer communication style from message patterns
	if analytics.MessageCount > 0 {
		avgWordsPerMessage := float64(len(strings.Fields(analytics.Summary))) / float64(analytics.MessageCount)
		if avgWordsPerMessage < 10 {
			updater.SetCommunicationStyle("brief")
		} else if avgWordsPerMessage > 30 {
			updater.SetCommunicationStyle("detailed")
		} else {
			updater.SetCommunicationStyle("balanced")
		}
	}

	// Update session metrics
	totalSessions := profile.TotalSessions + 1
	newAvgDuration := (profile.AvgSessionDuration*float64(profile.TotalSessions) + float64(analytics.SessionDuration)) / float64(totalSessions)
	newAvgMessages := (profile.AvgMessagesPerSession*float64(profile.TotalSessions) + float64(analytics.MessageCount)) / float64(totalSessions)

	updater.SetTotalSessions(totalSessions)
	updater.SetAvgSessionDuration(newAvgDuration)
	updater.SetAvgMessagesPerSession(newAvgMessages)

	// Update product counters
	updater.SetTotalProductsViewed(profile.TotalProductsViewed + analytics.ProductsShown)
	updater.SetTotalProductsClicked(profile.TotalProductsClicked + analytics.ProductsClicked)

	// Calculate success rate
	successfulSessions := profile.TotalSessions
	if analytics.FoundProduct {
		successfulSessions++
	}
	newSuccessRate := float64(successfulSessions) / float64(totalSessions)
	updater.SetSuccessRate(newSuccessRate)

	// Extract common keywords from this session
	keywords := s.extractKeywordsFromTopics(analytics.KeyTopics)
	existingKeywords := profile.CommonKeywords
	updatedKeywords := s.mergeKeywords(existingKeywords, keywords)
	updater.SetCommonKeywords(updatedKeywords)

	// Update time patterns
	timePatterns := profile.TimePatterns
	hour := analytics.SessionStartedAt.Hour()
	hourKey := fmt.Sprintf("hour_%d", hour)
	count, _ := timePatterns[hourKey].(int)
	timePatterns[hourKey] = count + 1
	updater.SetTimePatterns(timePatterns)

	// Update last learning timestamp
	updater.SetLastLearningUpdate(time.Now())

	// Save updates
	_, err = updater.Save(s.ctx)
	if err != nil {
		return fmt.Errorf("failed to update user behavior profile: %w", err)
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf("user:behavior:%s", userID)
	s.redis.Del(s.ctx, cacheKey)

	fmt.Printf("🧠 Learned from session %s for user %s\n", sessionID, userID)
	return nil
}

// GetRecommendedCategories returns top N recommended categories for a user
func (s *UserBehaviorService) GetRecommendedCategories(userID uuid.UUID, limit int) ([]string, error) {
	profile, err := s.GetOrCreateProfile(userID)
	if err != nil {
		return nil, err
	}

	// Sort categories by weight
	type categoryWeight struct {
		Category string
		Weight   float64
	}

	var categories []categoryWeight
	for cat, weight := range profile.CategoryPreferences {
		categories = append(categories, categoryWeight{Category: cat, Weight: weight})
	}

	sort.Slice(categories, func(i, j int) bool {
		return categories[i].Weight > categories[j].Weight
	})

	// Return top N
	result := make([]string, 0, limit)
	for i := 0; i < limit && i < len(categories); i++ {
		result = append(result, categories[i].Category)
	}

	return result, nil
}

// GetPreferredBrands returns top N preferred brands for a user
func (s *UserBehaviorService) GetPreferredBrands(userID uuid.UUID, limit int) ([]string, error) {
	profile, err := s.GetOrCreateProfile(userID)
	if err != nil {
		return nil, err
	}

	// Sort brands by frequency
	type brandFreq struct {
		Brand string
		Count int
	}

	var brands []brandFreq
	for brand, count := range profile.BrandPreferences {
		brands = append(brands, brandFreq{Brand: brand, Count: count})
	}

	sort.Slice(brands, func(i, j int) bool {
		return brands[i].Count > brands[j].Count
	})

	// Return top N
	result := make([]string, 0, limit)
	for i := 0; i < limit && i < len(brands); i++ {
		result = append(result, brands[i].Brand)
	}

	return result, nil
}

// GetPriceRange returns the typical price range for a user in a category
func (s *UserBehaviorService) GetPriceRange(userID uuid.UUID, category string) (min float64, max float64, err error) {
	profile, err := s.GetOrCreateProfile(userID)
	if err != nil {
		return 0, 0, err
	}

	rangeData, exists := profile.PriceRanges[category]
	if !exists {
		return 0, 0, fmt.Errorf("no price range data for category %s", category)
	}

	rangeMap := rangeData.(map[string]interface{})
	min = rangeMap["min"].(float64)
	max = rangeMap["max"].(float64)

	return min, max, nil
}

// GetProfileSummary returns a human-readable summary of the user's behavior profile
func (s *UserBehaviorService) GetProfileSummary(userID uuid.UUID) (map[string]interface{}, error) {
	profile, err := s.GetOrCreateProfile(userID)
	if err != nil {
		return nil, err
	}

	// Get top categories
	topCategories, _ := s.GetRecommendedCategories(userID, 5)

	// Get top brands
	topBrands, _ := s.GetPreferredBrands(userID, 5)

	summary := map[string]interface{}{
		"communication_style":       profile.CommunicationStyle,
		"total_sessions":            profile.TotalSessions,
		"avg_session_duration_mins": profile.AvgSessionDuration / 60,
		"avg_messages_per_session":  profile.AvgMessagesPerSession,
		"success_rate":              profile.SuccessRate * 100, // As percentage
		"total_products_viewed":     profile.TotalProductsViewed,
		"total_products_clicked":    profile.TotalProductsClicked,
		"top_categories":            topCategories,
		"top_brands":                topBrands,
		"common_keywords":           profile.CommonKeywords,
		"last_updated":              profile.LastLearningUpdate,
	}

	return summary, nil
}

// GetCachedProfile retrieves profile from cache or database
func (s *UserBehaviorService) GetCachedProfile(userID uuid.UUID) (*ent.UserBehaviorProfile, error) {
	cacheKey := fmt.Sprintf("user:behavior:%s", userID)

	// Try cache first
	cached, err := s.redis.Get(s.ctx, cacheKey).Bytes()
	if err == nil {
		var profile ent.UserBehaviorProfile
		if err := json.Unmarshal(cached, &profile); err == nil {
			return &profile, nil
		}
	}

	// Get from database
	profile, err := s.GetOrCreateProfile(userID)
	if err != nil {
		return nil, err
	}

	// Cache for 1 hour
	jsonData, _ := json.Marshal(profile)
	s.redis.Set(s.ctx, cacheKey, jsonData, time.Hour)

	return profile, nil
}

// extractKeywordsFromTopics extracts individual keywords from topic strings
func (s *UserBehaviorService) extractKeywordsFromTopics(topics []string) []string {
	keywords := make([]string, 0)
	for _, topic := range topics {
		words := strings.Fields(strings.ToLower(topic))
		keywords = append(keywords, words...)
	}
	return keywords
}

// mergeKeywords combines existing keywords with new ones, keeping most frequent
func (s *UserBehaviorService) mergeKeywords(existing, new []string) []string {
	// Count frequency
	freq := make(map[string]int)
	for _, kw := range existing {
		freq[kw]++
	}
	for _, kw := range new {
		freq[kw]++
	}

	// Sort by frequency
	type kwFreq struct {
		Keyword string
		Count   int
	}
	var keywords []kwFreq
	for kw, count := range freq {
		keywords = append(keywords, kwFreq{Keyword: kw, Count: count})
	}

	sort.Slice(keywords, func(i, j int) bool {
		return keywords[i].Count > keywords[j].Count
	})

	// Return top 50 keywords
	result := make([]string, 0, 50)
	for i := 0; i < 50 && i < len(keywords); i++ {
		result = append(result, keywords[i].Keyword)
	}

	return result
}
