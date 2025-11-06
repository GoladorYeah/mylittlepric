package services

import (
	"fmt"
	"strings"

	"mylittleprice/internal/models"
)

// ContextDepth represents how much context should be sent to the AI
type ContextDepth int

const (
	// ContextDepthMinimal - Only last 1-2 messages + last product
	// Used for simple modifications like "cheaper", "different color"
	ContextDepthMinimal ContextDepth = 1

	// ContextDepthMedium - Last 3-4 messages + preferences + summary
	// Used for clarifications and follow-up questions
	ContextDepthMedium ContextDepth = 2

	// ContextDepthFull - Full cycle history (up to 6) + complete context
	// Used for complex queries and new categories
	ContextDepthFull ContextDepth = 3
)

// ContextOptimizerService determines optimal context depth for each request
type ContextOptimizerService struct {
	embedding *EmbeddingService
}

// NewContextOptimizerService creates a new context optimizer
func NewContextOptimizerService(embedding *EmbeddingService) *ContextOptimizerService {
	return &ContextOptimizerService{
		embedding: embedding,
	}
}

// DecideContextDepth analyzes the user message and determines optimal context depth
func (c *ContextOptimizerService) DecideContextDepth(
	userMessage string,
	session *models.ChatSession,
) ContextDepth {

	msgLower := strings.ToLower(userMessage)

	// 1. Simple price/feature modifiers - MINIMAL context
	if c.isSimpleModifier(msgLower) {
		fmt.Printf("🎯 Context depth: MINIMAL (simple modifier detected)\n")
		return ContextDepthMinimal
	}

	// 2. Short question or confirmation - MINIMAL context
	if c.isShortQuestion(userMessage) {
		fmt.Printf("🎯 Context depth: MINIMAL (short question)\n")
		return ContextDepthMinimal
	}

	// 3. New category or topic change - FULL context
	if c.isNewCategory(userMessage, session) {
		fmt.Printf("🎯 Context depth: FULL (new category detected)\n")
		return ContextDepthFull
	}

	// 4. Complex query with multiple requirements - FULL context
	if c.isComplexQuery(userMessage) {
		fmt.Printf("🎯 Context depth: FULL (complex query)\n")
		return ContextDepthFull
	}

	// 5. Clarification or follow-up - MEDIUM context
	if c.isClarification(msgLower) {
		fmt.Printf("🎯 Context depth: MEDIUM (clarification)\n")
		return ContextDepthMedium
	}

	// Default: MEDIUM for most queries
	fmt.Printf("🎯 Context depth: MEDIUM (default)\n")
	return ContextDepthMedium
}

// isSimpleModifier checks if message is a simple modification request
func (c *ContextOptimizerService) isSimpleModifier(msgLower string) bool {
	simpleModifiers := []string{
		// Price modifiers
		"подешевле", "подороже", "дешевле", "дороже",
		"cheaper", "expensive", "more expensive", "less expensive",
		"більш дешев", "більш дорог", "дешевш", "дорожч",
		"lower price", "higher price",

		// Size/storage modifiers
		"больше памяти", "меньше памяти", "larger", "smaller",
		"more storage", "less storage", "більше пам'яті", "менше пам'яті",

		// Color/variant modifiers
		"другой цвет", "другого цвета", "other color", "different color",
		"інший колір", "іншого кольору",
		"другая модель", "другую модель", "other model",

		// Quantity modifiers
		"больше вариантов", "другие варианты", "more options", "other options",
		"більше варіантів", "інші варіанти",

		// Simple affirmations with modifiers
		"да, но подешевле", "yes, but cheaper", "так, але дешевше",
		"да, другой", "yes, different", "так, інший",
	}

	for _, modifier := range simpleModifiers {
		if strings.Contains(msgLower, modifier) {
			return true
		}
	}

	return false
}

// isShortQuestion checks if message is a short question/confirmation
func (c *ContextOptimizerService) isShortQuestion(msg string) bool {
	// Messages under 30 characters that are questions or confirmations
	if len(msg) <= 30 {
		msgLower := strings.ToLower(msg)

		shortPatterns := []string{
			"да", "yes", "так", "ок", "ok", "okay",
			"нет", "no", "ні",
			"покажи", "show", "покажи",
			"это", "this", "це",
			"первый", "второй", "first", "second", "перший", "другий",
		}

		for _, pattern := range shortPatterns {
			if strings.Contains(msgLower, pattern) {
				return true
			}
		}
	}

	return false
}

