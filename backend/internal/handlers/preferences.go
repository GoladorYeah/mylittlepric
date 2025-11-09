package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"mylittleprice/internal/container"
	"mylittleprice/internal/middleware"
	"mylittleprice/internal/models"

	"github.com/gofiber/fiber/v2"
)

type PreferencesHandler struct {
	container *container.Container
}

func NewPreferencesHandler(c *container.Container) *PreferencesHandler {
	return &PreferencesHandler{container: c}
}

// GetUserPreferences returns user preferences for authenticated user
// GET /api/user/preferences
func (h *PreferencesHandler) GetUserPreferences(c *fiber.Ctx) error {
	// Get user ID from JWT token
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Get preferences from database
	prefs, err := h.container.PreferencesService.GetUserPreferences(userID)
	if err != nil {
		fmt.Printf("❌ Error getting preferences for user %s: %v\n", userID.String(), err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get user preferences",
		})
	}

	// If no preferences found, return empty object (first time user)
	if prefs == nil {
		return c.JSON(fiber.Map{
			"preferences": nil,
			"has_preferences": false,
		})
	}

	// Return preferences
	return c.JSON(fiber.Map{
		"preferences": prefs,
		"has_preferences": true,
	})
}

// UpdateUserPreferences creates or updates user preferences
// PUT /api/user/preferences
func (h *PreferencesHandler) UpdateUserPreferences(c *fiber.Ctx) error {
	// Get user ID from JWT token
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Parse request body
	var update models.UserPreferencesUpdate
	if err := json.Unmarshal(c.Body(), &update); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate fields (optional but recommended)
	if update.Country != nil && len(*update.Country) > 2 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid country code (must be 2 characters)",
		})
	}

	if update.Currency != nil && len(*update.Currency) > 3 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid currency code (must be 3 characters)",
		})
	}

	if update.Theme != nil {
		validThemes := map[string]bool{"light": true, "dark": true, "system": true}
		if !validThemes[*update.Theme] {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid theme (must be 'light', 'dark', or 'system')",
			})
		}
	}

	// Upsert preferences (create or update)
	prefs, err := h.container.PreferencesService.UpsertUserPreferences(userID, &update)
	if err != nil {
		fmt.Printf("❌ Error upserting preferences for user %s: %v\n", userID.String(), err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update user preferences",
		})
	}

	fmt.Printf("✅ Upserted preferences for user %s\n", userID.String())

	return c.JSON(fiber.Map{
		"success": true,
		"preferences": prefs,
	})
}
