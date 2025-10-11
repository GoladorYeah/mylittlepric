package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	"google.golang.org/genai"

	"mylittleprice/internal/config"
	"mylittleprice/internal/models"
	"mylittleprice/internal/utils"
)

// GeminiService handles AI interactions with smart tool selection
type GeminiService struct {
	keyRotator        *utils.KeyRotator
	config            *config.Config
	groundingStrategy *GroundingStrategy
	groundingStats    *GroundingStats
	tokenStats        *TokenStats
	categoryPrompts   map[string]string // For custom category prompts
	promptMutex       sync.RWMutex      // Thread-safe access to prompts
	ctx               context.Context
}

// TokenStats tracks token usage for monitoring
type TokenStats struct {
	mu                    sync.RWMutex
	TotalRequests         int
	TotalInputTokens      int64
	TotalOutputTokens     int64
	TotalTokens           int64
	RequestsWithGrounding int
	AverageInputTokens    float64
	AverageOutputTokens   float64
}

// NewGeminiService creates a new Gemini service
func NewGeminiService(keyRotator *utils.KeyRotator, cfg *config.Config) *GeminiService {
	groundingMode := cfg.GeminiGroundingMode
	if groundingMode == "" {
		groundingMode = "balanced"
	}

	service := &GeminiService{
		keyRotator:        keyRotator,
		config:            cfg,
		groundingStrategy: NewGroundingStrategy(groundingMode),
		groundingStats:    &GroundingStats{ReasonCounts: make(map[string]int)},
		tokenStats:        &TokenStats{},
		categoryPrompts:   make(map[string]string), // Initialize prompts map
		ctx:               context.Background(),
	}

	// Load category-specific prompts from files (optional)
	service.loadCategoryPrompts()

	return service
}

// ═══════════════════════════════════════════════════════════════
// 🎯 MAIN PROCESSING METHOD (OPTIMIZED WITH CATEGORY CONTEXT)
// ═══════════════════════════════════════════════════════════════

