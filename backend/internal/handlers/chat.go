package handlers

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

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

	// Get optional user ID from context
	var userID *uuid.UUID
	if uid, ok := c.Locals("user_id").(uuid.UUID); ok {
		userID = &uid
	}

	// Process chat using shared processor
	processorReq := &ChatRequest{
		SessionID:       req.SessionID,
		UserID:          userID,
		Message:         req.Message,
		Country:         req.Country,
		Language:        req.Language,
		Currency:        req.Currency,
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

func (h *ChatHandler) GetSessionMessages(c *fiber.Ctx) error {
	sessionID := c.Query("session_id")
	if sessionID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "invalid_request",
			Message: "session_id is required",
		})
	}

	// Get messages from session
	messages, err := h.container.MessageService.GetMessages(sessionID)
	if err != nil {
		// Check if it's a "not found" error (new session with no messages yet)
		// This is not a critical error - just return empty messages
		fmt.Printf("ℹ️ No messages found for session %s (new session?): %v\n", sessionID, err)

		// Try to get session to verify it exists or will be created
		session, sessionErr := h.container.SessionService.GetSession(sessionID)
		if sessionErr != nil {
			// Session doesn't exist yet - this is OK for new sessions
			return c.JSON(fiber.Map{
				"messages":      []interface{}{},
				"session_id":    sessionID,
				"message_count": 0,
			})
		}

		// Session exists but has no messages yet
		return c.JSON(fiber.Map{
			"messages":      []interface{}{},
			"session_id":    sessionID,
			"message_count": session.MessageCount,
		})
	}

	// Get session info for additional context
	session, err := h.container.SessionService.GetSession(sessionID)
	if err != nil {
		// Session might not exist yet, return messages without session context
		return c.JSON(fiber.Map{
			"messages":      messages,
			"session_id":    sessionID,
			"message_count": len(messages),
		})
	}

	return c.JSON(fiber.Map{
		"messages":      messages,
		"session_id":    sessionID,
		"message_count": len(messages),
		"search_state":  session.SearchState,
	})
}

// GetMessagesSince retrieves messages created after a specific timestamp
// This is useful for reconnection scenarios
// GET /api/chat/messages/since?session_id=xxx&since=2024-01-01T00:00:00Z
func (h *ChatHandler) GetMessagesSince(c *fiber.Ctx) error {
	sessionID := c.Query("session_id")
	if sessionID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "validation_error",
			Message: "session_id is required",
		})
	}

	sinceStr := c.Query("since")
	if sinceStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "validation_error",
			Message: "since timestamp is required",
		})
	}

	// Parse timestamp
	since, err := time.Parse(time.RFC3339, sinceStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid since timestamp format. Use RFC3339 format (e.g., 2024-01-01T00:00:00Z)",
		})
	}

	// Get messages since timestamp
	messages, err := h.container.MessageService.GetMessagesSince(sessionID, since)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "server_error",
			Message: "Failed to retrieve messages",
		})
	}

	return c.JSON(fiber.Map{
		"messages":      messages,
		"session_id":    sessionID,
		"message_count": len(messages),
		"since":         since,
	})
}
