package handlers

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	"mylittleprice/internal/container"
	"mylittleprice/internal/models"
	"mylittleprice/internal/services"
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
	UserID          *uuid.UUID // Optional user ID for authenticated users
	Message         string
	Country         string
	Language        string
	Currency        string
	NewSearch       bool
	CurrentCategory string
}

// ChatProcessorResponse represents the standardized response from chat processing
type ChatProcessorResponse struct {
	Type         string
	Output       string
	QuickReplies []string
	Products     []models.ProductCard
	SearchType   string
	SessionID    string
	MessageCount int
	SearchState  *models.SearchStateResponse
	Error        *ErrorInfo
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
	if req.Currency == "" {
		req.Currency = p.container.Config.DefaultCurrency
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

	// Add user message to cycle history
	if err := p.container.SessionService.AddToCycleHistory(req.SessionID, "user", req.Message); err != nil {
		fmt.Printf("⚠️ Failed to add to cycle history: %v\n", err)
	}

	// Re-fetch session after updating cycle history
	session, err = p.container.SessionService.GetSession(req.SessionID)
	if err != nil {
		return &ChatProcessorResponse{
			Error: &ErrorInfo{
				Code:    "session_error",
				Message: "Failed to get session after update",
			},
		}
	}

	// Process with Universal Prompt System with retry logic
	var geminiResponse *models.GeminiResponse
	var geminiErr error
	const maxProcessingRetries = 2

	for attempt := 0; attempt <= maxProcessingRetries; attempt++ {
		if attempt > 0 {
			log.Printf("🔄 Retry processing attempt %d/%d", attempt+1, maxProcessingRetries+1)
		}

		geminiResponse, geminiErr = p.container.GeminiService.ProcessWithUniversalPrompt(
			req.Message,
			session,
		)

		// Success - break out of retry loop
		if geminiErr == nil && geminiResponse != nil {
			if attempt > 0 {
				log.Printf("✅ Processing succeeded on retry attempt %d", attempt+1)
			}
			break
		}

		// Log the error
		if geminiErr != nil {
			log.Printf("❌ Gemini processing error (attempt %d/%d): %v", attempt+1, maxProcessingRetries+1, geminiErr)
		}

		// If this is the last attempt, use fallback response
		if attempt == maxProcessingRetries {
			log.Printf("⚠️ All processing attempts failed, using fallback response")

			// Return helpful fallback response instead of error
			return &ChatProcessorResponse{
				Type:         "dialogue",
				Output:       "I'm having trouble processing your request right now. Could you please rephrase your question or try again in a moment?",
				QuickReplies: []string{"Start over", "Try again"},
				SessionID:    req.SessionID,
				MessageCount: session.MessageCount,
				SearchState: &models.SearchStateResponse{
					Status:      string(session.SearchState.Status),
					Category:    session.SearchState.Category,
					CanContinue: session.SearchState.SearchCount < p.container.SessionService.GetMaxSearches(),
					SearchCount: session.SearchState.SearchCount,
					MaxSearches: p.container.SessionService.GetMaxSearches(),
					Message:     "Temporary processing issue",
				},
			}
		}

		// Wait a bit before retry (500ms, 1s)
		if attempt < maxProcessingRetries {
			retryDelay := time.Duration(500*(attempt+1)) * time.Millisecond
			time.Sleep(retryDelay)
		}
	}

	// Log the Gemini response for debugging
	priceInfo := ""
	if geminiResponse.MinPrice != nil || geminiResponse.MaxPrice != nil {
		priceInfo = fmt.Sprintf(", price_range=%v-%v", geminiResponse.MinPrice, geminiResponse.MaxPrice)
	}
	fmt.Printf("📥 Gemini response: type=%s, category=%s, search_phrase=%s%s\n",
		geminiResponse.ResponseType, geminiResponse.Category, geminiResponse.SearchPhrase, priceInfo)

	// Update category
	if geminiResponse.Category != "" {
		session.SearchState.Category = geminiResponse.Category
	}

	// Create assistant message (but don't save yet - we may need to add products first)
	assistantMessage := &models.Message{
		ID:           uuid.New(),
		SessionID:    session.ID,
		Role:         "assistant",
		Content:      geminiResponse.Output,
		ResponseType: geminiResponse.ResponseType,
		QuickReplies: geminiResponse.QuickReplies,
		CreatedAt:    time.Now(),
	}

	// Build response
	response := &ChatProcessorResponse{
		Type:         geminiResponse.ResponseType,
		Output:       geminiResponse.Output,
		QuickReplies: geminiResponse.QuickReplies,
		SessionID:    req.SessionID,
		MessageCount: session.MessageCount + 1,
	}

	// Handle search (intermediate search for verification/grounding)
	if geminiResponse.ResponseType == "search" {
		priceRangeStr := ""
		if geminiResponse.MinPrice != nil || geminiResponse.MaxPrice != nil {
			priceRangeStr = fmt.Sprintf(", price_range=%v-%v", geminiResponse.MinPrice, geminiResponse.MaxPrice)
		}
		fmt.Printf("�� Search request detected: phrase='%s', type='%s'%s\n",
			geminiResponse.SearchPhrase, geminiResponse.SearchType, priceRangeStr)

		// Validate search phrase
		if geminiResponse.SearchPhrase == "" {
			fmt.Printf("⚠️ Empty search phrase in search request\n")
			response.Output = "I need more details about what product you're looking for. Could you be more specific?"
			response.Type = "dialogue"
		} else {
			products, translatedQuery, searchErr := p.performSearch(geminiResponse, req.Country, req.Language)
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
				// Add products to assistant message BEFORE saving
				assistantMessage.Products = products

				// NEW: Update last search in conversation context
				productInfoList := make([]models.ProductInfo, 0, len(products))
				for _, p := range products {
					price := parsePrice(p.Price)
					productInfoList = append(productInfoList, models.ProductInfo{
						Name:  p.Name,
						Price: price,
					})
				}
				contextExtractor := p.container.GeminiService.GetContextExtractor()
				contextExtractor.UpdateLastSearch(session, translatedQuery, geminiResponse.Category, productInfoList, "")

				// Save search history
				p.saveSearchHistory(req, session, geminiResponse, translatedQuery, products)
			}
		}
	}

	// Handle api_request (final product search to complete cycle)
	if geminiResponse.ResponseType == "api_request" {
		fmt.Printf("🎯 API request detected: api='%s', params=%+v\n",
			geminiResponse.API, geminiResponse.Params)

		if geminiResponse.API == "google_shopping" && geminiResponse.Params != nil {
			// Extract query from params
			query, ok := geminiResponse.Params["q"].(string)
			if !ok || query == "" {
				fmt.Printf("⚠️ Missing or invalid 'q' parameter in api_request\n")
				response.Output = "I need more details about what product you're looking for. Could you be more specific?"
				response.Type = "dialogue"
			} else {
				// Perform the final search
				fmt.Printf("🛍️ Final product search: '%s'\n", query)

				// Create a temporary GeminiResponse for search
				searchResp := &models.GeminiResponse{
					SearchPhrase: query,
					SearchType:   "exact", // Default to exact for final searches
					Category:     geminiResponse.Category,
					PriceFilter:  geminiResponse.PriceFilter,
				}

				products, translatedQuery, searchErr := p.performSearch(searchResp, req.Country, req.Language)
				if searchErr != nil {
					log.Printf("⚠️ Final search failed: %v", searchErr)
					response.Output = "Sorry, I couldn't find any products. Please try different keywords."
					response.Type = "text"
				} else if len(products) > 0 {
					response.Products = products
					response.SearchType = "exact"
					response.Output = geminiResponse.Output // Use AI's message if provided

					// Update last product
					if len(products) > 0 {
						price := parsePrice(products[0].Price)
						session.SearchState.LastProduct = &models.ProductInfo{
							Name:  products[0].Name,
							Price: price,
						}
					}

					session.SearchState.SearchCount++
					// Add products to assistant message BEFORE saving
					assistantMessage.Products = products

					// NEW: Update last search in conversation context
					productInfoList := make([]models.ProductInfo, 0, len(products))
					for _, p := range products {
						price := parsePrice(p.Price)
						productInfoList = append(productInfoList, models.ProductInfo{
							Name:  p.Name,
							Price: price,
						})
					}
					contextExtractor := p.container.GeminiService.GetContextExtractor()
					contextExtractor.UpdateLastSearch(session, translatedQuery, searchResp.Category, productInfoList, "")

					// Save search history
					p.saveSearchHistory(req, session, searchResp, translatedQuery, products)

					fmt.Printf("✅ Cycle completed with %d products\n", len(products))
				} else {
					response.Output = "I couldn't find that exact product. Would you like to see similar alternatives?"
					response.Type = "dialogue"
				}
			}
		} else {
			fmt.Printf("⚠️ Unsupported API: %s\n", geminiResponse.API)
			response.Output = "I encountered an error processing your request. Please try again."
			response.Type = "dialogue"
		}
	}

	// Save assistant message (now with products if it was a search)
	if err := p.container.SessionService.AddMessage(req.SessionID, assistantMessage); err != nil {
		fmt.Printf("⚠️ Failed to store assistant message: %v\n", err)
	}

	// Add assistant response to cycle history
	if err := p.container.SessionService.AddToCycleHistory(req.SessionID, "assistant", geminiResponse.Output); err != nil {
		fmt.Printf("⚠️ Failed to add assistant response to cycle history: %v\n", err)
	}

	// Re-fetch session to get updated cycle history
	session, err = p.container.SessionService.GetSession(req.SessionID)
	if err != nil {
		fmt.Printf("⚠️ Failed to re-fetch session after adding to history: %v\n", err)
	}

	// NEW: Update conversation context periodically
	contextExtractor := p.container.GeminiService.GetContextExtractor()
	contextOptimizer := p.container.GeminiService.GetContextOptimizer()

	if contextOptimizer.ShouldUpdateContext(session) {
		fmt.Printf("🧠 Updating conversation context...\n")
		if err := contextExtractor.UpdateConversationContext(session, session.CycleState.CycleHistory); err != nil {
			fmt.Printf("⚠️ Failed to update conversation context: %v\n", err)
		} else {
			// Save updated context
			if err := p.container.SessionService.UpdateSession(session); err != nil {
				fmt.Printf("⚠️ Failed to save updated context: %v\n", err)
			}
		}
	}

	// Check if we need to start a new cycle (iteration limit reached)
	// This checks BEFORE incrementing, so iteration 6 will trigger a new cycle
	shouldStartNewCycle, err := p.container.SessionService.IncrementCycleIteration(req.SessionID)
	if err != nil {
		fmt.Printf("⚠️ Failed to increment cycle iteration: %v\n", err)
	}

	if shouldStartNewCycle {
		fmt.Printf("🔄 Iteration limit reached (%d), starting new cycle\n", services.MaxIterations)

		// Collect products from last cycle
		products := []models.ProductInfo{}
		if session.SearchState.LastProduct != nil {
			products = append(products, *session.SearchState.LastProduct)
		}

		// Start new cycle with context carryover
		if err := p.container.SessionService.StartNewCycle(req.SessionID, req.Message, products); err != nil {
			fmt.Printf("⚠️ Failed to start new cycle: %v\n", err)
		}
	}

	// Re-fetch session after cycle operations (either increment or new cycle)
	session, err = p.container.SessionService.GetSession(req.SessionID)
	if err != nil {
		fmt.Printf("⚠️ Failed to re-fetch session after cycle operations: %v\n", err)
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
		// Try to get existing session
		session, err = p.container.SessionService.GetSession(req.SessionID)
		if err != nil {
			// Session not found - create new one with the SAME ID
			fmt.Printf("⚠️ Session %s not found in Redis, creating new session with same ID\n", req.SessionID)
			session, err = p.container.SessionService.CreateSession(req.SessionID, req.Country, req.Language, req.Currency)
			if err != nil {
				return nil, err
			}
		} else {
			// Session exists - preserve language, country, and currency from session
			// Only update if explicitly changed by user (non-empty AND different from session)
			updated := false

			// Update language if changed and not empty
			if req.Language != "" && req.Language != session.LanguageCode {
				fmt.Printf("🗣️ Updating session language from %s to %s\n", session.LanguageCode, req.Language)
				session.LanguageCode = req.Language
				updated = true
			}

			// Update currency if changed and not empty
			if req.Currency != "" && req.Currency != session.Currency {
				fmt.Printf("💱 Updating session currency from %s to %s\n", session.Currency, req.Currency)
				session.Currency = req.Currency
				updated = true
			}

			// Update country if changed and not empty
			if req.Country != "" && req.Country != session.CountryCode {
				fmt.Printf("🌍 Updating session country from %s to %s\n", session.CountryCode, req.Country)
				session.CountryCode = req.Country
				updated = true
			}

			if updated {
				p.container.SessionService.UpdateSession(session)
			}

			// Override request with session values to ensure consistency
			req.Language = session.LanguageCode
			req.Country = session.CountryCode
			req.Currency = session.Currency
		}
	} else {
		// No session ID provided - generate new one
		req.SessionID = uuid.New().String()
		session, err = p.container.SessionService.CreateSession(req.SessionID, req.Country, req.Language, req.Currency)
		if err != nil {
			return nil, err
		}
	}

	return session, nil
}

