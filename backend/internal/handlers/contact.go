package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"mylittleprice/internal/container"
)

type ContactHandler struct {
	container *container.Container
}

func NewContactHandler(c *container.Container) *ContactHandler {
	return &ContactHandler{container: c}
}

type ContactFormRequest struct {
	Name    string `json:"name" validate:"required"`
	Email   string `json:"email" validate:"required,email"`
	Subject string `json:"subject" validate:"required"`
	Message string `json:"message" validate:"required"`
}

// SubmitContactForm handles contact form submissions
// POST /api/contact
func (h *ContactHandler) SubmitContactForm(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	var req ContactFormRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate required fields
	if req.Name == "" || req.Email == "" || req.Subject == "" || req.Message == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "All fields are required",
		})
	}

	// Basic email validation
	if len(req.Email) < 3 || !contains(req.Email, "@") {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid email address",
		})
	}

	// Store in Redis with 90-day expiration for record keeping
	contactData, err := json.Marshal(req)
	if err != nil {
		slog.Error("Failed to marshal contact form", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to process contact form",
		})
	}

	key := "contact_form:" + time.Now().Format("20060102_150405") + ":" + sanitizeEmail(req.Email)
	if err := h.container.Redis.Set(ctx, key, contactData, 90*24*time.Hour).Err(); err != nil {
		slog.Error("Failed to store contact form in Redis", "error", err)
		// Don't fail - we can still send to Discord
	}

	// Log the contact form submission
	slog.Info("Contact form submitted",
		"name", req.Name,
		"email", req.Email,
		"subject", req.Subject,
		"timestamp", time.Now().Format(time.RFC3339),
	)

	// Send Discord notification if webhook URL is configured
	if h.container.Config.ContactWebhookURL != "" {
		go h.sendDiscordNotification(req, key)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Thank you for your message! We'll get back to you soon.",
	})
}

// sendDiscordNotification sends a contact form notification to Discord webhook
func (h *ContactHandler) sendDiscordNotification(req ContactFormRequest, contactID string) {
	// Truncate message if too long
	message := req.Message
	if len(message) > 1000 {
		message = message[:997] + "..."
	}

	// Build fields for Discord embed
	fields := []map[string]interface{}{
		{
			"name":   "👤 Name",
			"value":  req.Name,
			"inline": true,
		},
		{
			"name":   "📧 Email",
			"value":  req.Email,
			"inline": true,
		},
		{
			"name":   "📋 Subject",
			"value":  req.Subject,
			"inline": false,
		},
		{
			"name":   "💬 Message",
			"value":  message,
			"inline": false,
		},
	}

	// Create Discord embed
	embed := map[string]interface{}{
		"title":       "📬 New Contact Form Submission",
		"description": fmt.Sprintf("**%s** sent a message", req.Name),
		"color":       3447003, // Blue color
		"fields":      fields,
		"footer": map[string]interface{}{
			"text": fmt.Sprintf("Contact ID: %s", contactID),
		},
		"timestamp": time.Now().Format(time.RFC3339),
	}

	payload := map[string]interface{}{
		"embeds": []map[string]interface{}{embed},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		slog.Error("Failed to marshal Discord payload", "error", err)
		return
	}

	resp, err := http.Post(h.container.Config.ContactWebhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		slog.Error("Failed to send Discord notification", "error", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		slog.Error("Discord webhook returned error", "status", resp.StatusCode)
		return
	}

	slog.Info("Discord notification sent successfully", "contact_id", contactID)
}

// Helper functions
func contains(s, substr string) bool {
	for i := 0; i < len(s)-len(substr)+1; i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func sanitizeEmail(email string) string {
	// Replace @ and . with underscores for safe key naming
	result := ""
	for _, ch := range email {
		if ch == '@' || ch == '.' {
			result += "_"
		} else {
			result += string(ch)
		}
	}
	return result
}
