package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// SessionSignature provides HMAC-based session ID signing and verification
// Format: sessionID.timestamp.signature
type SessionSignature struct {
	secretKey []byte
}

// NewSessionSignature creates a new session signature handler
func NewSessionSignature(secretKey string) *SessionSignature {
	return &SessionSignature{
		secretKey: []byte(secretKey),
	}
}

// SignSessionID creates a signed session ID with timestamp
// Format: sessionID.timestamp.signature
// Example: abc123.1699999999.base64signature
func (s *SessionSignature) SignSessionID(sessionID string, userID *uuid.UUID) string {
	timestamp := time.Now().Unix()

	// Create payload: sessionID.timestamp.userID
	var payload string
	if userID != nil {
		payload = fmt.Sprintf("%s.%d.%s", sessionID, timestamp, userID.String())
	} else {
		payload = fmt.Sprintf("%s.%d", sessionID, timestamp)
	}

	// Generate HMAC signature
	h := hmac.New(sha256.New, s.secretKey)
	h.Write([]byte(payload))
	signature := base64.URLEncoding.EncodeToString(h.Sum(nil))

	// Return signed session ID
	if userID != nil {
		return fmt.Sprintf("%s.%d.%s.%s", sessionID, timestamp, userID.String(), signature)
	}
	return fmt.Sprintf("%s.%d.%s", sessionID, timestamp, signature)
}

// VerifyAndExtractSessionID verifies the signed session ID and extracts the original session ID
// Returns: sessionID, userID (or nil), error
func (s *SessionSignature) VerifyAndExtractSessionID(signedSessionID string, maxAge time.Duration) (string, *uuid.UUID, error) {
	parts := strings.Split(signedSessionID, ".")

	// Check format: sessionID.timestamp.signature OR sessionID.timestamp.userID.signature
	if len(parts) != 3 && len(parts) != 4 {
		return "", nil, fmt.Errorf("invalid signed session ID format")
	}

	sessionID := parts[0]
	timestamp := parts[1]

	var userID *uuid.UUID
	var signature string

	if len(parts) == 4 {
		// Format with userID: sessionID.timestamp.userID.signature
		userIDStr := parts[2]
		signature = parts[3]

		parsed, err := uuid.Parse(userIDStr)
		if err != nil {
			return "", nil, fmt.Errorf("invalid user ID in signed session")
		}
		userID = &parsed
	} else {
		// Format without userID: sessionID.timestamp.signature
		signature = parts[2]
	}

	// Parse timestamp
	var ts int64
	_, err := fmt.Sscanf(timestamp, "%d", &ts)
	if err != nil {
		return "", nil, fmt.Errorf("invalid timestamp in signed session")
	}

	// Check if signature is expired
	if maxAge > 0 {
		age := time.Since(time.Unix(ts, 0))
		if age > maxAge {
			return "", nil, fmt.Errorf("signed session ID expired (age: %v, max: %v)", age, maxAge)
		}
	}

	// Verify signature
	var payload string
	if userID != nil {
		payload = fmt.Sprintf("%s.%s.%s", sessionID, timestamp, userID.String())
	} else {
		payload = fmt.Sprintf("%s.%s", sessionID, timestamp)
	}

	h := hmac.New(sha256.New, s.secretKey)
	h.Write([]byte(payload))
	expectedSignature := base64.URLEncoding.EncodeToString(h.Sum(nil))

	if !hmac.Equal([]byte(signature), []byte(expectedSignature)) {
		return "", nil, fmt.Errorf("invalid signature")
	}

	return sessionID, userID, nil
}

// IsSignedSessionID checks if a session ID is signed (without full verification)
func (s *SessionSignature) IsSignedSessionID(sessionID string) bool {
	parts := strings.Split(sessionID, ".")
	return len(parts) == 3 || len(parts) == 4
}

// ExtractRawSessionID extracts the raw session ID from a signed session ID (without verification)
// Use this only when you don't need to verify the signature
func ExtractRawSessionID(signedSessionID string) string {
	parts := strings.Split(signedSessionID, ".")
	if len(parts) >= 3 {
		return parts[0]
	}
	return signedSessionID
}
