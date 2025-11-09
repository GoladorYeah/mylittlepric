// backend/internal/services/serp.go
package services

import (
	"fmt"
	"strings"
	"time"

	g "github.com/serpapi/google-search-results-golang"

	"mylittleprice/internal/config"
	"mylittleprice/internal/domain"
	"mylittleprice/internal/models"
	"mylittleprice/internal/utils"
)

type SerpService struct {
	keyRotator *utils.KeyRotator
	config     *config.Config
}

type SearchResult struct {
	Products        []domain.ShoppingItem
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

func (s *SerpService) SearchProducts(query, searchType, country string, minPrice, maxPrice *float64) ([]models.ProductCard, int, error) {
	// Validate input
	if err := validateSearchQuery(query); err != nil {
		return nil, -1, fmt.Errorf("invalid search query: %w", err)
	}

	// Try up to total number of keys + 2 (for network retries)
	maxRetries := s.keyRotator.GetTotalKeys() + 1
	var lastErr error
	var lastKeyIndex int = -1
	var lastWasQuotaError bool = false

	for attempt := 0; attempt <= maxRetries; attempt++ {
		// Only apply backoff for network errors, not quota errors
		if attempt > 0 && !lastWasQuotaError {
			// Exponential backoff: 500ms, 1s, 2s
			backoffDuration := time.Duration(500*(1<<uint(attempt-1))) * time.Millisecond
			if backoffDuration > 2*time.Second {
				backoffDuration = 2 * time.Second
			}
			fmt.Printf("   ⏳ SERP retry attempt %d/%d after %v...\n", attempt+1, maxRetries+1, backoffDuration)
			time.Sleep(backoffDuration)
		} else if attempt > 0 && lastWasQuotaError {
			fmt.Printf("   🔄 Trying next key (attempt %d/%d)...\n", attempt+1, maxRetries+1)
		}

		apiKey, keyIndex, err := s.keyRotator.GetNextKey()
		if err != nil {
			return nil, -1, fmt.Errorf("failed to get API key: %w", err)
		}
		lastKeyIndex = keyIndex
		lastWasQuotaError = false

		// ✅ ЛОГИРУЕМ ОРИГИНАЛЬНЫЙ ЗАПРОС
		if attempt == 0 {
			fmt.Printf("\n🔍 SERP API Request:\n")
			fmt.Printf("   Original Query: %s\n", query)
			fmt.Printf("   Type: %s\n", searchType)
			fmt.Printf("   Country: %s\n", country)
			fmt.Printf("   Language: %s\n", getLanguageForCountry(country))
			if minPrice != nil || maxPrice != nil {
				fmt.Printf("   Price Range: %v - %v\n", minPrice, maxPrice)
			}
		}
		fmt.Printf("   Key Index: %d (attempt %d)\n", keyIndex, attempt+1)

		parameter := map[string]string{
			"engine": "google_shopping",
			"q":      query,
			"gl":     country,
			"hl":     getLanguageForCountry(country),
		}

		// Add price range filters if provided
		if minPrice != nil {
			parameter["min_price"] = fmt.Sprintf("%.0f", *minPrice)
		}
		if maxPrice != nil {
			parameter["max_price"] = fmt.Sprintf("%.0f", *maxPrice)
		}

		search := g.NewGoogleSearch(parameter, apiKey)

		startTime := time.Now()
		data, err := search.GetJSON()
		elapsed := time.Since(startTime)

		if err != nil {
			lastErr = err
			fmt.Printf("   ❌ SERP API Error (%.2fs, attempt %d/%d): %v\n", elapsed.Seconds(), attempt+1, maxRetries+1, err)

			// Check if error is retryable
			errMsg := err.Error()
			isQuotaError := strings.Contains(errMsg, "run out of searches") ||
				strings.Contains(errMsg, "quota exceeded") ||
				strings.Contains(errMsg, "limit exceeded") ||
				strings.Contains(errMsg, "rate limit")

			isNetworkError := strings.Contains(errMsg, "timeout") ||
				strings.Contains(errMsg, "503") ||
				strings.Contains(errMsg, "502") ||
				strings.Contains(errMsg, "500")

			if isQuotaError {
				// Mark this key as exhausted
				fmt.Printf("   ⚠️ Quota error detected for key %d\n", keyIndex)
				if markErr := s.keyRotator.MarkKeyAsExhausted(keyIndex); markErr != nil {
					fmt.Printf("   ⚠️ Failed to mark key as exhausted: %v\n", markErr)
				}
				// Try next key immediately (don't wait for backoff)
				lastWasQuotaError = true
				if attempt < maxRetries {
					continue
				}
			} else if isNetworkError {
				// Retryable network error - continue to next attempt
				if attempt < maxRetries {
					continue
				}
			}

			// Non-retryable error or last attempt
			return nil, keyIndex, fmt.Errorf("SERP API error: %w", err)
		}

		// Success!
		if attempt > 0 {
			fmt.Printf("   ✅ SERP request succeeded on attempt %d\n", attempt+1)
		}
		fmt.Printf("   ⏱️ Response time: %.2fs\n", elapsed.Seconds())

		shoppingItems := []domain.ShoppingItem{}

		if shoppingResults, ok := data["shopping_results"].([]interface{}); ok {
			fmt.Printf("   📦 Raw results: %d products\n", len(shoppingResults))

			for _, item := range shoppingResults {
				if itemMap, ok := item.(map[string]interface{}); ok {
					shoppingItem := domain.ShoppingItem{
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
						PageToken:   getStringFromInterface(itemMap["immersive_product_page_token"]),
					}
					shoppingItems = append(shoppingItems, shoppingItem)
				}
			}
		} else {
			fmt.Printf("   ⚠️ No shopping_results in response\n")
		}

		result := s.validateRelevance(query, shoppingItems, searchType)

		if !result.IsRelevant {
			fmt.Printf("   ⚠️ No relevant results for '%s' (score: %.2f)\n", query, result.RelevanceScore)
			return nil, keyIndex, fmt.Errorf("no relevant products found")
		}

		cards := s.convertToProductCards(result.Products, searchType)

		fmt.Printf("   ✅ Found %d relevant products (score: %.2f)\n\n", len(cards), result.RelevanceScore)

		return cards, keyIndex, nil
	}

	// All retries failed
	if lastErr != nil {
		return nil, lastKeyIndex, fmt.Errorf("SERP API failed after %d retries: %w", maxRetries+1, lastErr)
	}
	return nil, lastKeyIndex, fmt.Errorf("SERP API failed after %d retries", maxRetries+1)
}

func (s *SerpService) validateRelevance(query string, items []domain.ShoppingItem, searchType string) SearchResult {
	if len(items) == 0 {
		return SearchResult{
			Products:        []domain.ShoppingItem{},
			RelevanceScore:  0.0,
			IsRelevant:      false,
			AlternativeHint: "No products found",
		}
	}

	// ✅ НОВАЯ ЛОГИКА: Берем товары по позициям 1-10 от SerpAPI (уже отсортированы)
	// Вместо фильтрации по score, используем оригинальный порядок от Google Shopping

	maxProducts := 10 // Всегда берем до 10 товаров
	relevantProducts := []domain.ShoppingItem{}

	// Берем товары в оригинальном порядке (позиции 1-10)
	productCount := min(maxProducts, len(items))
	for i := 0; i < productCount; i++ {
		relevantProducts = append(relevantProducts, items[i])
	}

	// ✅ Логируем взятые товары
	fmt.Printf("   📊 Taking products at positions 1-%d (total available: %d):\n", productCount, len(items))
	for i := 0; i < min(5, productCount); i++ {
		title := items[i].Title
		if len(title) > 60 {
			title = title[:60] + "..."
		}
		fmt.Printf("      Position %d: %s\n", i+1, title)
	}
	if productCount > 5 {
		fmt.Printf("      ... and %d more\n", productCount-5)
	}

	// Считаем релевантным если есть хоть какие-то продукты
	isRelevant := len(relevantProducts) > 0

	result := SearchResult{
		Products:       relevantProducts,
		RelevanceScore: 1.0, // Не используется при новой логике
		IsRelevant:     isRelevant,
	}

	if !isRelevant && len(items) > 0 {
		result.AlternativeHint = fmt.Sprintf("Found similar products but exact match not available. Best alternative: %s", items[0].Title)
	}

	return result
}
func (s *SerpService) calculateRelevanceScore(queryWords []string, item domain.ShoppingItem) float32 {
	titleLower := strings.ToLower(item.Title)
	var score float32 = 0.0

	// ✅ 1. Полное совпадение всей фразы (бонус)
	queryStr := strings.Join(queryWords, " ")
	if strings.Contains(titleLower, queryStr) {
		score += float32(s.config.SerpScorePhraseMatch)
	}

	// ✅ 2. Все слова присутствуют (хороший сигнал)
	allWordsPresent := true
	for _, word := range queryWords {
		if len(word) <= s.config.SerpMinWordLength {
			continue // Игнорируем короткие слова
		}
		if !strings.Contains(titleLower, word) {
			allWordsPresent = false
			break
		}
	}
	if allWordsPresent {
		score += float32(s.config.SerpScoreAllWords)
	}

	// ✅ 3. Частичное совпадение слов (даже если не все слова есть)
	matchedWords := 0
	importantMatchedWords := 0
	for _, word := range queryWords {
		if len(word) <= 2 || isCommonWord(word) {
			continue
		}
		if strings.Contains(titleLower, word) {
			matchedWords++
			// Дополнительный бонус за важные слова (бренды, типы продуктов)
			if !isCommonWord(word) {
				importantMatchedWords++
			}
		}
	}

	significantWords := 0
	for _, word := range queryWords {
		if len(word) > s.config.SerpMinWordLength && !isCommonWord(word) {
			significantWords++
		}
	}

	if significantWords > 0 {
		matchRatio := float32(matchedWords) / float32(significantWords)
		score += matchRatio * float32(s.config.SerpScorePartialWords)
	}

	// ✅ 4. Порядок слов (менее важно)
	if len(queryWords) >= 2 {
		titleWords := strings.Fields(titleLower)
		orderScore := s.calculateWordOrderScore(queryWords, titleWords)
		score += orderScore * float32(s.config.SerpScoreWordOrderWeight)
	}

	// ✅ 5. Бренды (важно для точности)
	brands := []string{
		"apple", "iphone", "ipad", "macbook", "samsung", "galaxy",
		"google", "pixel", "xiaomi", "oneplus", "sony", "dell",
		"hp", "lenovo", "asus", "acer", "msi", "lg", "huawei",
		"nike", "adidas", "puma", "reebok", "under", "armour",
	}
	for _, brand := range brands {
		for _, word := range queryWords {
			if word == brand && strings.Contains(titleLower, brand) {
				score += float32(s.config.SerpScoreBrandMatch)
				break
			}
		}
	}

	// ✅ 6. Номера моделей (если есть в запросе, должны совпадать)
	modelNumbers := s.extractModelNumbers(queryWords)
	if len(modelNumbers) > 0 {
		hasModelMatch := false
		for _, modelNum := range modelNumbers {
			if strings.Contains(titleLower, modelNum) {
				hasModelMatch = true
				break
			}
		}
		if hasModelMatch {
			score += float32(s.config.SerpScoreModelMatch)
		}
		// Убираем штраф если модель не совпала - может быть похожий продукт
	}

	// ✅ 7. УБИРАЕМ ЖЕСТКИЙ ШТРАФ за лишние слова
	// Это слишком строго для гибкого поиска

	// Ограничиваем score в пределах [0, 1]
	if score < 0 {
		score = 0
	}
	if score > 1 {
		score = 1
	}

	return score
}

func (s *SerpService) calculateWordOrderScore(queryWords, titleWords []string) float32 {
	if len(queryWords) < 2 || len(titleWords) < 2 {
		return 0
	}

	matches := 0
	for i := 0; i < len(queryWords)-1; i++ {
		word1 := queryWords[i]
		word2 := queryWords[i+1]

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

		if pos1 != -1 && pos2 != -1 && pos1 < pos2 {
			matches++
		}
	}

	return float32(matches) / float32(len(queryWords)-1)
}

func (s *SerpService) extractModelNumbers(words []string) []string {
	modelNumbers := []string{}

	for _, word := range words {
		hasDigit := false
		for _, char := range word {
			if char >= '0' && char <= '9' {
				hasDigit = true
				break
			}
		}

		if hasDigit && len(word) >= s.config.SerpModelNumberMinLength {
			modelNumbers = append(modelNumbers, word)
		}
	}

	return modelNumbers
}

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

func (s *SerpService) GetProductDetailsByToken(pageToken string) (map[string]interface{}, int, error) {
	maxRetries := s.keyRotator.GetTotalKeys() + 1
	var lastErr error
	var lastKeyIndex int = -1
	var lastWasQuotaError bool = false

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 && !lastWasQuotaError {
			backoffDuration := time.Duration(500*(1<<uint(attempt-1))) * time.Millisecond
			if backoffDuration > 2*time.Second {
				backoffDuration = 2 * time.Second
			}
			fmt.Printf("   ⏳ Product details retry attempt %d/%d after %v...\n", attempt+1, maxRetries+1, backoffDuration)
			time.Sleep(backoffDuration)
		} else if attempt > 0 && lastWasQuotaError {
			fmt.Printf("   🔄 Trying next key for product details (attempt %d/%d)...\n", attempt+1, maxRetries+1)
		}

		apiKey, keyIndex, err := s.keyRotator.GetNextKey()
		if err != nil {
			return nil, -1, fmt.Errorf("failed to get API key: %w", err)
		}
		lastKeyIndex = keyIndex
		lastWasQuotaError = false

		parameter := map[string]string{
			"engine":      "google_immersive_product",
			"page_token":  pageToken,
			"more_stores": "true",
		}

		search := g.NewGoogleSearch(parameter, apiKey)
		startTime := time.Now()
		data, err := search.GetJSON()
		elapsed := time.Since(startTime)

		if err != nil {
			lastErr = err
			fmt.Printf("   ❌ Product details error (%.2fs, attempt %d/%d): %v\n", elapsed.Seconds(), attempt+1, maxRetries+1, err)

			errMsg := err.Error()
			isQuotaError := strings.Contains(errMsg, "run out of searches") ||
				strings.Contains(errMsg, "quota exceeded") ||
				strings.Contains(errMsg, "limit exceeded") ||
				strings.Contains(errMsg, "rate limit")

			isNetworkError := strings.Contains(errMsg, "timeout") ||
				strings.Contains(errMsg, "503") ||
				strings.Contains(errMsg, "502") ||
				strings.Contains(errMsg, "500")

			if isQuotaError {
				fmt.Printf("   ⚠️ Quota error detected for key %d\n", keyIndex)
				if markErr := s.keyRotator.MarkKeyAsExhausted(keyIndex); markErr != nil {
					fmt.Printf("   ⚠️ Failed to mark key as exhausted: %v\n", markErr)
				}
				lastWasQuotaError = true
				if attempt < maxRetries {
					continue
				}
			} else if isNetworkError {
				if attempt < maxRetries {
					continue
				}
			}

			return nil, keyIndex, fmt.Errorf("SERP API error: %w", err)
		}

		// Success!
		if attempt > 0 {
			fmt.Printf("   ✅ Product details request succeeded on attempt %d\n", attempt+1)
		}
		return data, keyIndex, nil
	}

	if lastErr != nil {
		return nil, lastKeyIndex, fmt.Errorf("product details failed after %d retries: %w", maxRetries+1, lastErr)
	}
	return nil, lastKeyIndex, fmt.Errorf("product details failed after %d retries", maxRetries+1)
}

