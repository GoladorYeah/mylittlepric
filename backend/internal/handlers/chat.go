// backend/internal/handlers/chat.go
package handlers

import (
	"fmt"
	"log"
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
		req.Language = "en"
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
			fmt.Printf("⚠️ Failed to start new search: %v\n", err)
		}
	}

	if session.SearchState.SearchCount >= h.container.SessionService.GetMaxSearches() {
		return c.Status(fiber.StatusOK).JSON(models.ChatResponse{
			Type:         "text",
			Output:       "You have reached the maximum number of searches. Please start a new search.",
			SessionID:    req.SessionID,
			MessageCount: session.MessageCount,
			SearchState: &models.SearchStateResponse{
				Status:      string(session.SearchState.Status),
				CanContinue: false,
				SearchCount: session.SearchState.SearchCount,
				MaxSearches: h.container.SessionService.GetMaxSearches(),
				Message:     "Search limit reached",
			},
		})
	}

	userMessage := &models.Message{
		ID:        uuid.New(),
		SessionID: session.ID,
		Role:      "user",
		Content:   req.Message,
		CreatedAt: time.Now(),
	}

	if err := h.container.SessionService.AddMessage(req.SessionID, userMessage); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "storage_error",
			Message: "Failed to store message",
		})
	}

	if err := h.container.SessionService.IncrementMessageCount(req.SessionID); err != nil {
		fmt.Printf("⚠️ Failed to increment message count: %v\n", err)
	}

	conversationHistory, err := h.container.SessionService.GetConversationHistory(req.SessionID)
	if err != nil {
		conversationHistory = []map[string]string{}
	}

	geminiResponse, _, err := h.container.GeminiService.ProcessMessageWithContext(
		req.Message,
		conversationHistory,
		req.Country,
		req.Language,
		session.SearchState.Category,
		session.SearchState.LastProduct,
	)

	if err != nil {
		log.Printf("❌ Gemini processing error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "processing_error",
			Message: fmt.Sprintf("AI processing failed: %v", err),
		})
	}

	if geminiResponse == nil {
		log.Printf("❌ Gemini returned nil response")
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "processing_error",
			Message: "AI returned empty response",
		})
	}

	session.SearchState.Category = geminiResponse.Category

	assistantMessage := &models.Message{
		ID:           uuid.New(),
		SessionID:    session.ID,
		Role:         "assistant",
		Content:      geminiResponse.Output,
		ResponseType: geminiResponse.ResponseType,
		QuickReplies: geminiResponse.QuickReplies,
		CreatedAt:    time.Now(),
	}

	if err := h.container.SessionService.AddMessage(req.SessionID, assistantMessage); err != nil {
		fmt.Printf("⚠️ Failed to store assistant message: %v\n", err)
	}

	response := models.ChatResponse{
		Type:         geminiResponse.ResponseType,
		Output:       geminiResponse.Output,
		QuickReplies: geminiResponse.QuickReplies,
		SessionID:    req.SessionID,
		MessageCount: session.MessageCount + 1,
	}

	if geminiResponse.ResponseType == "search" {
		products, searchErr := h.performSearch(geminiResponse, req.Country, req.Language)
		if searchErr != nil {
			log.Printf("⚠️ Search failed: %v", searchErr)
			response.Output = "Sorry, I couldn't find any products. Please try different keywords."
			response.Type = "text"
		} else if len(products) > 0 {
			response.Products = products
			response.SearchType = geminiResponse.SearchType

			if len(products) > 0 {
				priceStr := products[0].Price
				priceStr = strings.ReplaceAll(priceStr, "$", "")
				priceStr = strings.ReplaceAll(priceStr, "€", "")
				priceStr = strings.ReplaceAll(priceStr, "£", "")
				priceStr = strings.ReplaceAll(priceStr, "CHF", "")
				priceStr = strings.TrimSpace(priceStr)
				priceStr = strings.ReplaceAll(priceStr, ",", "")

				price, _ := strconv.ParseFloat(priceStr, 64)

				session.SearchState.LastProduct = &models.ProductInfo{
					Name:  products[0].Name,
					Price: price,
				}
			}

			session.SearchState.SearchCount++
			assistantMessage.Products = products
		}
	}

	session.SearchState.Status = models.SearchStatusIdle
	if err := h.container.SessionService.UpdateSession(session); err != nil {
		fmt.Printf("⚠️ Failed to update session: %v\n", err)
	}

	response.SearchState = &models.SearchStateResponse{
		Status:      string(session.SearchState.Status),
		Category:    session.SearchState.Category,
		CanContinue: session.SearchState.SearchCount < h.container.SessionService.GetMaxSearches(),
		SearchCount: session.SearchState.SearchCount,
		MaxSearches: h.container.SessionService.GetMaxSearches(),
	}

	return c.JSON(response)
}

func (h *ChatHandler) performSearch(geminiResp *models.GeminiResponse, country, language string) ([]models.ProductCard, error) {
	products, _, err := h.container.SerpService.SearchProducts(
		geminiResp.SearchPhrase,
		geminiResp.SearchType,
		country,
	)

	if err != nil {
		return nil, err
	}

	return products, nil
}
