package handlers

import (
	"github.com/gofiber/fiber/v2"

	"mylittleprice/internal/container"
	"mylittleprice/internal/models"
)

type ChatHandler struct {
	container *container.Container
	processor *ChatProcessor
}

func NewChatHandler(c *container.Container) *ChatHandler {
	return &ChatHandler{
		container: c,
		processor: NewChatProcessor(c),
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

	// Process chat using shared processor
	processorReq := &ChatRequest{
		SessionID:       req.SessionID,
		Message:         req.Message,
		Country:         req.Country,
		Language:        req.Language,
		NewSearch:       req.NewSearch,
		CurrentCategory: "",
	}

	result := h.processor.ProcessChat(processorReq)

	// Handle errors
	if result.Error != nil {
		statusCode := fiber.StatusInternalServerError
		if result.Error.Code == "validation_error" {
			statusCode = fiber.StatusBadRequest
		}
		return c.Status(statusCode).JSON(models.ErrorResponse{
			Error:   result.Error.Code,
			Message: result.Error.Message,
		})
	}

	// Build response
	response := models.ChatResponse{
		Type:         result.Type,
		Output:       result.Output,
		QuickReplies: result.QuickReplies,
		Products:     result.Products,
		SearchType:   result.SearchType,
		SessionID:    result.SessionID,
		MessageCount: result.MessageCount,
		SearchState:  result.SearchState,
	}

	return c.JSON(response)
}
