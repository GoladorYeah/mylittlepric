package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"google.golang.org/genai"

	"mylittleprice/internal/config"
	"mylittleprice/internal/models"
	"mylittleprice/internal/utils"
)

type GeminiService struct {
	keyRotator        *utils.KeyRotator
	config            *config.Config
	groundingStrategy *GroundingStrategy
	groundingStats    *GroundingStats
	tokenStats        *TokenStats
	categoryPrompts   map[string]string
	promptMutex       sync.RWMutex
	ctx               context.Context
}

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
		categoryPrompts:   make(map[string]string),
		ctx:               context.Background(),
	}

	service.loadCategoryPrompts()
	return service
}

func (g *GeminiService) ProcessMessageWithContext(
	userMessage string,
	conversationHistory []map[string]string,
	country string,
	language string,
	currentCategory string,
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

	if currentCategory != "" {
		fmt.Printf("   📂 Using existing category: %s (NO function call needed)\n", currentCategory)
	}

	decision := g.groundingStrategy.Decide(userMessage, conversationHistory)
	g.groundingStats.RecordDecision(decision)

	useGrounding := decision.ShouldUseGrounding && g.config.GeminiUseGrounding

	if useGrounding {
		fmt.Printf("   🔍 Strategy: GROUNDING (reason: %s, confidence: %.2f)\n", decision.Reason, decision.Confidence)
		return g.processWithGrounding(client, userMessage, conversationHistory, country, language, currentCategory, keyIndex)
	} else {
		fmt.Printf("   💬 Strategy: DIRECT JSON (reason: %s)\n", decision.Reason)
		return g.processWithDirectJSON(client, userMessage, conversationHistory, country, language, currentCategory, keyIndex)
	}
}

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

	if result.UsageMetadata != nil {
		g.logTokenUsage(result.UsageMetadata, true)
	}

	if len(result.Candidates) > 0 && result.Candidates[0].GroundingMetadata != nil {
		g.logGroundingUsage(result.Candidates[0].GroundingMetadata)
	}

	responseText := g.extractResponseText(result)

	if strings.TrimSpace(responseText) == "" {
		return nil, keyIndex, fmt.Errorf("empty response from Gemini")
	}

	return g.parseGeminiResponse(responseText, currentCategory)
}

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

	if result.UsageMetadata != nil {
		g.logTokenUsage(result.UsageMetadata, false)
	}

	responseText := g.extractResponseText(result)

	if strings.TrimSpace(responseText) == "" {
		return nil, keyIndex, fmt.Errorf("empty response from Gemini")
	}

	return g.parseGeminiResponse(responseText, currentCategory)
}

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

	if currentCategory != "" {
		g.promptMutex.RLock()
		categoryPrompt, hasCustomPrompt := g.categoryPrompts[currentCategory]
		g.promptMutex.RUnlock()

		if hasCustomPrompt {
			prompt = strings.ReplaceAll(categoryPrompt, "{country}", country)
			prompt = strings.ReplaceAll(prompt, "{language}", languageName)
			prompt = strings.ReplaceAll(prompt, "{currency}", currency)
		} else {
			prompt = g.buildDefaultCategoryPrompt(country, languageName, currency, currentCategory)
		}
	} else {
		prompt = g.buildFirstMessagePrompt(country, languageName, currency)
	}

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

