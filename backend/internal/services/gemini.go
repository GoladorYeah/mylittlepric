package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"google.golang.org/genai"

	"mylittleprice/internal/config"
	"mylittleprice/internal/models"
	"mylittleprice/internal/utils"
)

type GeminiService struct {
	keyRotator     *utils.KeyRotator
	config         *config.Config
	promptManager  *PromptManager
	groundingStats *GroundingStats
	tokenStats     *TokenStats
	ctx            context.Context
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

type GroundingStats struct {
	TotalDecisions    int
	GroundingEnabled  int
	GroundingDisabled int
	ReasonCounts      map[string]int
	AverageConfidence float32
}

func NewGeminiService(keyRotator *utils.KeyRotator, cfg *config.Config) *GeminiService {
	return &GeminiService{
		keyRotator:     keyRotator,
		config:         cfg,
		promptManager:  NewPromptManager(),
		groundingStats: &GroundingStats{ReasonCounts: make(map[string]int)},
		tokenStats:     &TokenStats{},
		ctx:            context.Background(),
	}
}

func (g *GeminiService) ProcessMessageWithContext(
	userMessage string,
	conversationHistory []map[string]string,
	country string,
	language string,
	currentCategory string,
	lastProduct *models.ProductInfo,
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

	currency := getCurrencyForCountry(country)

	useGrounding := g.shouldUseGrounding(userMessage, conversationHistory, currentCategory)

	promptKey := g.promptManager.GetPromptKey(currentCategory)
	systemPrompt := g.promptManager.GetPrompt(promptKey, country, language, currency, currentCategory)

	lastProductStr := ""
	if lastProduct != nil {
		lastProductStr = fmt.Sprintf("%s (%.2f %s)", lastProduct.Name, lastProduct.Price, currency)
	}

	systemPrompt = strings.ReplaceAll(systemPrompt, "{last_product}", lastProductStr)

	fullPrompt := g.buildFullPrompt(systemPrompt, conversationHistory, userMessage, useGrounding)

	if useGrounding {
		fmt.Printf("   🔍 GROUNDING\n")
		g.groundingStats.GroundingEnabled++
		return g.processWithGrounding(client, fullPrompt, keyIndex, currentCategory)
	}

	fmt.Printf("   💬 DIRECT\n")
	g.groundingStats.GroundingDisabled++
	return g.processDirect(client, fullPrompt, keyIndex, currentCategory)
}

func (g *GeminiService) buildFullPrompt(systemPrompt string, history []map[string]string, userMessage string, useGrounding bool) string {
	prompt := systemPrompt

	if len(history) > 0 {
		maxHistory := 4
		startIdx := 0
		if len(history) > maxHistory {
			startIdx = len(history) - maxHistory
		}

		prompt += "\n\n# CONVERSATION HISTORY\n"
		for i := startIdx; i < len(history); i++ {
			msg := history[i]
			if msg["role"] == "user" {
				prompt += fmt.Sprintf("User: %s\n", msg["content"])
			} else if msg["role"] == "assistant" {
				prompt += fmt.Sprintf("You: %s\n", msg["content"])
			}
		}
	}

	if useGrounding {
		prompt += "\n# GROUNDING INSTRUCTION\nUse Google Search ONCE with query: '[brand] [product] models 2024 2025'. Get current models and prices. Do NOT search multiple times.\n"
	}

	prompt += fmt.Sprintf("\n# USER MESSAGE\nUser: %s\n\nRespond with JSON only:\n", userMessage)
	return prompt
}

func (g *GeminiService) shouldUseGrounding(userMessage string, history []map[string]string, currentCategory string) bool {
	if !g.config.GeminiUseGrounding {
		return false
	}

	messageLower := strings.ToLower(userMessage)
	wordCount := len(strings.Fields(userMessage))

	greetings := []string{"hello", "hi", "hey", "thanks", "bye", "yes", "no", "ok", "cheaper", "expensive", "more expensive", "less expensive"}
	for _, greeting := range greetings {
		if messageLower == greeting || strings.Contains(messageLower, greeting) {
			return false
		}
	}

	if currentCategory == "electronics" && len(history) >= 1 {
		lastQuestion := g.getLastAssistantMessage(history)
		if lastQuestion != "" {
			lastQLower := strings.ToLower(lastQuestion)
			if strings.Contains(lastQLower, "brand") {
				return true
			}
		}
	}

	if currentCategory == "generic_model" && len(history) >= 1 {
		lastQuestion := g.getLastAssistantMessage(history)
		if lastQuestion != "" {
			lastQLower := strings.ToLower(lastQuestion)
			if strings.Contains(lastQLower, "what type") || strings.Contains(lastQLower, "which type") || strings.Contains(lastQLower, "what size") {
				return true
			}
		}
	}

	if wordCount < g.config.GeminiGroundingMinWords {
		return false
	}

	recencyKeywords := []string{"latest", "newest", "new", "2024", "2025", "2026", "current"}
	for _, kw := range recencyKeywords {
		if strings.Contains(messageLower, kw) {
			return true
		}
	}

	return false
}

func (g *GeminiService) getLastAssistantMessage(history []map[string]string) string {
	for i := len(history) - 1; i >= 0; i-- {
		if history[i]["role"] == "assistant" {
			return history[i]["content"]
		}
	}
	return ""
}

func (g *GeminiService) processDirect(client *genai.Client, prompt string, keyIndex int, currentCategory string) (*models.GeminiResponse, int, error) {
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
		return nil, keyIndex, fmt.Errorf("empty response")
	}

	return g.parseGeminiResponse(responseText, currentCategory)
}

func (g *GeminiService) processWithGrounding(client *genai.Client, prompt string, keyIndex int, currentCategory string) (*models.GeminiResponse, int, error) {
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

	var geminiResp models.GeminiResponse
	if err := json.Unmarshal([]byte(cleaned), &geminiResp); err != nil {
		return nil, -1, fmt.Errorf("JSON parse error: %w", err)
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
