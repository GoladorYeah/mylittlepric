package middleware

import (
	"fmt"
	"net/http"
	"time"

	"mylittleprice/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// SessionService interface to avoid circular dependency
type SessionService interface {
	GetSession(sessionID string) (interface{}, error)
}

// SessionOwnershipValidator validates that the user has access to the requested session
type SessionOwnershipValidator struct {
	sessionService SessionService
	Signer         *utils.SessionSignature // Public for WebSocket handler access
}

// NewSessionOwnershipValidator creates a new session ownership validator
func NewSessionOwnershipValidator(sessionService SessionService, secretKey string) *SessionOwnershipValidator {
	return &SessionOwnershipValidator{
		sessionService: sessionService,
		Signer:         utils.NewSessionSignature(secretKey),
	}
}

// ValidateSessionOwnership is a middleware that validates session ownership
// It extracts session_id from query/body and verifies the user has access to it
func (v *SessionOwnershipValidator) ValidateSessionOwnership() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Extract session_id from query or body
		sessionID := c.Query("session_id")
		if sessionID == "" {
			// Try to get from body
			var body struct {
				SessionID string `json:"session_id"`
			}
			if err := c.BodyParser(&body); err == nil {
				sessionID = body.SessionID
			}
		}

		// If no session_id provided, skip validation (might be creating new session)
		if sessionID == "" {
			return c.Next()
		}

		// Check if session ID is signed
		if v.Signer.IsSignedSessionID(sessionID) {
			// Verify signed session ID
			rawSessionID, embeddedUserID, err := v.Signer.VerifyAndExtractSessionID(sessionID, 24*time.Hour)
			if err != nil {
				return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
					"error": "Invalid or expired session signature",
				})
			}

			// Store raw session ID for handlers
			c.Locals("session_id", rawSessionID)

			// If signed session has user_id, verify it matches authenticated user
			if embeddedUserID != nil {
				userUUID, ok := GetUserID(c)
				if !ok {
					return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
						"error": "Session requires authentication",
					})
				}

				if *embeddedUserID != userUUID {
					return c.Status(http.StatusForbidden).JSON(fiber.Map{
						"error": "You don't have permission to access this session",
					})
				}
			}

			return c.Next()
		}

		// Unsigned session ID - validate ownership via database
		rawSessionID := utils.ExtractRawSessionID(sessionID)
		c.Locals("session_id", rawSessionID)

		// Get session from database
		sessionData, err := v.sessionService.GetSession(rawSessionID)
		if err != nil {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"error": "Session not found",
			})
		}

		// Extract UserID from session (assuming it has a UserID field)
		// We use type assertion to access the UserID
		if sessionMap, ok := sessionData.(map[string]interface{}); ok {
			if userIDStr, ok := sessionMap["UserID"].(string); ok && userIDStr != "" {
				sessionUserID, err := uuid.Parse(userIDStr)
				if err == nil {
					// Session has a user_id, verify it matches authenticated user
					userUUID, ok := GetUserID(c)
					if !ok {
						return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
							"error": "Session requires authentication",
						})
					}

					if sessionUserID != userUUID {
						return c.Status(http.StatusForbidden).JSON(fiber.Map{
							"error": "You don't have permission to access this session",
						})
					}
				}
			}
		}

		return c.Next()
	}
}

// ValidateSessionOwnershipStrict is a stricter version that requires signed session IDs
func (v *SessionOwnershipValidator) ValidateSessionOwnershipStrict() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Extract session_id from query or body
		sessionID := c.Query("session_id")
		if sessionID == "" {
			var body struct {
				SessionID string `json:"session_id"`
			}
			if err := c.BodyParser(&body); err == nil {
				sessionID = body.SessionID
			}
		}

		if sessionID == "" {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "session_id is required",
			})
		}

		// Require signed session ID
		if !v.Signer.IsSignedSessionID(sessionID) {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "Unsigned session IDs are not allowed. Please use signed session IDs.",
			})
		}

		// Verify signed session ID
		rawSessionID, embeddedUserID, err := v.Signer.VerifyAndExtractSessionID(sessionID, 24*time.Hour)
		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired session signature",
			})
		}

		// Store raw session ID for handlers
		c.Locals("session_id", rawSessionID)

		// If session has user_id, verify it matches
		if embeddedUserID != nil {
			userUUID, ok := GetUserID(c)
			if !ok {
				return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
					"error": "Session requires authentication",
				})
			}

			if *embeddedUserID != userUUID {
				return c.Status(http.StatusForbidden).JSON(fiber.Map{
					"error": "You don't have permission to access this session",
				})
			}
		}

		return c.Next()
	}
}

// SignSessionIDForUser creates a signed session ID for the authenticated user
func (v *SessionOwnershipValidator) SignSessionIDForUser(c *fiber.Ctx, sessionID string) (string, error) {
	userUUID, ok := GetUserID(c)
	if !ok {
		// Anonymous user - sign without user_id
		return v.Signer.SignSessionID(sessionID, nil), nil
	}

	// Authenticated user - sign with user_id
	return v.Signer.SignSessionID(sessionID, &userUUID), nil
}

// ValidateWebSocketSessionOwnership validates session ownership for WebSocket connections
// Used in WebSocket upgrade handler
func (v *SessionOwnershipValidator) ValidateWebSocketSessionOwnership(sessionID string, userID *uuid.UUID) error {
	// If session ID is signed, verify it
	if v.Signer.IsSignedSessionID(sessionID) {
		rawSessionID, embeddedUserID, err := v.Signer.VerifyAndExtractSessionID(sessionID, 24*time.Hour)
		if err != nil {
			return fmt.Errorf("invalid session signature: %w", err)
		}

		// If signed session has user_id, verify it matches
		if embeddedUserID != nil {
			if userID == nil {
				return fmt.Errorf("session requires authentication")
			}
			if *embeddedUserID != *userID {
				return fmt.Errorf("session belongs to different user")
			}
		}

		sessionID = rawSessionID
	}

	// Get session from database
	sessionData, err := v.sessionService.GetSession(sessionID)
	if err != nil {
		return fmt.Errorf("session not found: %w", err)
	}

	// Extract UserID from session
	if sessionMap, ok := sessionData.(map[string]interface{}); ok {
		if userIDStr, ok := sessionMap["UserID"].(string); ok && userIDStr != "" {
			sessionUserID, err := uuid.Parse(userIDStr)
			if err == nil {
				// Session has user_id, verify it matches
				if userID == nil {
					return fmt.Errorf("session requires authentication")
				}
				if sessionUserID != *userID {
					return fmt.Errorf("session belongs to different user")
				}
			}
		}
	}

	return nil
}

// GetSessionIDFromLocals extracts the validated session ID from fiber context locals
func GetSessionIDFromLocals(c *fiber.Ctx) (string, bool) {
	sessionID, ok := c.Locals("session_id").(string)
	return sessionID, ok
}
