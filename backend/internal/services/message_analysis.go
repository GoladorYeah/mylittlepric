// backend/internal/services/message_analysis.go
package services

import (
	"regexp"
	"strconv"
	"strings"
)

// MessageAnalysisService provides advanced message analysis for intent detection
type MessageAnalysisService struct {
	embedding *EmbeddingService
}

// NewMessageAnalysisService creates a new message analysis service
func NewMessageAnalysisService(embedding *EmbeddingService) *MessageAnalysisService {
	return &MessageAnalysisService{
		embedding: embedding,
	}
}

// IntentResult represents detected intent and confidence
type IntentResult struct {
	Intent     string  // "exploration", "purchase", "comparison", "information", "clarification"
	Confidence float64 // 0.0 - 1.0
	Context    map[string]interface{}
}

// DetectIntent analyzes a message and detects user intent
func (s *MessageAnalysisService) DetectIntent(message string) *IntentResult {
	message = strings.ToLower(message)

	// Score each intent type
	scores := map[string]float64{
		"exploration":    s.scoreExplorationIntent(message),
		"purchase":       s.scorePurchaseIntent(message),
		"comparison":     s.scoreComparisonIntent(message),
		"information":    s.scoreInformationIntent(message),
		"clarification":  s.scoreClarificationIntent(message),
		"specification":  s.scoreSpecificationIntent(message),
		"budget_inquiry": s.scoreBudgetIntent(message),
	}

	// Find highest scoring intent
	maxIntent := "exploration"
	maxScore := scores[maxIntent]

	for intent, score := range scores {
		if score > maxScore {
			maxIntent = intent
			maxScore = score
		}
	}

	// Build context based on detected intent
	context := s.buildIntentContext(message, maxIntent)

	return &IntentResult{
		Intent:     maxIntent,
		Confidence: maxScore,
		Context:    context,
	}
}

// scoreExplorationIntent scores exploration intent
func (s *MessageAnalysisService) scoreExplorationIntent(message string) float64 {
	keywords := []string{
		"looking for", "need", "want", "searching", "find", "show me",
		"recommend", "suggest", "help me find", "what are", "any",
	}
	return s.scoreKeywords(message, keywords, 0.15)
}

// scorePurchaseIntent scores purchase intent
func (s *MessageAnalysisService) scorePurchaseIntent(message string) float64 {
	keywords := []string{
		"buy", "purchase", "order", "get", "shop", "where to buy",
		"best deal", "cheapest", "price", "cost", "affordable",
	}
	return s.scoreKeywords(message, keywords, 0.2)
}

// scoreComparisonIntent scores comparison intent
func (s *MessageAnalysisService) scoreComparisonIntent(message string) float64 {
	keywords := []string{
		"compare", "comparison", "difference", "better", "vs", "versus",
		"which", "between", "or", "choice", "decide", "pros and cons",
	}
	return s.scoreKeywords(message, keywords, 0.25)
}

// scoreInformationIntent scores information seeking intent
func (s *MessageAnalysisService) scoreInformationIntent(message string) float64 {
	keywords := []string{
		"how", "what", "why", "when", "explain", "tell me",
		"information", "details", "learn", "know", "understand",
	}
	return s.scoreKeywords(message, keywords, 0.15)
}

// scoreClarificationIntent scores clarification intent
func (s *MessageAnalysisService) scoreClarificationIntent(message string) float64 {
	keywords := []string{
		"yes", "no", "ok", "sure", "maybe", "i mean", "actually",
		"sorry", "correction", "i said", "right", "correct",
	}
	return s.scoreKeywords(message, keywords, 0.3)
}

// scoreSpecificationIntent scores specification inquiry intent
func (s *MessageAnalysisService) scoreSpecificationIntent(message string) float64 {
	keywords := []string{
		"specs", "specifications", "features", "performance",
		"capacity", "size", "weight", "dimensions", "technical",
		"processor", "ram", "storage", "battery", "screen",
	}
	return s.scoreKeywords(message, keywords, 0.2)
}

// scoreBudgetIntent scores budget inquiry intent
func (s *MessageAnalysisService) scoreBudgetIntent(message string) float64 {
	keywords := []string{
		"budget", "under", "below", "less than", "maximum",
		"afford", "cheap", "expensive", "cost", "dollars",
	}
	score := s.scoreKeywords(message, keywords, 0.2)

	// Boost if numbers are mentioned
	if s.containsPriceNumbers(message) {
		score += 0.3
	}

	return score
}

