// backend/internal/services/gemini.go
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
	client            *genai.Client
	keyRotator        *utils.KeyRotator
	config            *config.Config
	promptManager     *PromptManager
	groundingStats    *GroundingStats
	groundingStrategy *GroundingStrategy // ← ДОБАВИТЬ ЭТО
	tokenStats        *TokenStats
	embedding         *EmbeddingService
	ctx               context.Context
	mu                sync.RWMutex
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

// В конструкторе NewGeminiService добавьте инициализацию:
func NewGeminiService(keyRotator *utils.KeyRotator, cfg *config.Config, embedding *EmbeddingService) *GeminiService {
	ctx := context.Background()

	apiKey, _, err := keyRotator.GetNextKey()
	if err != nil {
		panic(fmt.Errorf("failed to get initial API key: %w", err))
	}

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		panic(fmt.Errorf("failed to create Gemini client: %w", err))
	}

	return &GeminiService{
		client:            client,
		keyRotator:        keyRotator,
		config:            cfg,
		promptManager:     NewPromptManager(),
		groundingStats:    &GroundingStats{ReasonCounts: make(map[string]int)},
		groundingStrategy: NewGroundingStrategy(embedding), // ← ДОБАВИТЬ ЭТО
		tokenStats:        &TokenStats{},
		embedding:         embedding,
		ctx:               ctx,
	}
}

func (g *GeminiService) rotateClient() error {
	g.mu.Lock()
	defer g.mu.Unlock()

	apiKey, _, err := g.keyRotator.GetNextKey()
	if err != nil {
		return fmt.Errorf("failed to get API key: %w", err)
	}

	client, err := genai.NewClient(g.ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return fmt.Errorf("failed to create Gemini client: %w", err)
	}

	g.client = client
	return nil
}

