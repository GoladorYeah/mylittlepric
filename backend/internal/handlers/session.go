package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"mylittleprice/internal/container"
	"mylittleprice/internal/middleware"
	"mylittleprice/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type SessionHandler struct {
	container *container.Container
}

func NewSessionHandler(c *container.Container) *SessionHandler {
	return &SessionHandler{container: c}
}

// GetActiveSession returns the most recent active session for authenticated user
// GET /api/sessions/active
func (h *SessionHandler) GetActiveSession(c *fiber.Ctx) error {
	// Get user ID from JWT token
	userUUID, ok := middleware.GetUserID(c)
	if !ok {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Get active session from database
	session, err := h.container.SessionService.GetActiveSessionForUser(userUUID)
	if err != nil {
		fmt.Printf("❌ Error getting active session for user %s: %v\n", userUUID.String(), err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get active session",
		})
	}

	// No active session found - return null (this is OK, user can start new session)
	if session == nil {
		return c.JSON(fiber.Map{
			"session": nil,
			"has_active_session": false,
		})
	}

	// Get message count from Redis
	messages, err := h.container.SessionService.GetMessages(session.SessionID)
	if err != nil {
		fmt.Printf("⚠️ Failed to get message count for session %s: %v\n", session.SessionID, err)
		messages = []*models.Message{} // Empty array
	}

	// Return session info
	return c.JSON(fiber.Map{
		"session": fiber.Map{
			"session_id":    session.SessionID,
			"message_count": len(messages),
			"search_state": fiber.Map{
				"status":   session.SearchState.Status,
				"category": session.SearchState.Category,
			},
			"created_at": session.CreatedAt,
			"updated_at": session.UpdatedAt,
			"expires_at": session.ExpiresAt,
		},
		"has_active_session": true,
	})
}

// LinkSessionToUser links current anonymous session to authenticated user
// POST /api/sessions/link
func (h *SessionHandler) LinkSessionToUser(c *fiber.Ctx) error {
	// Get user ID from JWT token
	userUUID, ok := middleware.GetUserID(c)
	if !ok {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Parse request body
	var req struct {
		SessionID string `json:"session_id"`
	}

	if err := json.Unmarshal(c.Body(), &req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.SessionID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "session_id is required",
		})
	}

	// Link session to user
	if err := h.container.SessionService.LinkSessionToUser(req.SessionID, userUUID); err != nil {
		fmt.Printf("❌ Error linking session %s to user %s: %v\n", req.SessionID, userUUID.String(), err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to link session to user",
		})
	}

	fmt.Printf("✅ Linked session %s to user %s\n", req.SessionID, userUUID.String())

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Session linked to user successfully",
	})
}
