package handlers

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"mylittleprice/internal/container"
	"mylittleprice/internal/models"
)

// ChatHandler handles chat-related requests
type ChatHandler struct {
	container *container.Container
}

// NewChatHandler creates a new chat handler with dependency injection
func NewChatHandler(c *container.Container) *ChatHandler {
	return &ChatHandler{
		container: c,
	}
}

func (h *ChatHandler) HandleChat(c *fiber.Ctx) error {
	var req models.ChatRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "invalid_request",
			Message: "Failed to parse request body",
		})
	}

	// ═══════════════════════════════════════════════════════════
	// STEP 1: VALIDATE REQUEST
	// ═══════════════════════════════════════════════════════════

	if req.Message == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "validation_error",
			Message: "Message is required",
		})
	}

	// Set defaults
	if req.Country == "" {
		req.Country = "CH"
	}

	if req.Language == "" {
		req.Language = h.container.SessionService.GetLanguageForCountry(req.Country)
	}

	// ═══════════════════════════════════════════════════════════
	// STEP 2: GET OR CREATE SESSION
	// ═══════════════════════════════════════════════════════════

	var session *models.ChatSession
	var err error

	if req.SessionID != "" {
		session, err = h.container.SessionService.GetSession(req.SessionID)
		if err != nil {
			session, err = h.container.SessionService.CreateSession(req.SessionID, req.Country, req.Language)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
					Error:   "session_error",
					Message: "Failed to create session",
				})
			}
		}
	} else {
		req.SessionID = uuid.New().String()
		session, err = h.container.SessionService.CreateSession(req.SessionID, req.Country, req.Language)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
				Error:   "session_error",
				Message: "Failed to create session",
			})
		}
	}

	// ═══════════════════════════════════════════════════════════
	// STEP 3: HANDLE NEW SEARCH REQUEST
	// ═══════════════════════════════════════════════════════════

	if req.NewSearch {
		fmt.Printf("🔄 New search requested for session %s\n", req.SessionID)

		canStart, message := h.container.SessionService.CanStartNewSearch(req.SessionID)
		if !canStart {
			return c.Status(fiber.StatusTooManyRequests).JSON(models.ErrorResponse{
				Error:   "max_searches_reached",
				Message: message,
			})
		}

		if err := h.container.SessionService.StartNewSearch(req.SessionID); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
				Error:   "session_error",
				Message: "Failed to start new search",
			})
		}

		session, _ = h.container.SessionService.GetSession(req.SessionID)
		fmt.Printf("✅ New search started (count: %d/%d)\n",
			session.SearchState.SearchCount,
			h.container.SessionService.GetMaxSearches())
	} else if session.SearchState.Status == models.SearchStatusIdle {
		session.SearchState.Status = models.SearchStatusInProgress
		h.container.SessionService.UpdateSession(session)
		fmt.Printf("✅ Search status changed: idle → in_progress\n")
	}

	// ═══════════════════════════════════════════════════════════
	// STEP 4: CHECK IF SEARCH IS BLOCKED
	// ═══════════════════════════════════════════════════════════

	if h.container.SessionService.IsSearchCompleted(req.SessionID) {
		fmt.Printf("⛔ Search completed, chat blocked for session %s\n", req.SessionID)

		return c.JSON(models.ChatResponse{
			Type:         "search_blocked",
			Output:       "This search is complete. To search for another product, please click 'New Search' button.",
			SessionID:    req.SessionID,
			MessageCount: session.MessageCount,
			SearchState:  h.container.SessionService.GetSearchStateInfo(req.SessionID),
		})
	}

	// ═══════════════════════════════════════════════════════════
	// STEP 5: CHECK MESSAGE LIMIT
	// ═══════════════════════════════════════════════════════════

	canSend, err := h.container.SessionService.CanSendMessage(req.SessionID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "session_error",
			Message: "Failed to check message limit",
		})
	}

	if !canSend {
		return c.Status(fiber.StatusTooManyRequests).JSON(models.ErrorResponse{
			Error:   "message_limit_exceeded",
			Message: fmt.Sprintf("Maximum %d messages per session reached.", h.container.Config.MaxMessagesPerSession),
		})
	}

	// ═══════════════════════════════════════════════════════════
	// STEP 6: SAVE USER MESSAGE
	// ═══════════════════════════════════════════════════════════

	userMessage := &models.Message{
		ID:        uuid.New(),
		SessionID: session.ID,
		Role:      "user",
		Content:   req.Message,
		Category:  session.SearchState.Category,
		CreatedAt: time.Now(),
	}

	if err := h.container.SessionService.AddMessage(req.SessionID, userMessage); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "message_error",
			Message: "Failed to save message",
		})
	}

	if err := h.container.SessionService.IncrementMessageCount(req.SessionID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "session_error",
			Message: "Failed to update session",
		})
	}

	// ═══════════════════════════════════════════════════════════
	// STEP 7: GET CONVERSATION HISTORY
	// ═══════════════════════════════════════════════════════════

	history, err := h.container.SessionService.GetConversationHistory(req.SessionID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "session_error",
			Message: "Failed to get conversation history",
		})
	}

	// ═══════════════════════════════════════════════════════════
	// STEP 8: PROCESS WITH GEMINI AI
	// ═══════════════════════════════════════════════════════════

	fmt.Printf("🤖 Processing message for session %s\n", req.SessionID)
	if session.SearchState.Category != "" {
		fmt.Printf("   📂 Using existing category: %s\n", session.SearchState.Category)
	} else {
		fmt.Printf("   🆕 First message - will determine category\n")
	}

	startTime := time.Now()
	geminiResponse, keyIndex, err := h.container.GeminiService.ProcessMessageWithContext(
		req.Message,
		history,
		session.CountryCode,
		session.LanguageCode,
		session.SearchState.Category,
	)
	responseTime := time.Since(startTime)

	h.container.GeminiRotator.RecordUsage(keyIndex, err == nil, responseTime)

	if err != nil {
		fmt.Printf("❌ Gemini Error: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "ai_error",
			Message: fmt.Sprintf("Failed to process message with AI: %v", err),
		})
	}

	// ═══════════════════════════════════════════════════════════
	// STEP 9: SAVE CATEGORY IF NEWLY DETERMINED
	// ═══════════════════════════════════════════════════════════

	if geminiResponse.Category != "" && session.SearchState.Category == "" {
		h.container.SessionService.SetCategory(req.SessionID, geminiResponse.Category)
		session.SearchState.Category = geminiResponse.Category
		fmt.Printf("   📂 Category saved: %s\n", geminiResponse.Category)
	}

	// ═══════════════════════════════════════════════════════════
	// STEP 10: PROCESS RESPONSE BASED ON TYPE
	// ═══════════════════════════════════════════════════════════

	var response models.ChatResponse

	// Check if we should force search after too many questions
	if geminiResponse.ResponseType == "dialogue" && session.MessageCount >= 6 {
		fmt.Printf("⚠️  Too many questions (%d), forcing category search\n", session.MessageCount)

		searchPhrase := fmt.Sprintf("%s %s", session.SearchState.Category, strings.Join(session.SearchState.CollectedParams, " "))
		geminiResponse.ResponseType = "search"
		geminiResponse.SearchPhrase = strings.TrimSpace(searchPhrase)
		geminiResponse.SearchType = "category"
	}

	switch geminiResponse.ResponseType {
	case "dialogue":
		response = h.handleDialogueResponse(req, session, geminiResponse)

	case "search":
		searchResponse, err := h.handleSearchResponse(req, session, geminiResponse)
		if err != nil {
			return err
		}
		response = searchResponse

	default:
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "invalid_response",
			Message: "Invalid AI response type",
		})
	}

	// ═══════════════════════════════════════════════════════════
	// STEP 11: SAVE ASSISTANT MESSAGE
	// ═══════════════════════════════════════════════════════════

	assistantMessage := &models.Message{
		ID:           uuid.New(),
		SessionID:    session.ID,
		Role:         "assistant",
		Content:      geminiResponse.Output,
		ResponseType: geminiResponse.ResponseType,
		QuickReplies: geminiResponse.QuickReplies,
		Category:     session.SearchState.Category,
		CreatedAt:    time.Now(),
	}

	if err := h.container.SessionService.AddMessage(req.SessionID, assistantMessage); err != nil {
		fmt.Printf("⚠️  Warning: Failed to save assistant message: %v\n", err)
	}

	return c.JSON(response)
}