// performSearch executes product search with translation
func (p *ChatProcessor) performSearch(geminiResp *models.GeminiResponse, country, language string) ([]models.ProductCard, string, error) {
	// Translate query to English for better search results
	fmt.Printf("🔤 Translation check: '%s'\n", geminiResp.SearchPhrase)

	translatedQuery, err := p.container.GeminiService.TranslateToEnglish(geminiResp.SearchPhrase)
	if err != nil {
		fmt.Printf("⚠️ Translation failed: %v, using original query\n", err)
		translatedQuery = geminiResp.SearchPhrase
	} else if translatedQuery != geminiResp.SearchPhrase {
		fmt.Printf("🌐 Translated: '%s' → '%s'\n", geminiResp.SearchPhrase, translatedQuery)
	} else {
		fmt.Printf("✓ Query already in English: '%s'\n", translatedQuery)
	}

	// Log price range if provided
	if geminiResp.MinPrice != nil || geminiResp.MaxPrice != nil {
		fmt.Printf("💰 Price range: %v - %v\n", geminiResp.MinPrice, geminiResp.MaxPrice)
	}

	fmt.Printf("📤 Sending to SERP: '%s'\n", translatedQuery)

	products, _, err := p.container.SerpService.SearchWithCache(
		translatedQuery,
		geminiResp.SearchType,
		country,
		geminiResp.MinPrice,
		geminiResp.MaxPrice,
		p.container.CacheService,
	)

	if err != nil {
		return nil, translatedQuery, err
	}

	return products, translatedQuery, nil
}

