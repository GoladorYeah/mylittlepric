package services

import (
	"fmt"
	"regexp"
	"strings"
)

// Email validation regex (RFC 5322 simplified)
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// validateEmail checks if the email format is valid
func validateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("email cannot be empty")
	}

	email = strings.TrimSpace(email)

	if len(email) > 255 {
		return fmt.Errorf("email too long (max 255 characters)")
	}

	if !emailRegex.MatchString(email) {
		return fmt.Errorf("invalid email format")
	}

	return nil
}

// validatePassword checks if the password meets minimum requirements
func validatePassword(password string) error {
	if password == "" {
		return fmt.Errorf("password cannot be empty")
	}

	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}

	if len(password) > 128 {
		return fmt.Errorf("password too long (max 128 characters)")
	}

	return nil
}

// validateSessionID checks if the session ID format is valid
func validateSessionID(sessionID string) error {
	if sessionID == "" {
		return fmt.Errorf("session ID cannot be empty")
	}

	// Session IDs can be:
	// 1. Plain: UUID or similar (16-128 chars, alphanumeric + hyphens/underscores)
	// 2. Signed: sessionID.timestamp.signature or sessionID.timestamp.userID.signature

	// Maximum length increased to support signed session IDs
	if len(sessionID) < 16 || len(sessionID) > 512 {
		return fmt.Errorf("invalid session ID length")
	}

	// Allow alphanumeric, hyphens, underscores, and dots (for signed session IDs)
	// Also allow = and - for base64 URL encoding in signatures
	if !regexp.MustCompile(`^[a-zA-Z0-9\-_.=]+$`).MatchString(sessionID) {
		return fmt.Errorf("invalid session ID format")
	}

	return nil
}

// validateSearchQuery checks if the search query is valid
func validateSearchQuery(query string) error {
	if query == "" {
		return fmt.Errorf("search query cannot be empty")
	}

	query = strings.TrimSpace(query)

	if len(query) < 2 {
		return fmt.Errorf("search query too short (minimum 2 characters)")
	}

	if len(query) > 500 {
		return fmt.Errorf("search query too long (max 500 characters)")
	}

	return nil
}

// validateIDToken checks if the OAuth ID token format is valid
func validateIDToken(idToken string) error {
	if idToken == "" {
		return fmt.Errorf("ID token cannot be empty")
	}

	// JWT tokens have 3 parts separated by dots
	parts := strings.Split(idToken, ".")
	if len(parts) != 3 {
		return fmt.Errorf("invalid ID token format (expected JWT)")
	}

	// Each part should be non-empty
	for i, part := range parts {
		if part == "" {
			return fmt.Errorf("invalid ID token format (part %d is empty)", i+1)
		}
	}

	// Reasonable length check (JWTs are typically 100-2000 characters)
	if len(idToken) < 100 || len(idToken) > 4096 {
		return fmt.Errorf("invalid ID token length")
	}

	return nil
}
