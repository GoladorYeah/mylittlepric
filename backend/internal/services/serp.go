// backend/internal/services/serp.go
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

type SerpService struct {
	keyRotator *utils.KeyRotator
	config     *config.Config
}

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

func (s *SerpService) SearchProducts(query, searchType, country string) ([]models.ProductCard, int, error) {
	apiKey, keyIndex, err := s.keyRotator.GetNextKey()
	if err != nil {
		return nil, -1, fmt.Errorf("failed to get API key: %w", err)
	}

	// ✅ ЛОГИРУЕМ ОРИГИНАЛЬНЫЙ ЗАПРОС
	fmt.Printf("\n🔍 SERP API Request:\n")
	fmt.Printf("   Original Query: %s\n", query)
	fmt.Printf("   Type: %s\n", searchType)
	fmt.Printf("   Country: %s\n", country)
	fmt.Printf("   Language: %s\n", getLanguageForCountry(country))
	fmt.Printf("   Key Index: %d\n", keyIndex)

	parameter := map[string]string{
		"engine": "google_shopping",
		"q":      query,
		"gl":     country,
		"hl":     getLanguageForCountry(country),
	}

	search := g.NewGoogleSearch(parameter, apiKey)

	startTime := time.Now()
	data, err := search.GetJSON()
	elapsed := time.Since(startTime)

	if err != nil {
		fmt.Printf("   ❌ SERP API Error (%.2fs): %v\n", elapsed.Seconds(), err)
		return nil, keyIndex, fmt.Errorf("SERP API error: %w", err)
	}

	fmt.Printf("   ⏱️ Response time: %.2fs\n", elapsed.Seconds())

	shoppingItems := []types.ShoppingItem{}

	if shoppingResults, ok := data["shopping_results"].([]interface{}); ok {
		fmt.Printf("   📦 Raw results: %d products\n", len(shoppingResults))

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

func (s *SerpService) validateRelevance(query string, items []types.ShoppingItem, searchType string) SearchResult {
	if len(items) == 0 {
		return SearchResult{
			Products:        []types.ShoppingItem{},
			RelevanceScore:  0.0,
			IsRelevant:      false,
			AlternativeHint: "No products found",
		}
	}

	queryLower := strings.ToLower(query)
	queryWords := strings.Fields(queryLower)

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

	// Сортировка по убыванию score
	for i := 0; i < len(scoredProducts); i++ {
		for j := i + 1; j < len(scoredProducts); j++ {
			if scoredProducts[j].score > scoredProducts[i].score {
				scoredProducts[i], scoredProducts[j] = scoredProducts[j], scoredProducts[i]
			}
		}
	}

	// ✅ ВСТАВЬТЕ ЛОГ ЗДЕСЬ - ПОКАЗЫВАЕМ TOP-5 С ИХ SCORES
	fmt.Printf("   📊 Top 5 results:\n")
	topCount := min(5, len(scoredProducts))
	for i := 0; i < topCount; i++ {
		title := scoredProducts[i].item.Title
		if len(title) > 60 {
			title = title[:60] + "..."
		}
		fmt.Printf("      %d. [%.2f] %s\n", i+1, scoredProducts[i].score, title)
	}

	// ✅ СМЯГЧАЕМ THRESHOLDS - делаем поиск более гибким
	var threshold float32
	switch searchType {
	case "exact":
		threshold = 0.4 // было 0.7
	case "parameters":
		threshold = 0.2 // было 0.5
	case "category":
		threshold = 0.1 // было 0.3
	default:
		threshold = 0.2
	}

	relevantProducts := []types.ShoppingItem{}
	maxProducts := s.getMaxProducts(searchType)

	// ✅ ЕСЛИ НЕТ ПРОДУКТОВ С ДОСТАТОЧНЫМ SCORE, БЕРЕМ ЛУЧШИЕ
	if len(scoredProducts) > 0 {
		// Сначала берем все продукты выше threshold
		for i := 0; i < len(scoredProducts) && i < maxProducts; i++ {
			if scoredProducts[i].score >= threshold {
				relevantProducts = append(relevantProducts, scoredProducts[i].item)
			}
		}

		// Если ничего не нашли, берем хотя бы топ-3 результата
		if len(relevantProducts) == 0 && len(scoredProducts) > 0 {
			fmt.Printf("   💡 No products above threshold (%.2f), taking top results\n", threshold)
			topCount := min(3, len(scoredProducts))
			for i := 0; i < topCount; i++ {
				relevantProducts = append(relevantProducts, scoredProducts[i].item)
			}
		}
	}

	var avgScore float32
	if len(relevantProducts) > 0 {
		topCount := min(3, len(scoredProducts))
		for i := 0; i < topCount; i++ {
			avgScore += scoredProducts[i].score
		}
		avgScore /= float32(topCount)
	}

	// ✅ СМЯГЧАЕМ УСЛОВИЕ РЕЛЕВАНТНОСТИ
	// Теперь считаем релевантным если есть хоть какие-то продукты
	isRelevant := len(relevantProducts) > 0

	result := SearchResult{
		Products:       relevantProducts,
		RelevanceScore: avgScore,
		IsRelevant:     isRelevant,
	}

	if !isRelevant && len(items) > 0 {
		result.AlternativeHint = fmt.Sprintf("Found similar products but exact match not available. Best alternative: %s", items[0].Title)
	}

	return result
}
func (s *SerpService) calculateRelevanceScore(queryWords []string, item types.ShoppingItem) float32 {
	titleLower := strings.ToLower(item.Title)
	var score float32 = 0.0

	// ✅ 1. Полное совпадение всей фразы (бонус)
	queryStr := strings.Join(queryWords, " ")
	if strings.Contains(titleLower, queryStr) {
		score += 1.0
	}

	// ✅ 2. Все слова присутствуют (хороший сигнал)
	allWordsPresent := true
	for _, word := range queryWords {
		if len(word) <= 2 {
			continue // Игнорируем короткие слова
		}
		if !strings.Contains(titleLower, word) {
			allWordsPresent = false
			break
		}
	}
	if allWordsPresent {
		score += 0.6 // было 0.8
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
		if len(word) > 2 && !isCommonWord(word) {
			significantWords++
		}
	}

	if significantWords > 0 {
		matchRatio := float32(matchedWords) / float32(significantWords)
		score += matchRatio * 0.5 // До 0.5 баллов за частичное совпадение
	}

	// ✅ 4. Порядок слов (менее важно)
	if len(queryWords) >= 2 {
		titleWords := strings.Fields(titleLower)
		orderScore := s.calculateWordOrderScore(queryWords, titleWords)
		score += orderScore * 0.2 // было 0.3
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
				score += 0.3
				break
			}
		}
	}

	// ✅ 6. Номера моделей (если есть в запросе, должны совпадать)
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
			score += 0.3 // было 0.5
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

func extractModelNumbers(words []string) []string {
	modelNumbers := []string{}

	for _, word := range words {
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

func (s *SerpService) convertToProductCards(items []types.ShoppingItem, searchType string) []models.ProductCard {
	maxProducts := s.getMaxProducts(searchType)
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

func (s *SerpService) getMaxProducts(searchType string) int {
	switch searchType {
	case "exact":
		return 3
	case "parameters":
		return 6
	case "category":
		return 8
	default:
		return 6
	}
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

func (s *SerpService) SearchWithCache(query, searchType, country string, cacheService *CacheService) ([]models.ProductCard, int, error) {
	if cacheService != nil {
		cacheKey := fmt.Sprintf("search:%s:%s:%s", country, searchType, query)
		if cached, err := cacheService.GetSearchResults(cacheKey); err == nil && cached != nil {
			return cached, -1, nil
		}
	}

	cards, keyIndex, err := s.SearchProducts(query, searchType, country)
	if err != nil {
		return nil, keyIndex, err
	}

	if cacheService != nil {
		cacheKey := fmt.Sprintf("search:%s:%s:%s", country, searchType, query)
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