// scoreKeywords calculates score based on keyword matches
func (s *MessageAnalysisService) scoreKeywords(message string, keywords []string, baseScore float64) float64 {
	score := 0.0
	for _, keyword := range keywords {
		if strings.Contains(message, keyword) {
			score += baseScore
		}
	}
	// Clamp to 1.0
	if score > 1.0 {
		score = 1.0
	}
	return score
}

// buildIntentContext builds context information for the detected intent
func (s *MessageAnalysisService) buildIntentContext(message string, intent string) map[string]interface{} {
	context := make(map[string]interface{})

	// Extract price mentions
	prices := s.ExtractPrices(message)
	if len(prices) > 0 {
		context["prices"] = prices
	}

	// Extract brands
	brands := s.ExtractBrands(message)
	if len(brands) > 0 {
		context["brands"] = brands
	}

	// Extract categories using embeddings
	category := s.embedding.DetectCategory(message)
	if category != "" {
		context["category"] = category
	}

	// Intent-specific context
	switch intent {
	case "comparison":
		context["comparison_items"] = s.extractComparisonItems(message)
	case "specification":
		context["specs_interested"] = s.extractSpecifications(message)
	case "budget_inquiry":
		if len(prices) > 0 {
			context["budget_limit"] = prices[0]
		}
	}

	return context
}

// ExtractPrices extracts price values from message
func (s *MessageAnalysisService) ExtractPrices(message string) []float64 {
	prices := make([]float64, 0)

	// Patterns for price extraction
	patterns := []string{
		`\$\s*(\d+(?:,\d{3})*(?:\.\d{2})?)`,      // $1,000.00
		`(\d+(?:,\d{3})*(?:\.\d{2})?)\s*dollars`, // 1000 dollars
		`(\d+(?:,\d{3})*(?:\.\d{2})?)\s*USD`,     // 1000 USD
		`under\s+(\d+)`,                           // under 500
		`below\s+(\d+)`,                           // below 500
		`less\s+than\s+(\d+)`,                     // less than 500
		`around\s+(\d+)`,                          // around 500
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindAllStringSubmatch(message, -1)
		for _, match := range matches {
			if len(match) > 1 {
				// Remove commas and parse
				numStr := strings.ReplaceAll(match[1], ",", "")
				if price, err := strconv.ParseFloat(numStr, 64); err == nil {
					prices = append(prices, price)
				}
			}
		}
	}

	return prices
}

// ExtractBrands extracts brand mentions from message
func (s *MessageAnalysisService) ExtractBrands(message string) []string {
	brands := make([]string, 0)
	message = strings.ToLower(message)

	// Common brands (this should be expanded)
	knownBrands := []string{
		"apple", "samsung", "sony", "lg", "dell", "hp", "lenovo",
		"asus", "acer", "microsoft", "google", "amazon", "nike",
		"adidas", "canon", "nikon", "bose", "jbl", "logitech",
	}

	for _, brand := range knownBrands {
		if strings.Contains(message, brand) {
			brands = append(brands, brand)
		}
	}

	return brands
}

// containsPriceNumbers checks if message contains numeric values
func (s *MessageAnalysisService) containsPriceNumbers(message string) bool {
	re := regexp.MustCompile(`\$?\d+`)
	return re.MatchString(message)
}

// extractComparisonItems extracts items being compared
func (s *MessageAnalysisService) extractComparisonItems(message string) []string {
	items := make([]string, 0)

	// Look for " vs ", " or ", " versus " patterns
	patterns := []string{`(.+?)\s+vs\s+(.+)`, `(.+?)\s+or\s+(.+)`, `(.+?)\s+versus\s+(.+)`}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(message)
		if len(matches) > 2 {
			items = append(items, strings.TrimSpace(matches[1]))
			items = append(items, strings.TrimSpace(matches[2]))
			break
		}
	}

	return items
}

// extractSpecifications extracts technical specifications mentioned
func (s *MessageAnalysisService) extractSpecifications(message string) []string {
	specs := make([]string, 0)

	specKeywords := []string{
		"processor", "cpu", "ram", "memory", "storage", "ssd", "hdd",
		"screen", "display", "resolution", "battery", "camera",
		"weight", "size", "dimensions", "port", "connectivity",
	}

	for _, spec := range specKeywords {
		if strings.Contains(message, spec) {
			specs = append(specs, spec)
		}
	}

	return specs
}

