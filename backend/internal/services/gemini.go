package services

import (
	"context"
	"encoding/json"
	"fmt"
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
	promptManager     *PromptManager
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
		promptManager:     NewPromptManager(),
		ctx:               context.Background(),
	}

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

	fmt.Printf("🚀 Processing with category: %s\n", currentCategory)

	decision := g.groundingStrategy.Decide(userMessage, conversationHistory)
	g.groundingStats.RecordDecision(decision)

	useGrounding := decision.ShouldUseGrounding && g.config.GeminiUseGrounding

	if useGrounding {
		fmt.Printf("   🔍 GROUNDING (reason: %s)\n", decision.Reason)
		return g.processWithGrounding(client, userMessage, conversationHistory, country, language, currentCategory, keyIndex)
	} else {
		fmt.Printf("   💬 DIRECT (reason: %s)\n", decision.Reason)
		return g.processWithDirectJSON(client, userMessage, conversationHistory, country, language, currentCategory, keyIndex)
	}
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

	currency := getCurrencyForCountry(country)

	var promptKey string
	if currentCategory == "" {
		promptKey = "master"
	} else if currentCategory == "electronics" {
		promptKey = "specialized_electronics"
	} else if currentCategory == "generic_model" {
		promptKey = "specialized_generic"
	} else {
		promptKey = "specialized_parametric"
	}

	basePrompt := g.promptManager.GetPrompt(promptKey, country, language, currency, currentCategory)

	if basePrompt == "" {
		return nil, keyIndex, fmt.Errorf("prompt not found: %s", promptKey)
	}

	prompt := basePrompt

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

	prompt += fmt.Sprintf("\n# USER MESSAGE\nUser: %s\n\nRespond with JSON ONLY:\n", userMessage)

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

func (g *GeminiService) processWithGrounding(
	client *genai.Client,
	userMessage string,
	history []map[string]string,
	country string,
	language string,
	currentCategory string,
	keyIndex int,
) (*models.GeminiResponse, int, error) {

	currency := getCurrencyForCountry(country)
	languageName := getLanguageName(language)

	categoryInfo := ""
	if currentCategory != "" {
		categoryInfo = fmt.Sprintf("\nCATEGORY: %s (already determined)", currentCategory)
	}

	isBrandSelection := g.detectBrandSelection(userMessage, history)

	var prompt string

	if isBrandSelection {
		currentYear := time.Now().Year()
		prompt = fmt.Sprintf(`Shopping assistant for %s with Google Search.
Language: %s, Currency: %s%s

# TASK
User selected brand (%s). Use Google Search to find LATEST models.

Search for models from %d and %d!

Response JSON:
{
  "response_type": "dialogue",
  "output": "Which model? (max 150 chars)",
  "quick_replies": ["Latest Model 1 ($X)", "Latest Model 2 ($Y)", "Model 3 ($Z)", "Older Model"]
}

CRITICAL: quick_replies must have newest models first with prices!
`, country, languageName, currency, categoryInfo, userMessage, currentYear, currentYear+1)
	} else {
		prompt = fmt.Sprintf(`Shopping assistant for %s with Google Search.
Language: %s, Currency: %s%s

# TASK
User asked about specific product. Use Google Search to:
1. Verify product exists
2. Check availability
3. Confirm exact model name

Response JSON:

## If product EXISTS:
{"response_type":"search","search_phrase":"exact product","search_type":"exact"}

## If NOT available:
{"response_type":"dialogue","output":"Not available. Alternatives:","quick_replies":["Alt1","Alt2","Alt3"]}
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

	prompt += fmt.Sprintf("\n# USER MESSAGE\nUser: %s\n\nUse Google Search, respond JSON:\n", userMessage)

	tools := []*genai.Tool{
		{GoogleSearch: &genai.GoogleSearch{}},
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
		return nil, keyIndex, fmt.Errorf("empty response")
	}

	return g.parseGeminiResponse(responseText, currentCategory)
}

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

	fmt.Printf("   📊 Tokens: In=%d, Out=%d, Total=%d\n", inputTokens, outputTokens, totalTokens)

	g.recordTokenUsage(inputTokens, outputTokens, totalTokens, withGrounding)
}

func (g *GeminiService) logGroundingUsage(metadata *genai.GroundingMetadata) {
	fmt.Printf("   ✅ Grounding used!\n")

	if len(metadata.WebSearchQueries) > 0 {
		fmt.Printf("   🔎 Queries: %v\n", metadata.WebSearchQueries)
	}

	if len(metadata.GroundingChunks) > 0 {
		fmt.Printf("   📚 Chunks: %d\n", len(metadata.GroundingChunks))
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
