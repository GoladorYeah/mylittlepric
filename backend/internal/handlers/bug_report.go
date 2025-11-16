package handlers

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"mylittleprice/internal/container"
)

type BugReportHandler struct {
	container *container.Container
}

func NewBugReportHandler(c *container.Container) *BugReportHandler {
	return &BugReportHandler{container: c}
}

type BugReportContext struct {
	UserID         string                   `json:"user_id"`
	UserEmail      string                   `json:"user_email"`
	SessionID      string                   `json:"session_id"`
	URL            string                   `json:"url"`
	UserAgent      string                   `json:"user_agent"`
	ScreenRes      string                   `json:"screen_resolution"`
	ViewportSize   string                   `json:"viewport_size"`
	Timestamp      string                   `json:"timestamp"`
	LastMessages   []map[string]interface{} `json:"last_messages"`
}

type BugReportRequest struct {
	Description      string           `json:"description"`
	StepsToReproduce string           `json:"steps_to_reproduce"`
	Context          BugReportContext `json:"context"`
}

// SubmitBugReport handles bug report submissions
// POST /api/bug-report
func (h *BugReportHandler) SubmitBugReport(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 30*time.Second) // Increased timeout for file uploads
	defer cancel()

	// Parse form data
	description := c.FormValue("description")
	stepsToReproduce := c.FormValue("steps_to_reproduce")
	contextJSON := c.FormValue("context")

	// Validate required fields
	if description == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Description is required",
		})
	}

	// Parse context
	var bugContext BugReportContext
	if err := json.Unmarshal([]byte(contextJSON), &bugContext); err != nil {
		slog.Error("Failed to parse context", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid context data",
		})
	}

	// Handle file attachments
	form, err := c.MultipartForm()
	var attachments []string // Base64 encoded images for Discord
	var attachmentNames []string

	if err == nil && form != nil {
		files := form.File["attachments"]
		for i, fileHeader := range files {
			if i >= 3 { // Max 3 files
				break
			}

			// Check file size (max 5MB)
			if fileHeader.Size > 5*1024*1024 {
				continue
			}

			// Open file
			file, err := fileHeader.Open()
			if err != nil {
				slog.Error("Failed to open uploaded file", "error", err)
				continue
			}
			defer file.Close()

			// Read file content
			fileData, err := io.ReadAll(file)
			if err != nil {
				slog.Error("Failed to read file", "error", err)
				continue
			}

			// Encode to base64 for Discord
			encoded := base64.StdEncoding.EncodeToString(fileData)
			attachments = append(attachments, encoded)
			attachmentNames = append(attachmentNames, fileHeader.Filename)
		}
	}

	// Create request object
	req := BugReportRequest{
		Description:      description,
		StepsToReproduce: stepsToReproduce,
		Context:          bugContext,
	}

	// Store in database
	bugReportData, err := json.Marshal(req)
	if err != nil {
		slog.Error("Failed to marshal bug report", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to process bug report",
		})
	}

	// Store in Redis with 30-day expiration for quick access
	key := "bug_report:" + time.Now().Format("20060102_150405") + ":" + req.Context.UserID
	if err := h.container.Redis.Set(ctx, key, bugReportData, 30*24*time.Hour).Err(); err != nil {
		slog.Error("Failed to store bug report in Redis", "error", err)
		// Don't fail - we can still log it
	}

	// Log the bug report for external monitoring systems to pick up
	slog.Info("Bug report submitted",
		"user_id", req.Context.UserID,
		"user_email", req.Context.UserEmail,
		"session_id", req.Context.SessionID,
		"url", req.Context.URL,
		"description", req.Description,
		"steps", req.StepsToReproduce,
		"user_agent", req.Context.UserAgent,
		"screen_res", req.Context.ScreenRes,
		"viewport", req.Context.ViewportSize,
		"timestamp", req.Context.Timestamp,
		"attachments_count", len(attachments),
	)

	// Send Discord notification if webhook URL is configured
	if h.container.Config.DiscordWebhookURL != "" {
		go h.sendDiscordNotification(req, key, attachments, attachmentNames)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Bug report submitted successfully",
		"id":      key,
	})
}

// GetBugReports retrieves recent bug reports (admin only)
// GET /api/bug-reports
func (h *BugReportHandler) GetBugReports(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	// Get all bug report keys
	keys, err := h.container.Redis.Keys(ctx, "bug_report:*").Result()
	if err != nil {
		slog.Error("Failed to fetch bug report keys", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch bug reports",
		})
	}

	reports := make([]map[string]interface{}, 0, len(keys))
	for _, key := range keys {
		data, err := h.container.Redis.Get(ctx, key).Result()
		if err != nil {
			continue
		}

		var report map[string]interface{}
		if err := json.Unmarshal([]byte(data), &report); err != nil {
			continue
		}

		report["id"] = key
		reports = append(reports, report)
	}

	return c.JSON(fiber.Map{
		"reports": reports,
		"total":   len(reports),
	})
}