// ProcessMessageWithContext processes user message with category context
// This is the OPTIMIZED method that avoids unnecessary function calls
func (g *GeminiService) ProcessMessageWithContext(
	userMessage string,
	conversationHistory []map[string]string,
	country string,
	language string,
	currentCategory string, // 🔥 KEY OPTIMIZATION: Pass existing category
) (*models.GeminiResponse, int, error) {

	apiKey, keyIndex, err := g.keyRotator.GetNextKey()
	if err != nil {
		return nil, -1, fmt.Errorf("failed to get API key: %w", err)
	}

	client, err := genai.NewClient(g.ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, keyIndex, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	fmt.Printf("🚀 Processing message with smart tool selection...\n")

	// 🎯 DECISION: Use category context to determine strategy
	if currentCategory != "" {
		fmt.Printf("   📂 Using existing category: %s (NO function call needed)\n", currentCategory)
	}

	// Smart grounding decision
	decision := g.groundingStrategy.Decide(userMessage, conversationHistory)
	g.groundingStats.RecordDecision(decision)

	useGrounding := decision.ShouldUseGrounding && g.config.GeminiUseGrounding

	if useGrounding {
		fmt.Printf("   🔍 Strategy: GROUNDING (reason: %s, confidence: %.2f)\n", decision.Reason, decision.Confidence)
		return g.processWithGrounding(client, userMessage, conversationHistory, country, language, currentCategory, keyIndex)
	} else {
		fmt.Printf("   💬 Strategy: DIRECT JSON (reason: %s)\n", decision.Reason)
		// 🔥 OPTIMIZATION: Use direct JSON prompt instead of function calling
		return g.processWithDirectJSON(client, userMessage, conversationHistory, country, language, currentCategory, keyIndex)
	}
}

// ═══════════════════════════════════════════════════════════════
// PROCESSING STRATEGIES
// ═══════════════════════════════════════════════════════════════

// processWithGrounding handles requests with Google Search (for specific products)
func (g *GeminiService) processWithGrounding(
	client *genai.Client,
	userMessage string,
	history []map[string]string,
	country string,
	language string,
	currentCategory string,
	keyIndex int,
) (*models.GeminiResponse, int, error) {

	prompt := g.buildGroundingPrompt(history, userMessage, country, language, currentCategory)

	tools := []*genai.Tool{
		{
			GoogleSearch: &genai.GoogleSearch{},
		},
	}

	generateConfig := &genai.GenerateContentConfig{
		Temperature:     &g.config.GeminiTemperature,
		MaxOutputTokens: int32(g.config.GeminiMaxOutputTokens),
		Tools:           tools,
	}

	result, err := client.Models.GenerateContent(
		g.ctx,
		g.config.GeminiModel,
		genai.Text(prompt),
		generateConfig,
	)
	if err != nil {
		return nil, keyIndex, fmt.Errorf("Gemini API error: %w", err)
	}

	// Log token usage
	if result.UsageMetadata != nil {
		g.logTokenUsage(result.UsageMetadata, true)
	}

	// Log grounding usage
	if len(result.Candidates) > 0 && result.Candidates[0].GroundingMetadata != nil {
		g.logGroundingUsage(result.Candidates[0].GroundingMetadata)
	}

	responseText := g.extractResponseText(result)

	if strings.TrimSpace(responseText) == "" {
		return nil, keyIndex, fmt.Errorf("empty response from Gemini")
	}

	return g.parseGeminiResponse(responseText, currentCategory)
}

// processWithDirectJSON handles requests with direct JSON (optimized, no tools)
func (g *GeminiService) processWithDirectJSON(
	client *genai.Client,
	userMessage string,
	history []map[string]string,
	country string,
	language string,
	currentCategory string,
	keyIndex int,
) (*models.GeminiResponse, int, error) {

	prompt := g.buildDirectJSONPrompt(history, userMessage, country, language, currentCategory)

	// NO TOOLS - just direct text generation (faster & cheaper!)
	generateConfig := &genai.GenerateContentConfig{
		Temperature:     &g.config.GeminiTemperature,
		MaxOutputTokens: int32(g.config.GeminiMaxOutputTokens),
	}

	result, err := client.Models.GenerateContent(
		g.ctx,
		g.config.GeminiModel,
		genai.Text(prompt),
		generateConfig,
	)
	if err != nil {
		return nil, keyIndex, fmt.Errorf("Gemini API error: %w", err)
	}

	// Log token usage
	if result.UsageMetadata != nil {
		g.logTokenUsage(result.UsageMetadata, false)
	}

	responseText := g.extractResponseText(result)

	if strings.TrimSpace(responseText) == "" {
		return nil, keyIndex, fmt.Errorf("empty response from Gemini")
	}

	return g.parseGeminiResponse(responseText, currentCategory)
}

// ═══════════════════════════════════════════════════════════════
// PROMPT BUILDERS
// ═══════════════════════════════════════════════════════════════

// buildDirectJSONPrompt creates optimized prompt without function calling
func (g *GeminiService) buildDirectJSONPrompt(
	history []map[string]string,
	userMessage string,
	country string,
	language string,
	currentCategory string,
) string {

	languageName := getLanguageName(language)
	currency := getCurrencyForCountry(country)

	var prompt string

	// 🆕 Check if we have a custom prompt for this category
	if currentCategory != "" {
		g.promptMutex.RLock()
		categoryPrompt, hasCustomPrompt := g.categoryPrompts[currentCategory]
		g.promptMutex.RUnlock()

		if hasCustomPrompt {
			// Use custom category-specific prompt
			prompt = strings.ReplaceAll(categoryPrompt, "{country}", country)
			prompt = strings.ReplaceAll(prompt, "{language}", languageName)
			prompt = strings.ReplaceAll(prompt, "{currency}", currency)
		} else {
			// Use default prompt for this category
			prompt = g.buildDefaultCategoryPrompt(country, languageName, currency, currentCategory)
		}
	} else {
		// First message - detect category
		prompt = g.buildFirstMessagePrompt(country, languageName, currency)
	}

	// Add conversation history
	if len(history) > 0 {
		maxHistory := 4
		startIdx := 0
		if len(history) > maxHistory {
			startIdx = len(history) - maxHistory
		}

		prompt += "\n# CONVERSATION HISTORY\n"
		for i := startIdx; i < len(history); i++ {
			msg := history[i]
			if msg["role"] == "user" {
				prompt += fmt.Sprintf("User: %s\n", msg["content"])
			} else if msg["role"] == "assistant" {
				prompt += fmt.Sprintf("You: %s\n", msg["content"])
			}
		}
	}

	prompt += fmt.Sprintf("\n# USER MESSAGE\nUser: %s\n\nRespond with JSON ONLY (no explanations):\n", userMessage)

	return prompt
}

// buildFirstMessagePrompt creates prompt for first message (category detection)
func (g *GeminiService) buildFirstMessagePrompt(country, languageName, currency string) string {
	return fmt.Sprintf(`You are a shopping assistant for %s.
Language: %s, Currency: %s

# YOUR TASK
Analyze user's message and respond with JSON format.

# RESPONSE FORMAT

## If user mentions a product category, respond:
{
  "response_type": "dialogue",
  "output": "Ask clarifying question (max 200 chars)",
  "quick_replies": ["Option1", "Option2", "Option3", "Option4"],
  "category": "electronics|clothing|furniture|kitchen|sports|tools|decor"
}

## If user gives specific product (brand + model), respond:
{
  "response_type": "search",
  "search_phrase": "exact product name",
  "search_type": "exact",
  "category": "electronics|clothing|furniture|kitchen|sports|tools|decor"
}

# EXAMPLES

User: "need phone"
→ {"response_type":"dialogue","output":"Which brand?","quick_replies":["Apple","Samsung","Google","Xiaomi"],"category":"electronics"}

User: "iPhone 15 Pro"
→ {"response_type":"search","search_phrase":"iPhone 15 Pro","search_type":"exact","category":"electronics"}

`, country, languageName, currency)
}

// buildDefaultCategoryPrompt creates default prompt when no custom prompt exists
func (g *GeminiService) buildDefaultCategoryPrompt(country, languageName, currency, category string) string {
	return fmt.Sprintf(`You are a shopping assistant for %s.
Language: %s, Currency: %s
CATEGORY: %s (ALREADY DETERMINED)

# YOUR TASK
Help user find products. Ask MAXIMUM 2-3 questions, then SEARCH.

# CRITICAL RULES
1. After 2-3 questions → ALWAYS search
2. If user gave brand → ask for model, then SEARCH
3. If user gave brand + specifications → SEARCH immediately
4. Don't ask more than 3 questions total
5. Prefer searching over asking

# RESPONSE FORMAT

## If need ONE more detail (max 3 questions total):
{
  "response_type": "dialogue",
  "output": "One specific question (max 150 chars)",
  "quick_replies": ["Option1", "Option2", "Option3", "Option4"]
}

## If have enough info OR asked 2+ questions already:
{
  "response_type": "search",
  "search_phrase": "brand product specifications",
  "search_type": "exact|parameters|category"
}

# SEARCH TYPE
- "exact": User gave specific model (Brother CS10)
- "parameters": User gave brand + specs (Brother mechanical 10-15 stitches)
- "category": General category search (sewing machine)

# EXAMPLES FOR %s

User: "Brother"
→ Ask model: {"response_type":"dialogue","output":"Which Brother model?","quick_replies":["CS10","XR3774","GX37"]}

User: "CS10" (after asking brand)
→ SEARCH: {"response_type":"search","search_phrase":"Brother CS10","search_type":"exact"}

User: "mechanical with 15 stitches" (after asking brand)
→ SEARCH: {"response_type":"search","search_phrase":"Brother mechanical sewing machine 15 stitches","search_type":"parameters"}

User: "basic model"
→ SEARCH: {"response_type":"search","search_phrase":"Brother basic sewing machine","search_type":"category"}

REMEMBER: After 2-3 questions → ALWAYS SEARCH! Don't keep asking.

`, country, languageName, currency, category, category)
}

// 🆕 loadCategoryPrompts loads category-specific prompts from files
func (g *GeminiService) loadCategoryPrompts() {
	categories := []string{
		"electronics", "clothing", "furniture", "kitchen",
		"sports", "tools", "decor",
	}

	g.promptMutex.Lock()
	defer g.promptMutex.Unlock()

	loadedCount := 0

	for _, category := range categories {
		promptPath := fmt.Sprintf("internal/services/prompts/%s_prompt.txt", category)
		content, err := os.ReadFile(promptPath)
		if err != nil {
			// File doesn't exist - will use default prompt
			continue
		}

		g.categoryPrompts[category] = string(content)
		loadedCount++
	}

	if loadedCount > 0 {
		fmt.Printf("✅ Loaded %d custom category prompts\n", loadedCount)
	} else {
		fmt.Printf("ℹ️  No custom prompts found - using default prompts\n")
	}
}

// buildGroundingPrompt creates prompt for grounding (specific products)
func (g *GeminiService) buildGroundingPrompt(
	history []map[string]string,
	userMessage string,
	country string,
	language string,
	currentCategory string,
) string {

	languageName := getLanguageName(language)
	currency := getCurrencyForCountry(country)

	categoryInfo := ""
	if currentCategory != "" {
		categoryInfo = fmt.Sprintf("\nCATEGORY: %s (already determined)", currentCategory)
	}

	prompt := fmt.Sprintf(`You are a shopping assistant for %s with Google Search access.
Language: %s, Currency: %s%s

# YOUR TASK
User asked about a SPECIFIC product. Use Google Search to:
1. Verify the product exists
2. Check if it's currently available
3. Confirm the exact model name

Then respond with JSON format:

## If product EXISTS and is available:
{
  "response_type": "search",
  "search_phrase": "exact product name",
  "search_type": "exact"
}

## If product DOESN'T EXIST or is unavailable:
{
  "response_type": "dialogue",
  "output": "That product isn't available. Here are alternatives:",
  "quick_replies": ["Alternative 1", "Alternative 2", "Alternative 3"]
}

`, country, languageName, currency, categoryInfo)

	// Add history
	if len(history) > 0 {
		maxHistory := 3
		startIdx := 0
		if len(history) > maxHistory {
			startIdx = len(history) - maxHistory
		}

		prompt += "\n# RECENT CONVERSATION\n"
		for i := startIdx; i < len(history); i++ {
			msg := history[i]
			if msg["role"] == "user" {
				prompt += fmt.Sprintf("User: %s\n", msg["content"])
			} else if msg["role"] == "assistant" {
				prompt += fmt.Sprintf("You: %s\n", msg["content"])
			}
		}
	}

	prompt += fmt.Sprintf("\n# USER MESSAGE\nUser: %s\n\nUse Google Search to verify, then respond with JSON only:\n", userMessage)

	return prompt
}

// ═══════════════════════════════════════════════════════════════
// RESPONSE PARSING
// ═══════════════════════════════════════════════════════════════

// parseGeminiResponse parses JSON response from Gemini
func (g *GeminiService) parseGeminiResponse(responseText string, currentCategory string) (*models.GeminiResponse, int, error) {
	cleaned := strings.TrimSpace(responseText)
	cleaned = strings.TrimPrefix(cleaned, "```json")
	cleaned = strings.TrimPrefix(cleaned, "```")
	cleaned = strings.TrimSuffix(cleaned, "```")
	cleaned = strings.TrimSpace(cleaned)

	// Find JSON object
	if !strings.HasPrefix(cleaned, "{") {
		startIdx := strings.Index(cleaned, "{")
		if startIdx != -1 {
			cleaned = cleaned[startIdx:]
		}
	}

	if !strings.HasSuffix(cleaned, "}") {
		endIdx := strings.LastIndex(cleaned, "}")
		if endIdx != -1 {
			cleaned = cleaned[:endIdx+1]
		}
	}

	if cleaned == "" {
		return nil, -1, fmt.Errorf("empty response")
	}

	var geminiResp models.GeminiResponse
	if err := json.Unmarshal([]byte(cleaned), &geminiResp); err != nil {
		return nil, -1, fmt.Errorf("JSON parse error: %w\nResponse: %s", err, cleaned)
	}

	// Validate response type
	if geminiResp.ResponseType != "dialogue" && geminiResp.ResponseType != "search" {
		return nil, -1, fmt.Errorf("invalid response_type: %s", geminiResp.ResponseType)
	}

	// Validate dialogue response
	if geminiResp.ResponseType == "dialogue" && geminiResp.Output == "" {
		return nil, -1, fmt.Errorf("dialogue missing output")
	}

	// Validate search response
	if geminiResp.ResponseType == "search" {
		if geminiResp.SearchPhrase == "" {
			return nil, -1, fmt.Errorf("search missing search_phrase")
		}
		if geminiResp.SearchType == "" {
			geminiResp.SearchType = "parameters"
		}
	}

	// 🎯 PRESERVE CATEGORY: If category not in response but we have it in context
	if geminiResp.Category == "" && currentCategory != "" {
		geminiResp.Category = currentCategory
	}

	return &geminiResp, -1, nil
}

// ═══════════════════════════════════════════════════════════════
// HELPER METHODS
// ═══════════════════════════════════════════════════════════════

// extractResponseText extracts text from Gemini result
func (g *GeminiService) extractResponseText(result *genai.GenerateContentResponse) string {
	if result == nil || len(result.Candidates) == 0 {
		return ""
	}

	candidate := result.Candidates[0]
	if candidate.Content == nil || len(candidate.Content.Parts) == 0 {
		return ""
	}

	var responseText string
	for _, part := range candidate.Content.Parts {
		if part.Text != "" {
			responseText += part.Text
		}
	}

	return responseText
}

// logTokenUsage logs token usage statistics
func (g *GeminiService) logTokenUsage(metadata *genai.GenerateContentResponseUsageMetadata, withGrounding bool) {
	inputTokens := int64(metadata.PromptTokenCount)
	outputTokens := int64(metadata.CandidatesTokenCount)
	totalTokens := int64(metadata.TotalTokenCount)

	fmt.Printf("   📊 Tokens: Input=%d, Output=%d, Total=%d\n",
		inputTokens, outputTokens, totalTokens,
	)

	g.recordTokenUsage(inputTokens, outputTokens, totalTokens, withGrounding)
}

// logGroundingUsage logs grounding usage information
func (g *GeminiService) logGroundingUsage(metadata *genai.GroundingMetadata) {
	fmt.Printf("   ✅ Grounding was used!\n")

	if len(metadata.WebSearchQueries) > 0 {
		fmt.Printf("   🔎 Search queries: %v\n", metadata.WebSearchQueries)
	}

	if len(metadata.GroundingChunks) > 0 {
		fmt.Printf("   📚 Retrieved %d grounding chunks\n", len(metadata.GroundingChunks))
		if metadata.GroundingChunks[0].Web != nil {
			fmt.Printf("      📄 %s\n", metadata.GroundingChunks[0].Web.Title)
		}
	}
}

// recordTokenUsage records token statistics
func (g *GeminiService) recordTokenUsage(inputTokens, outputTokens, totalTokens int64, withGrounding bool) {
	g.tokenStats.mu.Lock()
	defer g.tokenStats.mu.Unlock()

	g.tokenStats.TotalRequests++
	g.tokenStats.TotalInputTokens += inputTokens
	g.tokenStats.TotalOutputTokens += outputTokens
	g.tokenStats.TotalTokens += totalTokens

	if withGrounding {
		g.tokenStats.RequestsWithGrounding++
	}

	if g.tokenStats.TotalRequests > 0 {
		g.tokenStats.AverageInputTokens = float64(g.tokenStats.TotalInputTokens) / float64(g.tokenStats.TotalRequests)
		g.tokenStats.AverageOutputTokens = float64(g.tokenStats.TotalOutputTokens) / float64(g.tokenStats.TotalRequests)
	}
}

// ═══════════════════════════════════════════════════════════════
// STATISTICS METHODS
// ═══════════════════════════════════════════════════════════════

// GetGroundingStats returns grounding statistics
func (g *GeminiService) GetGroundingStats() *GroundingStats {
	return g.groundingStats
}

// GetTokenStats returns token usage statistics
func (g *GeminiService) GetTokenStats() map[string]interface{} {
	g.tokenStats.mu.RLock()
	defer g.tokenStats.mu.RUnlock()

	groundingPercentage := float64(0)
	if g.tokenStats.TotalRequests > 0 {
		groundingPercentage = float64(g.tokenStats.RequestsWithGrounding) / float64(g.tokenStats.TotalRequests) * 100
	}

	return map[string]interface{}{
		"total_requests":          g.tokenStats.TotalRequests,
		"total_input_tokens":      g.tokenStats.TotalInputTokens,
		"total_output_tokens":     g.tokenStats.TotalOutputTokens,
		"total_tokens":            g.tokenStats.TotalTokens,
		"requests_with_grounding": g.tokenStats.RequestsWithGrounding,
		"average_input_tokens":    fmt.Sprintf("%.1f", g.tokenStats.AverageInputTokens),
		"average_output_tokens":   fmt.Sprintf("%.1f", g.tokenStats.AverageOutputTokens),
		"grounding_percentage":    fmt.Sprintf("%.1f%%", groundingPercentage),
	}
}

// ═══════════════════════════════════════════════════════════════
// DEPRECATED (for backwards compatibility)
// ═══════════════════════════════════════════════════════════════

// ProcessMessage is deprecated - use ProcessMessageWithContext instead
// Kept for backwards compatibility
func (g *GeminiService) ProcessMessage(
	userMessage string,
	conversationHistory []map[string]string,
	country string,
	language string,
) (*models.GeminiResponse, int, error) {
	return g.ProcessMessageWithContext(userMessage, conversationHistory, country, language, "")
}

// getLanguageName returns full language name from code
func getLanguageName(langCode string) string {
	languageNames := map[string]string{
		"de": "German", "fr": "French", "it": "Italian",
		"es": "Spanish", "pt": "Portuguese", "nl": "Dutch",
		"pl": "Polish", "cs": "Czech", "sv": "Swedish",
		"no": "Norwegian", "da": "Danish", "fi": "Finnish",
		"en": "English",
	}
	if name, ok := languageNames[langCode]; ok {
		return name
	}
	return "English"
}
