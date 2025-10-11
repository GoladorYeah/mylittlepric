package services

import (
	"fmt"
	"strings"
	"time"

	g "github.com/serpapi/google-search-results-golang"

	"mylittleprice/internal/config"
	"mylittleprice/internal/models"
	"mylittleprice/internal/utils"
	"mylittleprice/pkg/types"
)

// SerpService handles all SERP API interactions with relevance validation
type SerpService struct {
	keyRotator *utils.KeyRotator
	config     *config.Config
}

// SearchResult contains search results with relevance score
type SearchResult struct {
	Products        []types.ShoppingItem
	RelevanceScore  float32
	IsRelevant      bool
	AlternativeHint string
}

func NewSerpService(keyRotator *utils.KeyRotator, cfg *config.Config) *SerpService {
	return &SerpService{
		keyRotator: keyRotator,
		config:     cfg,
	}
}

// SearchProducts searches for products with relevance validation
func (s *SerpService) SearchProducts(query, searchType, country string) ([]models.ProductCard, int, error) {
	// Get API key
	apiKey, keyIndex, err := s.keyRotator.GetNextKey()
	if err != nil {
		return nil, -1, fmt.Errorf("failed to get API key: %w", err)
	}

	// Build search parameters
	parameter := map[string]string{
		"engine": "google_shopping",
		"q":      query,
		"gl":     country,
		"hl":     getLanguageForCountry(country),
	}

	// Execute search
	search := g.NewGoogleSearch(parameter, apiKey)
	data, err := search.GetJSON()
	if err != nil {
		return nil, keyIndex, fmt.Errorf("SERP API error: %w", err)
	}

	// Extract shopping results
	shoppingItems := []types.ShoppingItem{}

	if shoppingResults, ok := data["shopping_results"].([]interface{}); ok {
		for _, item := range shoppingResults {
			if itemMap, ok := item.(map[string]interface{}); ok {
				shoppingItem := types.ShoppingItem{
					Position:    getIntFromInterface(itemMap["position"]),
					Title:       getStringFromInterface(itemMap["title"]),
					Link:        getStringFromInterface(itemMap["link"]),
					ProductLink: getStringFromInterface(itemMap["product_link"]),
					ProductID:   getStringFromInterface(itemMap["product_id"]),
					Thumbnail:   getStringFromInterface(itemMap["thumbnail"]),
					Price:       getStringFromInterface(itemMap["price"]),
					Merchant:    getStringFromInterface(itemMap["source"]),
					Rating:      getFloat32FromInterface(itemMap["rating"]),
					Reviews:     getIntFromInterface(itemMap["reviews"]),
					SerpAPILink: getStringFromInterface(itemMap["serpapi_product_api"]),
				}
				shoppingItems = append(shoppingItems, shoppingItem)
			}
		}
	}

	// 🎯 НОВАЯ ЛОГИКА: Валидация релевантности
	result := s.validateRelevance(query, shoppingItems, searchType)

	if !result.IsRelevant {
		fmt.Printf("   ⚠️ No relevant results for '%s' (score: %.2f)\n", query, result.RelevanceScore)
		return nil, keyIndex, fmt.Errorf("no relevant products found")
	}

	// Convert to ProductCards (теперь из types.ShoppingItem)
	cards := s.convertToProductCards(result.Products, searchType)

	fmt.Printf("   ✅ Found %d relevant products (score: %.2f)\n", len(cards), result.RelevanceScore)

	return cards, keyIndex, nil
}