// sendDiscordNotification sends a bug report notification to Discord webhook
func (h *BugReportHandler) sendDiscordNotification(req BugReportRequest, reportID string, attachments []string, attachmentNames []string) {
	// Truncate description if too long
	description := req.Description
	if len(description) > 1000 {
		description = description[:997] + "..."
	}

	// Truncate steps if too long
	steps := req.StepsToReproduce
	if len(steps) > 500 {
		steps = steps[:497] + "..."
	}

	// Build fields for Discord embed
	fields := []map[string]interface{}{
		{
			"name":   "👤 User",
			"value":  fmt.Sprintf("%s\n`%s`", req.Context.UserEmail, req.Context.UserID),
			"inline": true,
		},
		{
			"name":   "📍 Location",
			"value":  fmt.Sprintf("[%s](%s)", req.Context.URL, req.Context.URL),
			"inline": true,
		},
		{
			"name":   "🔧 Session ID",
			"value":  fmt.Sprintf("`%s`", req.Context.SessionID),
			"inline": false,
		},
	}

	// Add steps if provided
	if steps != "" {
		fields = append(fields, map[string]interface{}{
			"name":   "📝 Steps to Reproduce",
			"value":  steps,
			"inline": false,
		})
	}

	// Add technical details
	fields = append(fields, map[string]interface{}{
		"name": "💻 Technical Details",
		"value": fmt.Sprintf(
			"**Browser:** %s\n**Screen:** %s\n**Viewport:** %s",
			req.Context.UserAgent,
			req.Context.ScreenRes,
			req.Context.ViewportSize,
		),
		"inline": false,
	})

	// Create Discord embed
	embed := map[string]interface{}{
		"title":       "🐛 New Bug Report",
		"description": description,
		"color":       15158332, // Red color
		"fields":      fields,
		"footer": map[string]interface{}{
			"text": fmt.Sprintf("Report ID: %s", reportID),
		},
		"timestamp": req.Context.Timestamp,
	}

	// Set first image as embed image if available
	if len(attachments) > 0 {
		embed["image"] = map[string]interface{}{
			"url": fmt.Sprintf("attachment://%s", attachmentNames[0]),
		}
	}

	// Create webhook payload with multipart form data for attachments
	if len(attachments) > 0 {
		h.sendDiscordWithAttachments(embed, reportID, attachments, attachmentNames)
	} else {
		h.sendDiscordWithoutAttachments(embed, reportID)
	}
}

// sendDiscordWithoutAttachments sends Discord notification without files
func (h *BugReportHandler) sendDiscordWithoutAttachments(embed map[string]interface{}, reportID string) {
	payload := map[string]interface{}{
		"embeds": []map[string]interface{}{embed},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		slog.Error("Failed to marshal Discord payload", "error", err)
		return
	}

	resp, err := http.Post(h.container.Config.DiscordWebhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		slog.Error("Failed to send Discord notification", "error", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		slog.Error("Discord webhook returned error", "status", resp.StatusCode)
		return
	}

	slog.Info("Discord notification sent successfully", "report_id", reportID)
}

// sendDiscordWithAttachments sends Discord notification with file attachments
func (h *BugReportHandler) sendDiscordWithAttachments(embed map[string]interface{}, reportID string, attachments []string, attachmentNames []string) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	// Add JSON payload
	payload := map[string]interface{}{
		"embeds": []map[string]interface{}{embed},
	}
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		slog.Error("Failed to marshal Discord payload", "error", err)
		return
	}

	if err := writer.WriteField("payload_json", string(payloadJSON)); err != nil {
		slog.Error("Failed to write payload field", "error", err)
		return
	}

	// Add attachments
	for i, encodedData := range attachments {
		if i >= len(attachmentNames) {
			break
		}

		fileData, err := base64.StdEncoding.DecodeString(encodedData)
		if err != nil {
			slog.Error("Failed to decode attachment", "error", err)
			continue
		}

		part, err := writer.CreateFormFile(fmt.Sprintf("files[%d]", i), attachmentNames[i])
		if err != nil {
			slog.Error("Failed to create form file", "error", err)
			continue
		}

		if _, err := part.Write(fileData); err != nil {
			slog.Error("Failed to write file data", "error", err)
			continue
		}
	}

	if err := writer.Close(); err != nil {
		slog.Error("Failed to close multipart writer", "error", err)
		return
	}

	// Send request
	req, err := http.NewRequest("POST", h.container.Config.DiscordWebhookURL, &body)
	if err != nil {
		slog.Error("Failed to create HTTP request", "error", err)
		return
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("Failed to send Discord notification", "error", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		slog.Error("Discord webhook returned error", "status", resp.StatusCode)
		return
	}

	slog.Info("Discord notification with attachments sent successfully", "report_id", reportID, "attachments", len(attachments))
}