func (g *GeminiService) buildFirstMessagePrompt(country, languageName, currency string) string {
	return fmt.Sprintf(`You are a shopping assistant for %s.
Language: %s, Currency: %s

# YOUR TASK
Analyze user's message and respond with JSON format.

# CATEGORY MAPPING
- electronics: phones, laptops, tablets, TVs, cameras, headphones, smartwatches
- clothing: jackets, shirts, pants, shoes, dresses, accessories (bags, belts)
- furniture: sofas, tables, chairs, beds, desks, shelves, wardrobes
- kitchen: pans, pots, knives, appliances (coffee makers, blenders, toasters)
- sports: gym equipment, yoga mats, dumbbells, fitness trackers, bicycles
- tools: drills, saws, screwdrivers, power tools, hand tools
- decor: lamps, mirrors, vases, wall art, candles, frames
- textiles: pillows, blankets, bedding, towels, carpets, curtains

# RESPONSE FORMAT

## If user mentions a product category, respond:
{
  "response_type": "dialogue",
  "output": "Ask clarifying question (max 200 chars)",
  "quick_replies": ["Option1", "Option2", "Option3", "Option4"],
  "category": "electronics|clothing|furniture|kitchen|sports|tools|decor|textiles"
}

## If user gives specific product (brand + model), respond:
{
  "response_type": "search",
  "search_phrase": "exact product name",
  "search_type": "exact",
  "category": "electronics|clothing|furniture|kitchen|sports|tools|decor|textiles"
}

# EXAMPLES

User: "need phone"
→ {"response_type":"dialogue","output":"Which brand?","quick_replies":["Apple","Samsung","Google","Xiaomi"],"category":"electronics"}

User: "iPhone 15 Pro"
→ {"response_type":"search","search_phrase":"iPhone 15 Pro","search_type":"exact","category":"electronics"}

User: "need pillow"
→ {"response_type":"dialogue","output":"What size?","quick_replies":["50x70cm","70x70cm","40x60cm","Other"],"category":"textiles"}

User: "looking for lamp"
→ {"response_type":"dialogue","output":"What type?","quick_replies":["Table lamp","Floor lamp","Desk lamp","Wall lamp"],"category":"decor"}

`, country, languageName, currency)
}

func (g *GeminiService) buildDefaultCategoryPrompt(country, languageName, currency, category string) string {
	return fmt.Sprintf(`You are a shopping assistant for %s.
Language: %s, Currency: %s
CATEGORY: %s (ALREADY DETERMINED)

# YOUR TASK
Help user find products. Ask MAXIMUM 1-2 questions, then SEARCH.

# CRITICAL RULES
1. After 1-2 questions → ALWAYS search
2. If user gave specifications → SEARCH immediately
3. Don't ask more than 2 questions total
4. Prefer searching over asking

# RESPONSE FORMAT

## If need ONE more detail (max 2 questions total):
{
  "response_type": "dialogue",
  "output": "One specific question (max 150 chars)",
  "quick_replies": ["Option1", "Option2", "Option3", "Option4"]
}

## If have enough info OR asked 1+ questions already:
{
  "response_type": "search",
  "search_phrase": "product with all known specifications",
  "search_type": "exact|parameters|category"
}

# SEARCH TYPE
- "exact": User gave specific brand+model (Brother CS10)
- "parameters": User gave category+specs (pillow 50x70 memory foam)
- "category": General category search (pillow)

# EXAMPLES FOR %s

User: "soft" (after asking about pillow type)
→ SEARCH: {"response_type":"search","search_phrase":"pillow soft","search_type":"parameters"}

User: "memory foam" (first specification given)
→ SEARCH: {"response_type":"search","search_phrase":"pillow memory foam","search_type":"parameters"}

User: "50x70" (size specification)
→ SEARCH: {"response_type":"search","search_phrase":"pillow 50x70","search_type":"parameters"}

REMEMBER: After 1-2 questions → ALWAYS SEARCH! Don't keep asking.

`, country, languageName, currency, category, category)
}

