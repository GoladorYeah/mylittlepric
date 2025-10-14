package handlers

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"mylittleprice/internal/container"
	"mylittleprice/internal/models"
	"mylittleprice/internal/services"
)

type ChatHandler struct {
	container        *container.Container
	categoryDetector *services.CategoryDetector
}

func NewChatHandler(c *container.Container) *ChatHandler {
	return &ChatHandler{
		container:        c,
		categoryDetector: services.NewCategoryDetector(),
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

	if req.Message == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "validation_error",
			Message: "Message is required",
		})
	}

	if req.Country == "" {
		req.Country = "CH"
	}

	if req.Language == "" {
		req.Language = h.container.SessionService.GetLanguageForCountry(req.Country)
	}

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

	if req.NewSearch {
		fmt.Printf("🔄 New search for session %s\n", req.SessionID)

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
		fmt.Printf("✅ New search started (%d/%d)\n",
			session.SearchState.SearchCount,
			h.container.SessionService.GetMaxSearches())
	} else if session.SearchState.Status == models.SearchStatusIdle {
		session.SearchState.Status = models.SearchStatusInProgress
		h.container.SessionService.UpdateSession(session)
	}

	if h.container.SessionService.IsSearchCompleted(req.SessionID) {
		fmt.Printf("⛔ Search completed, blocked\n")

		return c.JSON(models.ChatResponse{
			Type:         "search_blocked",
			Output:       "Search complete. Click 'New Search' to continue.",
			SessionID:    req.SessionID,
			MessageCount: session.MessageCount,
			SearchState:  h.container.SessionService.GetSearchStateInfo(req.SessionID),
		})
	}

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
			Message: fmt.Sprintf("Max %d messages per session.", h.container.Config.MaxMessagesPerSession),
		})
	}

	if session.SearchState.Category == "" {
		detectedCategory := h.categoryDetector.DetectCategory(req.Message)
		if detectedCategory != "" {
			session.SearchState.Category = detectedCategory
			h.container.SessionService.SetCategory(req.SessionID, detectedCategory)
			fmt.Printf("   🎯 Auto-detected category: %s\n", detectedCategory)
		}
	}

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

	history, err := h.container.SessionService.GetConversationHistory(req.SessionID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "session_error",
			Message: "Failed to get history",
		})
	}

	fmt.Printf("🤖 Processing (category: %s)\n", session.SearchState.Category)

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
		fmt.Printf("❌ Gemini error: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "ai_error",
			Message: fmt.Sprintf("AI error: %v", err),
		})
	}

	if geminiResponse.Category != "" && session.SearchState.Category == "" {
		h.container.SessionService.SetCategory(req.SessionID, geminiResponse.Category)
		session.SearchState.Category = geminiResponse.Category
		fmt.Printf("   📂 Category set: %s\n", geminiResponse.Category)
	}

	var response models.ChatResponse

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
			Message: "Invalid response type",
		})
	}

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

func (h *ChatHandler) handleDialogueResponse(
	req models.ChatRequest,
	session *models.ChatSession,
	geminiResponse *models.GeminiResponse,
) models.ChatResponse {

	if h.isMeaningfulParam(req.Message) {
		h.container.SessionService.AddCollectedParam(req.SessionID, req.Message)
	}

	if session.MessageCount >= 8 {
		fmt.Printf("⚠️  Too many questions (%d), forcing search\n", session.MessageCount)

		params := strings.Join(session.SearchState.CollectedParams, " ")
		searchPhrase := fmt.Sprintf("%s %s", session.SearchState.Category, params)
		searchPhrase = strings.TrimSpace(searchPhrase)

		geminiResponse.ResponseType = "search"
		geminiResponse.SearchPhrase = searchPhrase
		geminiResponse.SearchType = "parameters"

		searchResponse, err := h.handleSearchResponse(req, session, geminiResponse)
		if err != nil {
			return models.ChatResponse{
				Type:         "text",
				Output:       geminiResponse.Output,
				QuickReplies: geminiResponse.QuickReplies,
				SessionID:    req.SessionID,
				MessageCount: session.MessageCount + 1,
				SearchState:  h.container.SessionService.GetSearchStateInfo(req.SessionID),
			}
		}
		return searchResponse
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

func (h *ChatHandler) handleSearchResponse(
	req models.ChatRequest,
	session *models.ChatSession,
	geminiResponse *models.GeminiResponse,
) (models.ChatResponse, error) {

	isValid, errMsg := h.container.Optimizer.ValidateQuery(geminiResponse.SearchPhrase)
	if !isValid {
		return models.ChatResponse{}, fiber.NewError(fiber.StatusBadRequest, errMsg)
	}

	if !h.container.Optimizer.IsProductQuery(geminiResponse.SearchPhrase) {
		return models.ChatResponse{
			Type:         "text",
			Output:       "Please ask about a specific product.",
			SessionID:    req.SessionID,
			MessageCount: session.MessageCount + 1,
			SearchState:  h.container.SessionService.GetSearchStateInfo(req.SessionID),
		}, nil
	}

	optimizedQuery := h.container.Optimizer.OptimizeQuery(geminiResponse.SearchPhrase, geminiResponse.SearchType)
	fmt.Printf("   🔍 Query: '%s' → '%s'\n", geminiResponse.SearchPhrase, optimizedQuery)

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
		fmt.Printf("   📊 SERP API (key %d, %dms)\n", keyIndex, responseTime.Milliseconds())
	} else {
		fmt.Printf("   💾 Cache hit (%dms)\n", responseTime.Milliseconds())
	}

	if err != nil {
		return models.ChatResponse{}, fiber.NewError(fiber.StatusInternalServerError, "Search failed")
	}

	if len(products) == 0 {
		return models.ChatResponse{
			Type:         "text",
			Output:       "No products found. Try different description.",
			SessionID:    req.SessionID,
			MessageCount: session.MessageCount + 1,
			SearchState:  h.container.SessionService.GetSearchStateInfo(req.SessionID),
		}, nil
	}

	h.container.SessionService.MarkSearchCompleted(req.SessionID)
	fmt.Printf("   ✅ Found %d products. Chat blocked.\n", len(products))

	output := "Here are the best options:"
	if geminiResponse.SearchType == "exact" {
		output = "Here is the product:"
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