// AnalyzeMessageSentiment analyzes sentiment of a message
func (s *MessageAnalysisService) AnalyzeMessageSentiment(message string) (sentiment string, score float64) {
	message = strings.ToLower(message)

	positiveWords := []string{
		"great", "excellent", "perfect", "love", "amazing", "awesome",
		"good", "nice", "thanks", "thank you", "appreciate", "helpful",
	}

	negativeWords := []string{
		"bad", "terrible", "awful", "hate", "poor", "expensive",
		"slow", "disappointed", "frustrating", "useless", "wrong",
	}

	positiveCount := 0
	negativeCount := 0

	for _, word := range positiveWords {
		if strings.Contains(message, word) {
			positiveCount++
		}
	}

	for _, word := range negativeWords {
		if strings.Contains(message, word) {
			negativeCount++
		}
	}

	total := positiveCount + negativeCount
	if total == 0 {
		return "neutral", 0.0
	}

	score = float64(positiveCount-negativeCount) / float64(total)

	if score > 0.3 {
		return "positive", score
	} else if score < -0.3 {
		return "negative", score
	}

	return "neutral", score
}

// DetectCommunicationStyle detects user's preferred communication style
func (s *MessageAnalysisService) DetectCommunicationStyle(messages []string) string {
	if len(messages) == 0 {
		return "balanced"
	}

	totalWords := 0
	totalMessages := len(messages)

	for _, msg := range messages {
		words := strings.Fields(msg)
		totalWords += len(words)
	}

	avgWordsPerMessage := float64(totalWords) / float64(totalMessages)

	if avgWordsPerMessage < 8 {
		return "brief"
	} else if avgWordsPerMessage > 25 {
		return "detailed"
	}

	return "balanced"
}

// IsQuestionNeedingClarification checks if message is a clarification question
func (s *MessageAnalysisService) IsQuestionNeedingClarification(message string) bool {
	message = strings.ToLower(message)

	clarificationPhrases := []string{
		"what do you mean",
		"can you explain",
		"i don't understand",
		"could you clarify",
		"not sure",
		"confused",
		"what's the difference",
	}

	for _, phrase := range clarificationPhrases {
		if strings.Contains(message, phrase) {
			return true
		}
	}

	return false
}

// ExtractProductRequirements extracts specific product requirements from message
func (s *MessageAnalysisService) ExtractProductRequirements(message string) map[string]interface{} {
	requirements := make(map[string]interface{})

	message = strings.ToLower(message)

	// Size requirements
	if strings.Contains(message, "small") || strings.Contains(message, "compact") {
		requirements["size"] = "small"
	} else if strings.Contains(message, "large") || strings.Contains(message, "big") {
		requirements["size"] = "large"
	}

	// Color preferences
	colors := []string{"black", "white", "red", "blue", "green", "silver", "gold", "gray"}
	for _, color := range colors {
		if strings.Contains(message, color) {
			requirements["color"] = color
			break
		}
	}

	// Speed/performance requirements
	if strings.Contains(message, "fast") || strings.Contains(message, "quick") {
		requirements["performance"] = "high"
	} else if strings.Contains(message, "slow") || strings.Contains(message, "basic") {
		requirements["performance"] = "low"
	}

	// Quality requirements
	if strings.Contains(message, "premium") || strings.Contains(message, "best") || strings.Contains(message, "high quality") {
		requirements["quality"] = "premium"
	} else if strings.Contains(message, "budget") || strings.Contains(message, "cheap") {
		requirements["quality"] = "budget"
	}

	// Brand preference
	brands := s.ExtractBrands(message)
	if len(brands) > 0 {
		requirements["preferred_brands"] = brands
	}

	// Price range
	prices := s.ExtractPrices(message)
	if len(prices) > 0 {
		requirements["max_price"] = prices[0]
	}

	return requirements
}

// ShouldRefineSearch determines if search should be refined based on user response
func (s *MessageAnalysisService) ShouldRefineSearch(message string) bool {
	message = strings.ToLower(message)

	refinementIndicators := []string{
		"not what i want",
		"show me different",
		"something else",
		"other options",
		"not quite",
		"not right",
		"try again",
		"refine",
		"narrow down",
		"more specific",
	}

	for _, indicator := range refinementIndicators {
		if strings.Contains(message, indicator) {
			return true
		}
	}

	return false
}
