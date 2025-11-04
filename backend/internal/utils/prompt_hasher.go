package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

// PromptHasher handles SHA-256 hashing for prompt versioning and drift detection
type PromptHasher struct{}

// NewPromptHasher creates a new PromptHasher instance
func NewPromptHasher() *PromptHasher {
	return &PromptHasher{}
}

// HashPrompt generates a SHA-256 hash of the given prompt text
func (h *PromptHasher) HashPrompt(promptText string) string {
	hasher := sha256.New()
	hasher.Write([]byte(promptText))
	return hex.EncodeToString(hasher.Sum(nil))
}

// HashPromptShort returns a short 12-character version of the hash for logging
func (h *PromptHasher) HashPromptShort(promptText string) string {
	fullHash := h.HashPrompt(promptText)
	if len(fullHash) > 12 {
		return fullHash[:12]
	}
	return fullHash
}