func (g *GeminiService) ProcessMessageWithContext(
	userMessage string,
	conversationHistory []map[string]string,
	country string,
	language string,
	currentCategory string,
	lastProduct *models.ProductInfo,
) (*models.GeminiResponse, int, error) {

	// УБРАЛИ rotateClient() отсюда!
	// if err := g.rotateClient(); err != nil {
	// 	return nil, -1, err
	// }

	if currentCategory == "" {
		detectedCategory := g.embedding.DetectCategory(userMessage)
		if detectedCategory != "" {
			currentCategory = detectedCategory
		}
	}

	promptKey := g.promptManager.GetPromptKey(currentCategory)
	systemPrompt := g.promptManager.GetPrompt(promptKey, country, language, currentCategory)

	lastProductStr := ""
	if lastProduct != nil {
		lastProductStr = fmt.Sprintf("%s (%.2f)", lastProduct.Name, lastProduct.Price)
	}

	systemPrompt = strings.ReplaceAll(systemPrompt, "{last_product}", lastProductStr)
	conversationContext := g.buildConversationContext(conversationHistory)

	prompt := systemPrompt + "\n\n# CONVERSATION HISTORY:\n" + conversationContext +
		"\n\nCurrent user message: " + userMessage +
		"\n\nCRITICAL INSTRUCTIONS:\n- You MUST respond with valid JSON only\n- If using grounding/search results, incorporate the information naturally\n- ALWAYS end your response with valid JSON in this exact format:\n{\"response_type\":\"dialogue\",\"output\":\"...\",\"quick_replies\":[...],\"category\":\"...\"}\nOR\n{\"response_type\":\"search\",\"search_phrase\":\"...\",\"search_type\":\"...\",\"category\":\"...\"}\n\nAnalyze the conversation history above. If the last assistant question was similar to what the current situation requires, provide a DIFFERENT question to move the conversation forward."

	temp := g.config.GeminiTemperature
	generateConfig := &genai.GenerateContentConfig{
		Temperature:     &temp,
		MaxOutputTokens: int32(g.config.GeminiMaxOutputTokens),
	}

	useGrounding := g.shouldUseGrounding(userMessage, conversationHistory, currentCategory)
	if useGrounding {
		generateConfig.Tools = []*genai.Tool{
			{GoogleSearch: &genai.GoogleSearch{}},
		}
	} else {
		generateConfig.ResponseMIMEType = "application/json"
	}

	g.mu.RLock()
	client := g.client
	g.mu.RUnlock()

	resp, err := client.Models.GenerateContent(
		g.ctx,
		g.config.GeminiModel,
		genai.Text(prompt),
		generateConfig,
	)

	// Если ошибка - пробуем ротировать ключ и повторить
	if err != nil {
		// Проверяем если это quota/rate limit ошибка
		if strings.Contains(err.Error(), "quota") ||
			strings.Contains(err.Error(), "429") ||
			strings.Contains(err.Error(), "RESOURCE_EXHAUSTED") {

			// Ротируем клиент
			if rotateErr := g.rotateClient(); rotateErr != nil {
				return nil, 0, fmt.Errorf("Gemini API error: %w, rotation failed: %v", err, rotateErr)
			}

			// Повторяем запрос с новым клиентом
			g.mu.RLock()
			client = g.client
			g.mu.RUnlock()

			resp, err = client.Models.GenerateContent(
				g.ctx,
				g.config.GeminiModel,
				genai.Text(prompt),
				generateConfig,
			)

			if err != nil {
				return nil, 0, fmt.Errorf("Gemini API error after rotation: %w", err)
			}
		} else {
			return nil, 0, fmt.Errorf("Gemini API error: %w", err)
		}
	}

	if resp == nil {
		return nil, 0, fmt.Errorf("Gemini returned nil response")
	}

	if resp.UsageMetadata != nil {
		g.updateTokenStats(resp.UsageMetadata, useGrounding)
	}

	if len(resp.Candidates) == 0 {
		return nil, 0, fmt.Errorf("no candidates in Gemini response")
	}

	candidate := resp.Candidates[0]
	if candidate.Content == nil || len(candidate.Content.Parts) == 0 {
		return nil, 0, fmt.Errorf("no content in Gemini response")
	}

	responseText := ""
	for _, part := range candidate.Content.Parts {
		if part.Text != "" {
			responseText += part.Text
		}
	}

	responseText = strings.TrimSpace(responseText)

	// Если использовали grounding, извлекаем JSON из текста
	if useGrounding {
		responseText = g.extractJSONFromText(responseText)
	}

	responseText = strings.Trim(responseText, "`")
	responseText = strings.TrimPrefix(responseText, "json")
	responseText = strings.TrimSpace(responseText)

	if responseText == "" {
		return nil, 0, fmt.Errorf("empty response text from Gemini")
	}

	var geminiResp models.GeminiResponse
	if err := json.Unmarshal([]byte(responseText), &geminiResp); err != nil {
		return nil, 0, fmt.Errorf("failed to parse Gemini JSON response: %w (response: %s)", err, responseText)
	}

	if geminiResp.ResponseType == "" {
		return nil, 0, fmt.Errorf("missing response_type in Gemini response")
	}

	return &geminiResp, 0, nil
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

// Замените функцию shouldUseGrounding:
func (g *GeminiService) shouldUseGrounding(userMessage string, history []map[string]string, category string) bool {
	if !g.config.GeminiUseGrounding {
		return false
	}

	// Используем умную стратегию
	decision := g.groundingStrategy.ShouldUseGrounding(userMessage, history, category)

	// Обновляем статистику
	g.groundingStats.TotalDecisions++
	if decision.UseGrounding {
		g.groundingStats.GroundingEnabled++
	} else {
		g.groundingStats.GroundingDisabled++
	}
	g.groundingStats.ReasonCounts[decision.Reason]++

	// Обновляем среднюю уверенность
	if g.groundingStats.TotalDecisions > 0 {
		currentAvg := g.groundingStats.AverageConfidence
		n := float32(g.groundingStats.TotalDecisions)
		g.groundingStats.AverageConfidence = (currentAvg*(n-1) + decision.Confidence) / n
	}

	return decision.UseGrounding
}

func (g *GeminiService) extractJSONFromText(text string) string {
	// Ищем JSON в тексте
	// Вариант 1: JSON в конце
	startIdx := strings.LastIndex(text, "{")
	endIdx := strings.LastIndex(text, "}")

	if startIdx != -1 && endIdx != -1 && endIdx > startIdx {
		return text[startIdx : endIdx+1]
	}

	// Вариант 2: JSON в code block
	if strings.Contains(text, "```json") {
		start := strings.Index(text, "```json")
		end := strings.Index(text[start+7:], "```")
		if end != -1 {
			return strings.TrimSpace(text[start+7 : start+7+end])
		}
	}

	// Вариант 3: JSON в code block без языка
	if strings.Contains(text, "```") {
		start := strings.Index(text, "```")
		end := strings.Index(text[start+3:], "```")
		if end != -1 {
			jsonCandidate := strings.TrimSpace(text[start+3 : start+3+end])
			if strings.HasPrefix(jsonCandidate, "{") {
				return jsonCandidate
			}
		}
	}

	// Если не нашли, возвращаем как есть
	return text
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
