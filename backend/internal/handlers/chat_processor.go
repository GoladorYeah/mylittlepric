package handlers

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	"mylittleprice/internal/container"
	"mylittleprice/internal/models"
)

// ChatProcessor handles the core chat processing logic shared between REST and WebSocket handlers
type ChatProcessor struct {
	container *container.Container
}

// NewChatProcessor creates a new chat processor
func NewChatProcessor(c *container.Container) *ChatProcessor {
	return &ChatProcessor{
		container: c,
	}
}

// ChatRequest represents a standardized chat request
type ChatRequest struct {
	SessionID       string
	Message         string
	Country         string
	Language        string
	NewSearch       bool
	CurrentCategory string
}

// ChatProcessorResponse represents the standardized response from chat processing
type ChatProcessorResponse struct {
	Type            string
	Output          string
	QuickReplies    []string
	Products        []models.ProductCard
	SearchType      string
	SessionID       string
	MessageCount    int
	SearchState     *models.SearchStateResponse
	Error           *ErrorInfo
}

// ErrorInfo contains error details
type ErrorInfo struct {
	Code    string
	Message string
}

// ProcessChat handles the main chat processing logic
func (p *ChatProcessor) ProcessChat(req *ChatRequest) *ChatProcessorResponse {
	// Validate input
	if req.Message == "" {
		return &ChatProcessorResponse{
			Error: &ErrorInfo{
				Code:    "validation_error",
				Message: "Message is required",
			},
		}
	}

	// Set defaults
	if req.Country == "" {
		req.Country = p.container.Config.DefaultCountry
	}
	if req.Language == "" {
		req.Language = p.container.Config.DefaultLanguage
	}

	// Get or create session
	session, err := p.getOrCreateSession(req)
	if err != nil {
		return &ChatProcessorResponse{
			Error: &ErrorInfo{
				Code:    "session_error",
				Message: "Failed to create session",
			},
		}
	}

	// Handle new search
	if req.NewSearch {
		fmt.Printf("🔄 New search for session %s\n", req.SessionID)
		if err := p.container.SessionService.StartNewSearch(req.SessionID); err != nil {
			fmt.Printf("⚠️ Failed to start new search: %v\n", err)
		}
		session, _ = p.container.SessionService.GetSession(req.SessionID)
	}

	// Handle category update
	if req.CurrentCategory != "" && req.CurrentCategory != session.SearchState.Category {
		session.SearchState.Category = req.CurrentCategory
		p.container.SessionService.UpdateSession(session)
	}

	// Check search limit
	if session.SearchState.SearchCount >= p.container.SessionService.GetMaxSearches() {
		return &ChatProcessorResponse{
			Type:         "text",
			Output:       "You have reached the maximum number of searches. Please start a new search.",
			SessionID:    req.SessionID,
			MessageCount: session.MessageCount,
			SearchState: &models.SearchStateResponse{
				Status:      string(session.SearchState.Status),
				Category:    session.SearchState.Category,
				CanContinue: false,
				SearchCount: session.SearchState.SearchCount,
				MaxSearches: p.container.SessionService.GetMaxSearches(),
				Message:     "Search limit reached",
			},
		}
	}

	// Store user message
	userMessage := &models.Message{
		ID:        uuid.New(),
		SessionID: session.ID,
		Role:      "user",
		Content:   req.Message,
		CreatedAt: time.Now(),
	}

	if err := p.container.SessionService.AddMessage(req.SessionID, userMessage); err != nil {
		return &ChatProcessorResponse{
			Error: &ErrorInfo{
				Code:    "storage_error",
				Message: "Failed to store message",
			},
		}
	}

	if err := p.container.SessionService.IncrementMessageCount(req.SessionID); err != nil {
		fmt.Printf("⚠️ Failed to increment message count: %v\n", err)
	}

	// Get conversation history
	conversationHistory, err := p.container.SessionService.GetConversationHistory(req.SessionID)
	if err != nil {
		conversationHistory = []map[string]string{}
	}

	// Process with Gemini
	geminiResponse, _, err := p.container.GeminiService.ProcessMessageWithContext(
		req.Message,
		conversationHistory,
		req.Country,
		req.Language,
		session.SearchState.Category,
		session.SearchState.LastProduct,
	)

	if err != nil {
		log.Printf("❌ Gemini processing error: %v", err)
		return &ChatProcessorResponse{
			Error: &ErrorInfo{
				Code:    "processing_error",
				Message: fmt.Sprintf("AI processing failed: %v", err),
			},
		}
	}

	if geminiResponse == nil {
		log.Printf("❌ Gemini returned nil response")
		return &ChatProcessorResponse{
			Error: &ErrorInfo{
				Code:    "processing_error",
				Message: "AI returned empty response",
			},
		}
	}

	// Update category
	if geminiResponse.Category != "" {
		session.SearchState.Category = geminiResponse.Category
	}

	// Store assistant message
	assistantMessage := &models.Message{
		ID:           uuid.New(),
		SessionID:    session.ID,
		Role:         "assistant",
		Content:      geminiResponse.Output,
		ResponseType: geminiResponse.ResponseType,
		QuickReplies: geminiResponse.QuickReplies,
		CreatedAt:    time.Now(),
	}

	if err := p.container.SessionService.AddMessage(req.SessionID, assistantMessage); err != nil {
		fmt.Printf("⚠️ Failed to store assistant message: %v\n", err)
	}

	// Build response
	response := &ChatProcessorResponse{
		Type:         geminiResponse.ResponseType,
		Output:       geminiResponse.Output,
		QuickReplies: geminiResponse.QuickReplies,
		SessionID:    req.SessionID,
		MessageCount: session.MessageCount + 1,
	}

	// Handle search
	if geminiResponse.ResponseType == "search" {
		products, searchErr := p.performSearch(geminiResponse, req.Country, req.Language)
		if searchErr != nil {
			log.Printf("⚠️ Search failed: %v", searchErr)
			response.Output = "Sorry, I couldn't find any products. Please try different keywords."
			response.Type = "text"
		} else if len(products) > 0 {
			response.Products = products
			response.SearchType = geminiResponse.SearchType

			// Update last product
			if len(products) > 0 {
				price := parsePrice(products[0].Price)
				session.SearchState.LastProduct = &models.ProductInfo{
					Name:  products[0].Name,
					Price: price,
				}
			}

			session.SearchState.SearchCount++
			assistantMessage.Products = products
		}
	}

	// Update session state
	session.SearchState.Status = models.SearchStatusIdle
	if err := p.container.SessionService.UpdateSession(session); err != nil {
		fmt.Printf("⚠️ Failed to update session: %v\n", err)
	}

	// Build search state response
	response.SearchState = &models.SearchStateResponse{
		Status:      string(session.SearchState.Status),
		Category:    session.SearchState.Category,
		CanContinue: session.SearchState.SearchCount < p.container.SessionService.GetMaxSearches(),
		SearchCount: session.SearchState.SearchCount,
		MaxSearches: p.container.SessionService.GetMaxSearches(),
	}

	return response
}

