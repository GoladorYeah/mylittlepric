package handlers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"mylittleprice/internal/container"
	"mylittleprice/internal/models"
)

type ChatHandler struct {
	container *container.Container
}

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

		if err := h.container.SessionService.StartNewSearch(req.SessionID); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
				Error:   "session_error",
				Message: "Failed to start new search",
			})
		}

		session, _ = h.container.SessionService.GetSession(req.SessionID)
		fmt.Printf("✅ New search started (%d)\n", session.SearchState.SearchCount)
	} else if session.SearchState.Status == models.SearchStatusIdle {
		session.SearchState.Status = models.SearchStatusInProgress
		h.container.SessionService.UpdateSession(session)
	}

	if h.container.SessionService.IsSearchCompleted(req.SessionID) {
		messageLower := strings.ToLower(req.Message)
		refinementWords := []string{"cheap", "expensive", "premium", "budget", "price", "cost", "less", "more"}
		isRefinement := false
		for _, word := range refinementWords {
			if strings.Contains(messageLower, word) {
				isRefinement = true
				break
			}
		}

		if isRefinement {
			h.container.SessionService.ResetSearchStatus(req.SessionID)
			session.SearchState.Status = models.SearchStatusInProgress
			h.container.SessionService.UpdateSession(session)
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
		session.SearchState.LastProduct,
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
		fmt.Printf("   📂 Category set: %s\n\n", geminiResponse.Category)
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

	if geminiResponse.PriceFilter != "" && session.SearchState.LastProduct != nil {
		optimizedQuery = session.SearchState.LastProduct.Name
	}

	fmt.Printf("   🔍 Query: '%s' → '%s'\n", geminiResponse.SearchPhrase, optimizedQuery)
	if geminiResponse.PriceFilter != "" {
		fmt.Printf("   💰 Price filter: %s (ref: %.2f)\n", geminiResponse.PriceFilter, session.SearchState.LastProduct.Price)
	}

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

	if geminiResponse.PriceFilter != "" && session.SearchState.LastProduct != nil {
		products = h.filterByPrice(products, session.SearchState.LastProduct.Price, geminiResponse.PriceFilter, session.Currency)
	}

	if len(products) > 0 {
		firstProduct := products[0]
		price := h.extractPrice(firstProduct.Price)
		h.container.SessionService.SetLastProduct(req.SessionID, firstProduct.Name, price)
	}

	h.container.SessionService.MarkSearchCompleted(req.SessionID)
	fmt.Printf("   ✅ Found %d products.\n", len(products))

	output := "Here are the best options:"
	if geminiResponse.SearchType == "exact" {
		output = "Here is the product:"
	}
	if geminiResponse.PriceFilter == "cheaper" {
		output = "Here are cheaper alternatives:"
	} else if geminiResponse.PriceFilter == "expensive" {
		output = "Here are more premium options:"
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

func (h *ChatHandler) filterByPrice(products []models.ProductCard, refPrice float64, filter string, currency string) []models.ProductCard {
	filtered := []models.ProductCard{}

	for _, product := range products {
		price := h.extractPrice(product.Price)

		if filter == "cheaper" && price < refPrice && price > refPrice*0.3 {
			filtered = append(filtered, product)
		} else if filter == "expensive" && price > refPrice && price < refPrice*3 {
			filtered = append(filtered, product)
		}
	}

	if len(filtered) == 0 {
		return products
	}

	return filtered
}

func (h *ChatHandler) extractPrice(priceStr string) float64 {
	cleaned := strings.TrimSpace(priceStr)
	cleaned = strings.ReplaceAll(cleaned, "CHF", "")
	cleaned = strings.ReplaceAll(cleaned, "EUR", "")
	cleaned = strings.ReplaceAll(cleaned, "USD", "")
	cleaned = strings.ReplaceAll(cleaned, "GBP", "")
	cleaned = strings.ReplaceAll(cleaned, "$", "")
	cleaned = strings.ReplaceAll(cleaned, "€", "")
	cleaned = strings.ReplaceAll(cleaned, "£", "")
	cleaned = strings.ReplaceAll(cleaned, ",", "")
	cleaned = strings.TrimSpace(cleaned)

	price, err := strconv.ParseFloat(cleaned, 64)
	if err != nil {
		return 0
	}

	return price
}