// isNewCategory checks if user is asking about a different product category
func (c *ContextOptimizerService) isNewCategory(msg string, session *models.ChatSession) bool {
	// If no category set yet, it's new
	if session.SearchState.Category == "" {
		return true
	}

	// Detect category using embedding service
	detectedCategory := ""
	if c.embedding != nil {
		detectedCategory = c.embedding.DetectCategory(msg)
	}

	// If detected category differs significantly from current, it's new
	if detectedCategory != "" && detectedCategory != session.SearchState.Category {
		// Check if it's a related subcategory or completely different
		if !c.isRelatedCategory(detectedCategory, session.SearchState.Category) {
			return true
		}
	}

	// Keywords indicating category change
	categoryChangeKeywords := []string{
		"а теперь", "а ещё", "also", "also need", "і ще",
		"другое", "something else", "щось інше",
		"вместо этого", "instead", "замість",
	}

	msgLower := strings.ToLower(msg)
	for _, kw := range categoryChangeKeywords {
		if strings.Contains(msgLower, kw) {
			return true
		}
	}

	return false
}

// isComplexQuery checks if query has multiple requirements
func (c *ContextOptimizerService) isComplexQuery(msg string) bool {
	// Count requirement indicators
	requirementCount := 0

	// Price mentions
	priceKeywords := []string{"price", "цена", "ціна", "budget", "бюджет", "$", "€", "₴", "uah", "usd"}
	for _, kw := range priceKeywords {
		if strings.Contains(strings.ToLower(msg), kw) {
			requirementCount++
			break
		}
	}

	// Feature mentions
	featureKeywords := []string{
		"with", "со", "із", "з",
		"memory", "storage", "память", "пам'ять",
		"screen", "display", "экран", "дисплей",
		"camera", "камера",
		"battery", "батарея", "акумулятор",
	}
	for _, kw := range featureKeywords {
		if strings.Contains(strings.ToLower(msg), kw) {
			requirementCount++
			break
		}
	}

	// Brand mentions
	brandKeywords := []string{"apple", "samsung", "xiaomi", "google", "oneplus", "sony", "lg"}
	for _, kw := range brandKeywords {
		if strings.Contains(strings.ToLower(msg), kw) {
			requirementCount++
			break
		}
	}

	// Condition mentions
	conditionKeywords := []string{"new", "новый", "новий", "warranty", "гарантия", "гарантія"}
	for _, kw := range conditionKeywords {
		if strings.Contains(strings.ToLower(msg), kw) {
			requirementCount++
			break
		}
	}

	// If 3+ requirements, it's complex
	if requirementCount >= 3 {
		return true
	}

	// Long messages (>100 chars) with multiple clauses
	if len(msg) > 100 && strings.Count(msg, ",") >= 2 {
		return true
	}

	return false
}

// isClarification checks if message is asking for clarification or providing more details
func (c *ContextOptimizerService) isClarification(msgLower string) bool {
	clarificationKeywords := []string{
		// Questions
		"какой", "какая", "какие", "which", "what", "який", "яка", "які",
		"сколько", "how much", "how many", "скільки",
		"когда", "when", "коли",
		"где", "where", "де",

		// Answers to questions
		"например", "for example", "наприклад",
		"я ищу", "i'm looking", "i need", "мне нужен", "мені потрібен",
		"хочу", "want", "хочу",
		"предпочитаю", "prefer", "віддаю перевагу",
	}

	for _, kw := range clarificationKeywords {
		if strings.Contains(msgLower, kw) {
			return true
		}
	}

	return false
}

// isRelatedCategory checks if two categories are related (e.g., "smartphones" and "brand_specific:apple_iphone")
func (c *ContextOptimizerService) isRelatedCategory(cat1, cat2 string) bool {
	// If one contains the other, they're related
	if strings.Contains(cat1, cat2) || strings.Contains(cat2, cat1) {
		return true
	}

	// Both are electronics-related
	electronicsCategories := []string{"smartphones", "laptops", "tablets", "headphones", "smartwatches"}
	isElectronics1 := false
	isElectronics2 := false

	for _, eCat := range electronicsCategories {
		if strings.Contains(strings.ToLower(cat1), eCat) {
			isElectronics1 = true
		}
		if strings.Contains(strings.ToLower(cat2), eCat) {
			isElectronics2 = true
		}
	}

	if isElectronics1 && isElectronics2 {
		return true
	}

	return false
}

// ShouldUpdateContext determines if conversation context should be updated
func (c *ContextOptimizerService) ShouldUpdateContext(session *models.ChatSession) bool {
	// Update context every 3-4 messages
	if session.CycleState.Iteration%3 == 0 {
		return true
	}

	// Update at end of cycle
	if session.CycleState.Iteration >= MaxIterations {
		return true
	}

	// Update if context is stale (>5 minutes)
	if session.ConversationContext != nil {
		timeSinceUpdate := session.UpdatedAt.Sub(session.ConversationContext.UpdatedAt)
		if timeSinceUpdate.Minutes() > 5 {
			return true
		}
	}

	return false
}