// validateRelevance checks if search results are relevant to the query
func (s *SerpService) validateRelevance(query string, items []types.ShoppingItem, searchType string) SearchResult {
	if len(items) == 0 {
		return SearchResult{
			Products:        []types.ShoppingItem{},
			RelevanceScore:  0.0,
			IsRelevant:      false,
			AlternativeHint: "No products found",
		}
	}

	// Normalize query for comparison
	queryLower := strings.ToLower(query)
	queryWords := strings.Fields(queryLower)

	// Score each product
	type scoredProduct struct {
		item  types.ShoppingItem
		score float32
	}

	scoredProducts := []scoredProduct{}

	for _, item := range items {
		score := s.calculateRelevanceScore(queryWords, item)
		scoredProducts = append(scoredProducts, scoredProduct{
			item:  item,
			score: score,
		})
	}

	// Sort by score (highest first)
	for i := 0; i < len(scoredProducts); i++ {
		for j := i + 1; j < len(scoredProducts); j++ {
			if scoredProducts[j].score > scoredProducts[i].score {
				scoredProducts[i], scoredProducts[j] = scoredProducts[j], scoredProducts[i]
			}
		}
	}

	// Determine relevance threshold based on search type
	var threshold float32
	switch searchType {
	case "exact":
		threshold = 0.7 // Высокий порог для точного поиска
	case "parameters":
		threshold = 0.5 // Средний порог для параметрического поиска
	case "category":
		threshold = 0.3 // Низкий порог для категорийного поиска
	default:
		threshold = 0.5
	}

	// Get best products above threshold
	relevantProducts := []types.ShoppingItem{}
	maxProducts := s.getMaxProducts(searchType)

	for i := 0; i < len(scoredProducts) && i < maxProducts; i++ {
		if scoredProducts[i].score >= threshold {
			relevantProducts = append(relevantProducts, scoredProducts[i].item)
		}
	}

	// Calculate overall relevance
	var avgScore float32
	if len(relevantProducts) > 0 {
		for i := 0; i < len(relevantProducts) && i < 3; i++ { // Top 3 average
			avgScore += scoredProducts[i].score
		}
		avgScore /= float32(min(3, len(relevantProducts)))
	}

	isRelevant := len(relevantProducts) > 0 && avgScore >= threshold

	result := SearchResult{
		Products:       relevantProducts,
		RelevanceScore: avgScore,
		IsRelevant:     isRelevant,
	}

	// Add alternative hint if not relevant
	if !isRelevant && len(items) > 0 {
		result.AlternativeHint = fmt.Sprintf("Found similar products but exact match not available. Best alternative: %s", items[0].Title)
	}

	return result
}

// calculateRelevanceScore calculates relevance score for a product
func (s *SerpService) calculateRelevanceScore(queryWords []string, item types.ShoppingItem) float32 {
	titleLower := strings.ToLower(item.Title)
	var score float32 = 0.0

	// 1. Exact title match (weight: 1.0)
	queryStr := strings.Join(queryWords, " ")
	if strings.Contains(titleLower, queryStr) {
		score += 1.0
	}

	// 2. All words present (weight: 0.8)
	allWordsPresent := true
	for _, word := range queryWords {
		if !strings.Contains(titleLower, word) {
			allWordsPresent = false
			break
		}
	}
	if allWordsPresent {
		score += 0.8
	}

	// 3. Individual word matches (weight: 0.1 each)
	matchedWords := 0
	for _, word := range queryWords {
		// Skip common words
		if len(word) <= 2 || isCommonWord(word) {
			continue
		}
		if strings.Contains(titleLower, word) {
			matchedWords++
		}
	}
	if len(queryWords) > 0 {
		score += float32(matchedWords) / float32(len(queryWords)) * 0.5
	}

	// 4. Word order matters (weight: 0.3)
	if len(queryWords) >= 2 {
		titleWords := strings.Fields(titleLower)
		orderScore := s.calculateWordOrderScore(queryWords, titleWords)
		score += orderScore * 0.3
	}

	// 5. Brand name priority (weight: 0.4)
	brands := []string{
		"apple", "iphone", "ipad", "macbook", "samsung", "galaxy",
		"google", "pixel", "xiaomi", "oneplus", "sony", "dell",
		"hp", "lenovo", "asus", "acer", "msi", "lg", "huawei",
	}
	for _, brand := range brands {
		for _, word := range queryWords {
			if word == brand && strings.Contains(titleLower, brand) {
				score += 0.4
				break
			}
		}
	}

	// 6. Model number precision (weight: 0.5)
	// Check if specific model numbers match (e.g., "15", "S24", "M3")
	modelNumbers := extractModelNumbers(queryWords)
	if len(modelNumbers) > 0 {
		hasModelMatch := false
		for _, modelNum := range modelNumbers {
			if strings.Contains(titleLower, modelNum) {
				hasModelMatch = true
				break
			}
		}
		if hasModelMatch {
			score += 0.5
		} else {
			// Penalty for missing model number
			score -= 0.3
		}
	}

	// 7. Penalty for extra words not in query
	titleWords := strings.Fields(titleLower)
	extraWordsPenalty := float32(0)
	for _, titleWord := range titleWords {
		if len(titleWord) > 3 && !isCommonWord(titleWord) {
			found := false
			for _, queryWord := range queryWords {
				if titleWord == queryWord || strings.Contains(titleWord, queryWord) || strings.Contains(queryWord, titleWord) {
					found = true
					break
				}
			}
			if !found {
				extraWordsPenalty += 0.05
			}
		}
	}
	score -= extraWordsPenalty

	// Normalize score to 0-1 range
	if score < 0 {
		score = 0
	}
	if score > 1 {
		score = 1
	}

	return score
}