// getOrCreateSession handles session retrieval or creation
func (p *ChatProcessor) getOrCreateSession(req *ChatRequest) (*models.ChatSession, error) {
	var session *models.ChatSession
	var err error

	if req.SessionID != "" {
		session, err = p.container.SessionService.GetSession(req.SessionID)
		if err != nil {
			req.SessionID = uuid.New().String()
			session, err = p.container.SessionService.CreateSession(req.SessionID, req.Country, req.Language)
			if err != nil {
				return nil, err
			}
		}
	} else {
		req.SessionID = uuid.New().String()
		session, err = p.container.SessionService.CreateSession(req.SessionID, req.Country, req.Language)
		if err != nil {
			return nil, err
		}
	}

	return session, nil
}

// performSearch executes product search with translation
func (p *ChatProcessor) performSearch(geminiResp *models.GeminiResponse, country, language string) ([]models.ProductCard, error) {
	// Translate query to English
	translatedQuery, err := p.container.GeminiService.TranslateToEnglish(geminiResp.SearchPhrase)
	if err != nil {
		fmt.Printf("⚠️ Translation failed, using original query: %v\n", err)
		translatedQuery = geminiResp.SearchPhrase
	} else if translatedQuery != geminiResp.SearchPhrase {
		fmt.Printf("🌐 Translated: '%s' → '%s'\n", geminiResp.SearchPhrase, translatedQuery)
	}

	products, _, err := p.container.SerpService.SearchWithCache(
		translatedQuery,
		geminiResp.SearchType,
		country,
		p.container.CacheService,
	)

	if err != nil {
		return nil, err
	}

	return products, nil
}

// parsePrice extracts numeric price from price string
func parsePrice(priceStr string) float64 {
	priceStr = strings.ReplaceAll(priceStr, "$", "")
	priceStr = strings.ReplaceAll(priceStr, "€", "")
	priceStr = strings.ReplaceAll(priceStr, "£", "")
	priceStr = strings.ReplaceAll(priceStr, "CHF", "")
	priceStr = strings.TrimSpace(priceStr)
	priceStr = strings.ReplaceAll(priceStr, ",", "")

	price, _ := strconv.ParseFloat(priceStr, 64)
	return price
}
