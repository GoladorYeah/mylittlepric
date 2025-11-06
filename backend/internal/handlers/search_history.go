package handlers

import (
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"mylittleprice/internal/app"
	"mylittleprice/internal/models"
)

type SearchHistoryHandler struct {
	container *app.Container
}

func NewSearchHistoryHandler(container *app.Container) *SearchHistoryHandler {
	return &SearchHistoryHandler{
		container: container,
	}
}

// GetSearchHistory retrieves search history for the authenticated user or anonymous session
// GET /api/search-history?limit=20&offset=0&session_id=xxx
func (h *SearchHistoryHandler) GetSearchHistory(c *fiber.Ctx) error {
	ctx := c.Context()

	// Get optional user from context (if authenticated)
	var userID *uuid.UUID
	if uid, ok := c.Locals("user_id").(uuid.UUID); ok {
		userID = &uid
	}

	// Parse query parameters
	limit, err := strconv.Atoi(c.Query("limit", "20"))
	if err != nil {
		limit = 20
	}

	offset, err := strconv.Atoi(c.Query("offset", "0"))
	if err != nil {
		offset = 0
	}

	// Get session_id from query params (for anonymous users)
	var sessionID *string
	if userID == nil {
		sessionIDStr := c.Query("session_id")
		if sessionIDStr != "" {
			sessionID = &sessionIDStr
		}
	}

	// Get search history
	history, err := h.container.SearchHistoryService.GetUserSearchHistory(ctx, userID, sessionID, limit, offset)
	if err != nil {
		log.Printf("Error getting search history: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "HISTORY_FETCH_ERROR",
			Message: "Failed to retrieve search history",
		})
	}

	return c.JSON(history)
}

// DeleteSearchHistory deletes a specific search history entry
// DELETE /api/search-history/:id
func (h *SearchHistoryHandler) DeleteSearchHistory(c *fiber.Ctx) error {
	ctx := c.Context()

	// Parse history ID
	historyID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "INVALID_ID",
			Message: "Invalid history ID",
		})
	}

	// Get optional user from context
	var userID *uuid.UUID
	if uid, ok := c.Locals("user_id").(uuid.UUID); ok {
		userID = &uid
	}

	// Delete the history entry
	err = h.container.SearchHistoryService.DeleteSearchHistory(ctx, historyID, userID)
	if err != nil {
		log.Printf("Error deleting search history: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "DELETE_ERROR",
			Message: err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Search history deleted successfully",
	})
}

// DeleteAllSearchHistory deletes all search history for the authenticated user
// DELETE /api/search-history
func (h *SearchHistoryHandler) DeleteAllSearchHistory(c *fiber.Ctx) error {
	ctx := c.Context()

	// Get user from context (authentication required for this endpoint)
	uid, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
			Error:   "UNAUTHORIZED",
			Message: "Authentication required to delete all history",
		})
	}

	// Delete all history
	err := h.container.SearchHistoryService.DeleteAllUserSearchHistory(ctx, uid)
	if err != nil {
		log.Printf("Error deleting all search history: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "DELETE_ERROR",
			Message: "Failed to delete search history",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "All search history deleted successfully",
	})
}

// TrackProductClick updates which product was clicked in a search
// POST /api/search-history/:id/click
func (h *SearchHistoryHandler) TrackProductClick(c *fiber.Ctx) error {
	ctx := c.Context()

	// Parse history ID
	historyID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "INVALID_ID",
			Message: "Invalid history ID",
		})
	}

	// Parse request body
	var req struct {
		ProductID string `json:"product_id" validate:"required"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "INVALID_REQUEST",
			Message: "Invalid request body",
		})
	}

	// Update clicked product
	err = h.container.SearchHistoryService.UpdateClickedProduct(ctx, historyID, req.ProductID)
	if err != nil {
		log.Printf("Error tracking product click: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "UPDATE_ERROR",
			Message: "Failed to track product click",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Product click tracked successfully",
	})
}