// calculateWordOrderScore checks if words appear in similar order
func (s *SerpService) calculateWordOrderScore(queryWords, titleWords []string) float32 {
	if len(queryWords) < 2 || len(titleWords) < 2 {
		return 0
	}

	matches := 0
	for i := 0; i < len(queryWords)-1; i++ {
		word1 := queryWords[i]
		word2 := queryWords[i+1]

		// Find positions in title
		pos1 := -1
		pos2 := -1

		for j, titleWord := range titleWords {
			if strings.Contains(titleWord, word1) {
				pos1 = j
			}
			if strings.Contains(titleWord, word2) {
				pos2 = j
			}
		}

		// Check if order is preserved
		if pos1 != -1 && pos2 != -1 && pos1 < pos2 {
			matches++
		}
	}

	return float32(matches) / float32(len(queryWords)-1)
}

// extractModelNumbers extracts model numbers from query words
func extractModelNumbers(words []string) []string {
	modelNumbers := []string{}

	for _, word := range words {
		// Check if word contains digits
		hasDigit := false
		for _, char := range word {
			if char >= '0' && char <= '9' {
				hasDigit = true
				break
			}
		}

		if hasDigit && len(word) >= 2 {
			modelNumbers = append(modelNumbers, word)
		}
	}

	return modelNumbers
}

// isCommonWord checks if word is common and should be ignored
func isCommonWord(word string) bool {
	commonWords := []string{
		"the", "a", "an", "and", "or", "but", "in", "on", "at", "to", "for",
		"of", "with", "by", "from", "as", "is", "was", "are", "were", "be",
		"have", "has", "had", "do", "does", "did", "will", "would", "could",
		"should", "may", "might", "can", "new", "latest", "best", "pro", "air",
		"version", "model", "series", "generation", "gen",
	}

	for _, common := range commonWords {
		if word == common {
			return true
		}
	}

	return false
}

// GetProductDetailsByToken fetches detailed product info using page_token
func (s *SerpService) GetProductDetailsByToken(pageToken string) (map[string]interface{}, int, error) {
	apiKey, keyIndex, err := s.keyRotator.GetNextKey()
	if err != nil {
		return nil, -1, fmt.Errorf("failed to get API key: %w", err)
	}

	parameter := map[string]string{
		"engine":      "google_immersive_product",
		"page_token":  pageToken,
		"more_stores": "true",
	}

	search := g.NewGoogleSearch(parameter, apiKey)
	data, err := search.GetJSON()
	if err != nil {
		return nil, keyIndex, fmt.Errorf("SERP API error: %w", err)
	}

	return data, keyIndex, nil
}