func (g *GeminiService) loadCategoryPrompts() {
	categories := []string{
		"electronics", "clothing", "furniture", "kitchen",
		"sports", "tools", "decor", "textiles",
	}

	g.promptMutex.Lock()
	defer g.promptMutex.Unlock()

	loadedCount := 0

	for _, category := range categories {
		promptPath := fmt.Sprintf("internal/services/prompts/%s_prompt.txt", category)
		content, err := os.ReadFile(promptPath)
		if err != nil {
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

	// Detect if this is a brand selection query
	isBrandSelection := g.detectBrandSelection(userMessage, history)

	var prompt string

	if isBrandSelection {
		// Special prompt for getting latest models after brand selection
		currentYear := time.Now().Year()
		prompt = fmt.Sprintf(`You are a shopping assistant for %s with Google Search access.
Language: %s, Currency: %s%s

# YOUR TASK
User selected a brand (%s). Use Google Search to find the LATEST available models for this brand.

IMPORTANT: Search for models available in %d and %d!

Then respond with JSON format:

{
  "response_type": "dialogue",
  "output": "Which model? (max 150 chars)",
  "quick_replies": ["Latest Model 1", "Latest Model 2", "Previous Model 1", "Previous Model 2"]
}

CRITICAL: 
- quick_replies MUST contain the newest/latest models first!
- Include both current year (%d) and previous year models
- Use REAL model names from Google Search results
- Order: Newest → Older

`, country, languageName, currency, categoryInfo, userMessage,
			currentYear, currentYear+1, currentYear)

	} else {
		// Regular grounding prompt for product verification
		prompt = fmt.Sprintf(`You are a shopping assistant for %s with Google Search access.
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
	}

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

	prompt += fmt.Sprintf("\n# USER MESSAGE\nUser: %s\n\nUse Google Search, then respond with JSON only:\n", userMessage)

	return prompt
}

// Helper function to detect brand selection
func (g *GeminiService) detectBrandSelection(userMessage string, history []map[string]string) bool {
	messageLower := strings.ToLower(userMessage)

	brands := []string{
		"apple", "iphone", "samsung", "google", "pixel",
		"xiaomi", "oneplus", "sony", "lg", "dell", "hp",
		"lenovo", "asus", "acer", "msi", "huawei",
	}

	isBrand := false
	for _, brand := range brands {
		if strings.Contains(messageLower, brand) {
			isBrand = true
			break
		}
	}

	if !isBrand {
		return false
	}

	// Check if previous question was about brand
	if len(history) > 0 {
		for i := len(history) - 1; i >= 0; i-- {
			if history[i]["role"] == "assistant" {
				lastQuestion := strings.ToLower(history[i]["content"])
				if strings.Contains(lastQuestion, "which brand") ||
					strings.Contains(lastQuestion, "what brand") {
					return true
				}
				break
			}
		}
	}

	return false
}

func (g *GeminiService) parseGeminiResponse(responseText string, currentCategory string) (*models.GeminiResponse, int, error) {
	cleaned := strings.TrimSpace(responseText)
	cleaned = strings.TrimPrefix(cleaned, "```json")
	cleaned = strings.TrimPrefix(cleaned, "```")
	cleaned = strings.TrimSuffix(cleaned, "```")
	cleaned = strings.TrimSpace(cleaned)

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

	if geminiResp.ResponseType != "dialogue" && geminiResp.ResponseType != "search" {
		return nil, -1, fmt.Errorf("invalid response_type: %s", geminiResp.ResponseType)
	}

	if geminiResp.ResponseType == "dialogue" && geminiResp.Output == "" {
		return nil, -1, fmt.Errorf("dialogue missing output")
	}

	if geminiResp.ResponseType == "search" {
		if geminiResp.SearchPhrase == "" {
			return nil, -1, fmt.Errorf("search missing search_phrase")
		}
		if geminiResp.SearchType == "" {
			geminiResp.SearchType = "parameters"
		}
	}

	if geminiResp.Category == "" && currentCategory != "" {
		geminiResp.Category = currentCategory
	}

	return &geminiResp, -1, nil
}

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

func (g *GeminiService) logTokenUsage(metadata *genai.GenerateContentResponseUsageMetadata, withGrounding bool) {
	inputTokens := int64(metadata.PromptTokenCount)
	outputTokens := int64(metadata.CandidatesTokenCount)
	totalTokens := int64(metadata.TotalTokenCount)

	fmt.Printf("   📊 Tokens: Input=%d, Output=%d, Total=%d\n",
		inputTokens, outputTokens, totalTokens,
	)

	g.recordTokenUsage(inputTokens, outputTokens, totalTokens, withGrounding)
}

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

func (g *GeminiService) GetGroundingStats() *GroundingStats {
	return g.groundingStats
}

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