// saveSearchHistory saves the search to history
func (p *ChatProcessor) saveSearchHistory(req *ChatRequest, session *models.ChatSession, geminiResp *models.GeminiResponse, translatedQuery string, products []models.ProductCard) {
	// Set currency from request or use default
	currency := req.Currency
	if currency == "" {
		currency = session.Currency
	}

	// Use session ID as string (no parsing needed)
	var sessionIDStr *string
	if req.SessionID != "" {
		sessionIDStr = &req.SessionID
	}

	history := &models.SearchHistory{
		UserID:         req.UserID,
		SessionID:      sessionIDStr,
		SearchQuery:    geminiResp.SearchPhrase,
		OptimizedQuery: &translatedQuery,
		SearchType:     geminiResp.SearchType,
		Category:       &geminiResp.Category,
		CountryCode:    req.Country,
		LanguageCode:   req.Language,
		Currency:       currency,
		ResultCount:    len(products),
		ProductsFound:  products,
	}

	// Save asynchronously to avoid blocking
	go func() {
		ctx := context.Background()
		if err := p.container.SearchHistoryService.SaveSearchHistory(ctx, history); err != nil {
			log.Printf("⚠️ Failed to save search history: %v", err)
		} else {
			log.Printf("📜 Search history saved: '%s' (%d results)", geminiResp.SearchPhrase, len(products))
		}
	}()
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