// handleDialogueResponse handles dialogue responses
func (h *ChatHandler) handleDialogueResponse(
	req models.ChatRequest,
	session *models.ChatSession,
	geminiResponse *models.GeminiResponse,
) models.ChatResponse {
	if h.isMeaningfulParam(req.Message) {
		h.container.SessionService.AddCollectedParam(req.SessionID, req.Message)
	}

	return models.ChatResponse{
		Type:         "text",
		Output:       geminiResponse.Output,
		QuickReplies: geminiResponse.QuickReplies,
		SessionID:    req.SessionID,
		MessageCount: session.MessageCount + 1,
		SearchState:  h.container.SessionService.GetSearchStateInfo(req.SessionID),
	}
}

// handleSearchResponse handles search responses
func (h *ChatHandler) handleSearchResponse(
	req models.ChatRequest,
	session *models.ChatSession,
	geminiResponse *models.GeminiResponse,
) (models.ChatResponse, error) {
	// Validate query
	isValid, errMsg := h.container.Optimizer.ValidateQuery(geminiResponse.SearchPhrase)
	if !isValid {
		return models.ChatResponse{}, fiber.NewError(fiber.StatusBadRequest, errMsg)
	}

	// Check if query is product-related
	if !h.container.Optimizer.IsProductQuery(geminiResponse.SearchPhrase) {
		return models.ChatResponse{
			Type:         "text",
			Output:       "Please ask about a specific product you'd like to find.",
			SessionID:    req.SessionID,
			MessageCount: session.MessageCount + 1,
			SearchState:  h.container.SessionService.GetSearchStateInfo(req.SessionID),
		}, nil
	}

	// Optimize query
	optimizedQuery := h.container.Optimizer.OptimizeQuery(geminiResponse.SearchPhrase, geminiResponse.SearchType)
	fmt.Printf("   🔍 Search query: '%s' → '%s'\n", geminiResponse.SearchPhrase, optimizedQuery)

	// Search products with caching
	startTime := time.Now()
	products, keyIndex, err := h.container.SerpService.SearchWithCache(
		optimizedQuery,
		geminiResponse.SearchType,
		session.CountryCode,
		h.container.CacheService,
	)
	responseTime := time.Since(startTime)

	if keyIndex != -1 {
		h.container.SerpRotator.RecordUsage(keyIndex, err == nil, responseTime)
		fmt.Printf("   📊 SERP API used (key %d, %dms)\n", keyIndex, responseTime.Milliseconds())
	} else {
		fmt.Printf("   💾 Cache hit (%dms)\n", responseTime.Milliseconds())
	}

	if err != nil {
		return models.ChatResponse{}, fiber.NewError(fiber.StatusInternalServerError, "Failed to search products")
	}

	if len(products) == 0 {
		return models.ChatResponse{
			Type:         "text",
			Output:       "Sorry, I couldn't find any products matching your criteria. Could you try describing what you're looking for differently?",
			SessionID:    req.SessionID,
			MessageCount: session.MessageCount + 1,
			SearchState:  h.container.SessionService.GetSearchStateInfo(req.SessionID),
		}, nil
	}

	// Mark search as completed
	h.container.SessionService.MarkSearchCompleted(req.SessionID)
	fmt.Printf("   ✅ Search completed! Found %d products. Chat blocked.\n", len(products))

	output := "Here are the best options I found for you:"
	if geminiResponse.SearchType == "exact" {
		output = "Here is the product you were looking for:"
	}

	return models.ChatResponse{
		Type:         "product_card",
		Output:       output,
		Products:     products,
		SearchType:   geminiResponse.SearchType,
		SessionID:    req.SessionID,
		MessageCount: session.MessageCount + 1,
		SearchState:  h.container.SessionService.GetSearchStateInfo(req.SessionID),
	}, nil
}

// isMeaningfulParam checks if user message is a meaningful parameter
func (h *ChatHandler) isMeaningfulParam(message string) bool {
	message = strings.TrimSpace(strings.ToLower(message))

	skipWords := []string{
		"hello", "hi", "hey", "yes", "no", "ok", "okay",
		"thanks", "thank you", "bye", "goodbye",
	}

	for _, word := range skipWords {
		if message == word {
			return false
		}
	}

	return len(message) > 2
}