// convertToProductCards converts ShoppingItems to ProductCards
func (s *SerpService) convertToProductCards(items []types.ShoppingItem, searchType string) []models.ProductCard {
	maxProducts := s.getMaxProducts(searchType)
	cards := make([]models.ProductCard, 0, maxProducts)

	for i, item := range items {
		if i >= maxProducts {
			break
		}

		pageToken := s.extractPageToken(item)

		badge := ""
		if item.Rating > 0 {
			badge = fmt.Sprintf("⭐ %.1f", item.Rating)
		}

		card := models.ProductCard{
			Name:        item.Title,
			Price:       item.Price,
			OldPrice:    item.OldPrice,
			Link:        item.ProductLink,
			Image:       item.Thumbnail,
			Description: item.Merchant,
			Badge:       badge,
			PageToken:   pageToken,
		}

		cards = append(cards, card)
	}

	return cards
}

// extractPageToken extracts page_token from shopping item
func (s *SerpService) extractPageToken(item types.ShoppingItem) string {
	if item.PageToken != "" {
		return item.PageToken
	}

	if item.SerpAPILink != "" {
		return extractTokenFromSerpAPILink(item.SerpAPILink)
	}

	if item.ProductID != "" {
		return item.ProductID
	}

	return ""
}

// getMaxProducts returns max number of products based on search type
func (s *SerpService) getMaxProducts(searchType string) int {
	switch searchType {
	case "exact":
		return 3 // Показываем топ-3 для точного поиска
	case "parameters":
		return 6
	case "category":
		return 8
	default:
		return 6
	}
}

// extractTokenFromSerpAPILink extracts page_token from serpapi URL
func extractTokenFromSerpAPILink(link string) string {
	tokenStart := findSubstring(link, "page_token=")
	if tokenStart == -1 {
		return ""
	}

	tokenStart += len("page_token=")
	tokenEnd := findSubstring(link[tokenStart:], "&")
	if tokenEnd == -1 {
		return link[tokenStart:]
	}

	return link[tokenStart : tokenStart+tokenEnd]
}

// findSubstring finds substring position
func findSubstring(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// getLanguageForCountry maps country code to language
func getLanguageForCountry(country string) string {
	languageMap := map[string]string{
		"CH": "de", "DE": "de", "AT": "de",
		"FR": "fr", "IT": "it", "ES": "es",
		"PT": "pt", "NL": "nl", "BE": "nl",
		"PL": "pl", "CZ": "cs", "SE": "sv",
		"NO": "no", "DK": "da", "FI": "fi",
		"GB": "en", "US": "en",
	}

	if lang, ok := languageMap[country]; ok {
		return lang
	}
	return "en"
}

// GetProductByPageToken is an alias
func (s *SerpService) GetProductByPageToken(pageToken string) (map[string]interface{}, int, error) {
	return s.GetProductDetailsByToken(pageToken)
}

// SearchWithCache searches products with optional caching
func (s *SerpService) SearchWithCache(query, searchType, country string, cacheService *CacheService) ([]models.ProductCard, int, error) {
	// Try cache
	if cacheService != nil {
		cacheKey := fmt.Sprintf("search:%s:%s:%s", country, searchType, query)
		if cached, err := cacheService.GetSearchResults(cacheKey); err == nil && cached != nil {
			return cached, -1, nil
		}
	}

	// Cache miss - fetch from API
	cards, keyIndex, err := s.SearchProducts(query, searchType, country)
	if err != nil {
		return nil, keyIndex, err
	}

	// Store in cache
	if cacheService != nil {
		cacheKey := fmt.Sprintf("search:%s:%s:%s", country, searchType, query)
		ttl := time.Duration(s.config.CacheSerpTTL) * time.Second
		_ = cacheService.SetSearchResults(cacheKey, cards, ttl)
	}

	return cards, keyIndex, nil
}

// Helper functions
func getStringFromInterface(val interface{}) string {
	if str, ok := val.(string); ok {
		return str
	}
	return ""
}

func getIntFromInterface(val interface{}) int {
	switch v := val.(type) {
	case int:
		return v
	case int64:
		return int(v)
	case float64:
		return int(v)
	default:
		return 0
	}
}

func getFloat32FromInterface(val interface{}) float32 {
	switch v := val.(type) {
	case float32:
		return v
	case float64:
		return float32(v)
	case int:
		return float32(v)
	default:
		return 0
	}
}
