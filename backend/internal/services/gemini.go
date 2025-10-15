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

	promptKey := g.promptManager.GetPromptKey(currentCategory)
	systemPrompt := g.promptManager.GetPrompt(promptKey, country, language, currentCategory)

	lastProductStr := ""
	if lastProduct != nil {
		lastProductStr = fmt.Sprintf("%s (%.2f)", lastProduct.Name, lastProduct.Price)
	}

	systemPrompt = strings.ReplaceAll(systemPrompt, "{last_product}", lastProductStr)

	conversationContext := g.buildConversationContext(conversationHistory)

	prompt := systemPrompt + "\n\n# CONVERSATION HISTORY:\n" + conversationContext + "\n\nCurrent user message: " + userMessage + "\n\nAnalyze the conversation history above. If the last assistant question was similar to what the current situation requires, provide a DIFFERENT question to move the conversation forward."

	temp := g.config.GeminiTemperature
	generateConfig := &genai.GenerateContentConfig{
		Temperature:      &temp,
		MaxOutputTokens:  int32(g.config.GeminiMaxOutputTokens),
		ResponseMIMEType: "application/json",
	}

	useGrounding := g.shouldUseGrounding(userMessage, conversationHistory, currentCategory)
	if useGrounding {
		generateConfig.Tools = []*genai.Tool{
			{GoogleSearchRetrieval: &genai.GoogleSearchRetrieval{}},
		}
	}

	resp, err := client.Models.GenerateContent(
		g.ctx,
		g.config.GeminiModel,
		genai.Text(prompt),
		generateConfig,
	)
	if err != nil {
		return nil, keyIndex, fmt.Errorf("Gemini API error: %w", err)
	}

	if resp == nil {
		return nil, keyIndex, fmt.Errorf("Gemini returned nil response")
	}

	if resp.UsageMetadata != nil {
		g.updateTokenStats(resp.UsageMetadata, useGrounding)
	}

	if len(resp.Candidates) == 0 {
		return nil, keyIndex, fmt.Errorf("no candidates in Gemini response")
	}

	candidate := resp.Candidates[0]
	if candidate.Content == nil || len(candidate.Content.Parts) == 0 {
		return nil, keyIndex, fmt.Errorf("no content in Gemini response")
	}

	responseText := ""
	for _, part := range candidate.Content.Parts {
		if part.Text != "" {
			responseText += part.Text
		}
	}

	responseText = strings.TrimSpace(responseText)
	responseText = strings.Trim(responseText, "`")
	responseText = strings.TrimPrefix(responseText, "json")
	responseText = strings.TrimSpace(responseText)

	if responseText == "" {
		return nil, keyIndex, fmt.Errorf("empty response text from Gemini")
	}

	var geminiResp models.GeminiResponse
	if err := json.Unmarshal([]byte(responseText), &geminiResp); err != nil {
		return nil, keyIndex, fmt.Errorf("failed to parse Gemini JSON response: %w (response: %s)", err, responseText)
	}

	if geminiResp.ResponseType == "" {
		return nil, keyIndex, fmt.Errorf("missing response_type in Gemini response")
	}

	return &geminiResp, keyIndex, nil
}

func (g *GeminiService) buildConversationContext(history []map[string]string) string {
	if len(history) == 0 {
		return "No previous messages"
	}

	var context strings.Builder
	for i, msg := range history {
		context.WriteString(fmt.Sprintf("%d. %s: %s\n", i+1, msg["role"], msg["content"]))
	}
	return context.String()
}

func (g *GeminiService) shouldUseGrounding(userMessage string, history []map[string]string, category string) bool {
	return false
}

func (g *GeminiService) updateTokenStats(metadata *genai.GenerateContentResponseUsageMetadata, withGrounding bool) {
	g.tokenStats.mu.Lock()
	defer g.tokenStats.mu.Unlock()

	g.tokenStats.TotalRequests++
	g.tokenStats.TotalInputTokens += int64(metadata.PromptTokenCount)

	outputTokens := int64(0)
	if metadata.TotalTokenCount > 0 && metadata.PromptTokenCount > 0 {
		outputTokens = int64(metadata.TotalTokenCount - metadata.PromptTokenCount)
	}
	g.tokenStats.TotalOutputTokens += outputTokens
	g.tokenStats.TotalTokens += int64(metadata.TotalTokenCount)

	if withGrounding {
		g.tokenStats.RequestsWithGrounding++
	}

	g.tokenStats.AverageInputTokens = float64(g.tokenStats.TotalInputTokens) / float64(g.tokenStats.TotalRequests)
	g.tokenStats.AverageOutputTokens = float64(g.tokenStats.TotalOutputTokens) / float64(g.tokenStats.TotalRequests)
}

func (g *GeminiService) GetTokenStats() *TokenStats {
	g.tokenStats.mu.RLock()
	defer g.tokenStats.mu.RUnlock()

	stats := *g.tokenStats
	return &stats
}

func (g *GeminiService) GetGroundingStats() *GroundingStats {
	return g.groundingStats
}