func (s *SerpService) convertToProductCards(items []domain.ShoppingItem, searchType string) []models.ProductCard {
	// ✅ НОВАЯ ЛОГИКА: Всегда берем до 10 товаров независимо от типа поиска
	maxProducts := 10
	cards := make([]models.ProductCard, 0, maxProducts)

	for i, item := range items {
		if i >= maxProducts {
			break
		}

		pageToken := item.PageToken
		if pageToken == "" {
			pageToken = s.extractPageToken(item)
		}

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

func (s *SerpService) extractPageToken(item domain.ShoppingItem) string {
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

func findSubstring(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

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

func (s *SerpService) GetProductByPageToken(pageToken string) (map[string]interface{}, int, error) {
	return s.GetProductDetailsByToken(pageToken)
}

func (s *SerpService) SearchWithCache(query, searchType, country string, minPrice, maxPrice *float64, cacheService *CacheService) ([]models.ProductCard, int, error) {
	// Build cache key including price range
	cacheKey := fmt.Sprintf("search:%s:%s:%s", country, searchType, query)
	if minPrice != nil {
		cacheKey += fmt.Sprintf(":min%.0f", *minPrice)
	}
	if maxPrice != nil {
		cacheKey += fmt.Sprintf(":max%.0f", *maxPrice)
	}

	if cacheService != nil {
		if cached, err := cacheService.GetSearchResults(cacheKey); err == nil && cached != nil {
			return cached, -1, nil
		}
	}

	cards, keyIndex, err := s.SearchProducts(query, searchType, country, minPrice, maxPrice)
	if err != nil {
		return nil, keyIndex, err
	}

	if cacheService != nil {
		ttl := time.Duration(s.config.CacheSerpTTL) * time.Second
		_ = cacheService.SetSearchResults(cacheKey, cards, ttl)
	}

	return cards, keyIndex, nil
}

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
