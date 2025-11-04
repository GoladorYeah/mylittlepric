package services

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"mylittleprice/internal/models"
	"mylittleprice/internal/utils"
)

const (
	PromptIDUniversal = "UniversalPrompt v1.0.1"
	MaxIterations     = 6
)

// UniversalPromptManager manages the Universal Prompt system with mini-kernel approach
type UniversalPromptManager struct {
	universalPrompt string
	miniKernel      string
	promptHasher    *utils.PromptHasher
	promptHash      string
	mu              sync.RWMutex
}

// NewUniversalPromptManager creates a new Universal Prompt Manager
func NewUniversalPromptManager() *UniversalPromptManager {
	upm := &UniversalPromptManager{
		promptHasher: utils.NewPromptHasher(),
	}
	upm.loadPrompts()
	return upm
}

// loadPrompts loads the universal prompt and mini-kernel from files
func (upm *UniversalPromptManager) loadPrompts() {
	upm.mu.Lock()
	defer upm.mu.Unlock()

	// Load universal prompt (sent once on session start)
	universalPath := "internal/services/prompts/universal_prompt.txt"
	universalContent, err := os.ReadFile(universalPath)
	if err != nil {
		panic(fmt.Errorf("CRITICAL: Failed to load universal prompt: %w", err))
	}
	upm.universalPrompt = string(universalContent)

	// Load mini-kernel (sent on every turn)
	kernelPath := "internal/services/prompts/mini_kernel.txt"
	kernelContent, err := os.ReadFile(kernelPath)
	if err != nil {
		panic(fmt.Errorf("CRITICAL: Failed to load mini-kernel: %w", err))
	}
	upm.miniKernel = string(kernelContent)

	// Generate hash for drift detection
	upm.promptHash = upm.promptHasher.HashPrompt(upm.universalPrompt)

	fmt.Printf("✅ Universal Prompt System loaded (hash: %s)\n", upm.promptHasher.HashPromptShort(upm.universalPrompt))
}

// GetSystemPrompt returns the full system prompt for NEW sessions
// This is sent ONCE when the session starts
func (upm *UniversalPromptManager) GetSystemPrompt(
	feLocation, feLanguage, feCurrency string,
) string {
	upm.mu.RLock()
	defer upm.mu.RUnlock()

	prompt := upm.universalPrompt
	prompt = strings.ReplaceAll(prompt, "{fe_location}", feLocation)
	prompt = strings.ReplaceAll(prompt, "{fe_language}", feLanguage)
	prompt = strings.ReplaceAll(prompt, "{fe_currency}", feCurrency)

	return prompt
}

// GetMiniKernel returns the mini-kernel for EVERY turn
// This ensures the rules are always in context
func (upm *UniversalPromptManager) GetMiniKernel(
	feLocation, feLanguage, feCurrency string,
	cycleState *models.CycleState,
) string {
	upm.mu.RLock()
	defer upm.mu.RUnlock()

	kernel := upm.miniKernel
	kernel = strings.ReplaceAll(kernel, "{fe_location}", feLocation)
	kernel = strings.ReplaceAll(kernel, "{fe_language}", feLanguage)
	kernel = strings.ReplaceAll(kernel, "{fe_currency}", feCurrency)
	kernel = strings.ReplaceAll(kernel, "{cycle_id}", fmt.Sprintf("%d", cycleState.CycleID))
	kernel = strings.ReplaceAll(kernel, "{iteration}", fmt.Sprintf("%d", cycleState.Iteration))
	kernel = strings.ReplaceAll(kernel, "{category}", getCategory(cycleState))

	return kernel
}

// BuildStateContext builds the state context object sent with each turn
// This includes cycle_history, last_cycle_context, and last_defined
func (upm *UniversalPromptManager) BuildStateContext(
	session *models.ChatSession,
) string {
	cycleState := &session.CycleState

	var sb strings.Builder

	sb.WriteString("=== CURRENT STATE ===\n")
	sb.WriteString(fmt.Sprintf("CYCLE_ID: %d\n", cycleState.CycleID))
	sb.WriteString(fmt.Sprintf("ITERATION: %d/%d\n", cycleState.Iteration, MaxIterations))
	sb.WriteString(fmt.Sprintf("CURRENT_CATEGORY: %s\n", getCategory(cycleState)))
	sb.WriteString("\n")

	// Cycle history (limited to last 6 messages to match MaxIterations)
	sb.WriteString("=== CYCLE_HISTORY (Current Cycle) ===\n")
	if len(cycleState.CycleHistory) == 0 {
		sb.WriteString("(empty - first message in cycle)\n")
	} else {
		// Show only last 6 messages to match max iterations per cycle
		// This ensures full visibility of current cycle while managing token usage
		maxRecentMessages := MaxIterations
		startIdx := len(cycleState.CycleHistory) - maxRecentMessages
		if startIdx < 0 {
			startIdx = 0
		}

		// If we're skipping messages, add a summary note
		if startIdx > 0 {
			sb.WriteString(fmt.Sprintf("(showing last %d of %d messages)\n",
				len(cycleState.CycleHistory)-startIdx, len(cycleState.CycleHistory)))
		}

		for i := startIdx; i < len(cycleState.CycleHistory); i++ {
			msg := cycleState.CycleHistory[i]
			sb.WriteString(fmt.Sprintf("%d. %s: %s\n", i+1, msg.Role, msg.Content))
		}
	}
	sb.WriteString("\n")

	// Last cycle context
	if cycleState.LastCycleContext != nil {
		sb.WriteString("=== LAST_CYCLE_CONTEXT ===\n")
		if len(cycleState.LastCycleContext.Groups) > 0 {
			sb.WriteString(fmt.Sprintf("Groups: %s\n", strings.Join(cycleState.LastCycleContext.Groups, ", ")))
		}
		if len(cycleState.LastCycleContext.Subgroups) > 0 {
			sb.WriteString(fmt.Sprintf("Subgroups: %s\n", strings.Join(cycleState.LastCycleContext.Subgroups, ", ")))
		}
		if len(cycleState.LastCycleContext.Products) > 0 {
			sb.WriteString("Products from last cycle:\n")
			for _, p := range cycleState.LastCycleContext.Products {
				sb.WriteString(fmt.Sprintf("  - %s (%.2f)\n", p.Name, p.Price))
			}
		}
		if cycleState.LastCycleContext.LastRequest != "" {
			sb.WriteString(fmt.Sprintf("Last request: %s\n", cycleState.LastCycleContext.LastRequest))
		}
		sb.WriteString("\n")
	}

	// Last defined products
	if len(cycleState.LastDefined) > 0 {
		sb.WriteString("=== LAST_DEFINED (confirmed products) ===\n")
		sb.WriteString(strings.Join(cycleState.LastDefined, ", ") + "\n")
		sb.WriteString("\n")
	}

	return sb.String()
}

// InitializeCycleState creates a new cycle state for a session
func (upm *UniversalPromptManager) InitializeCycleState() models.CycleState {
	return models.CycleState{
		CycleID:          1,
		Iteration:        1,
		CycleHistory:     []models.CycleMessage{},
		LastCycleContext: nil,
		LastDefined:      []string{},
		PromptID:         PromptIDUniversal,
		PromptHash:       upm.GetPromptHash(),
	}
}

// IncrementIteration increments the iteration counter
// Returns true if we should continue in the same cycle, false if we need a new cycle
func (upm *UniversalPromptManager) IncrementIteration(cycleState *models.CycleState) bool {
	// Check BEFORE incrementing to properly handle iteration 6
	if cycleState.Iteration >= MaxIterations {
		fmt.Printf("⚠️ Max iterations reached (%d), need new cycle\n", MaxIterations)
		return false // Need new cycle
	}

	cycleState.Iteration++
	fmt.Printf("📊 Cycle %d, Iteration %d/%d\n", cycleState.CycleID, cycleState.Iteration, MaxIterations)

	return true
}

// StartNewCycle starts a new cycle, carrying over context
func (upm *UniversalPromptManager) StartNewCycle(
	cycleState *models.CycleState,
	lastRequest string,
	products []models.ProductInfo,
) {
	// Save current cycle context
	lastContext := &models.LastCycleContext{
		Groups:      extractGroups(cycleState.CycleHistory),
		Subgroups:   extractSubgroups(cycleState.CycleHistory),
		Products:    products,
		LastRequest: lastRequest,
	}

	// Increment cycle ID
	cycleState.CycleID++
	cycleState.Iteration = 1
	cycleState.CycleHistory = []models.CycleMessage{}
	cycleState.LastCycleContext = lastContext

	fmt.Printf("🔄 Starting new Cycle %d (carried over context from previous cycle)\n", cycleState.CycleID)
}

// AddToCycleHistory adds a message to the current cycle history
func (upm *UniversalPromptManager) AddToCycleHistory(
	cycleState *models.CycleState,
	role, content string,
) {
	msg := models.CycleMessage{
		Role:      role,
		Content:   content,
		Timestamp: time.Now(),
	}
	cycleState.CycleHistory = append(cycleState.CycleHistory, msg)
}

// GetPromptHash returns the SHA-256 hash of the universal prompt
func (upm *UniversalPromptManager) GetPromptHash() string {
	upm.mu.RLock()
	defer upm.mu.RUnlock()
	return upm.promptHash
}

// GetPromptHashShort returns a short version of the hash for logging
func (upm *UniversalPromptManager) GetPromptHashShort() string {
	return upm.promptHasher.HashPromptShort(upm.universalPrompt)
}

// GetPromptID returns the prompt version identifier
func (upm *UniversalPromptManager) GetPromptID() string {
	return PromptIDUniversal
}

// Helper functions

func getCategory(cycleState *models.CycleState) string {
	// Try to extract category from cycle history
	// Look for the latest category mentioned in assistant messages
	for i := len(cycleState.CycleHistory) - 1; i >= 0; i-- {
		msg := cycleState.CycleHistory[i]
		if msg.Role == "assistant" {
			// Try to parse JSON to extract category
			var resp map[string]interface{}
			if err := json.Unmarshal([]byte(msg.Content), &resp); err == nil {
				if cat, ok := resp["category"].(string); ok && cat != "" {
					return cat
				}
			}
		}
	}
	return "unknown"
}

func extractGroups(history []models.CycleMessage) []string {
	groups := make(map[string]bool)
	// Simple extraction - in practice, you'd parse assistant responses for group mentions
	// For now, return empty as this is implementation-dependent
	result := make([]string, 0, len(groups))
	for g := range groups {
		result = append(result, g)
	}
	return result
}

func extractSubgroups(history []models.CycleMessage) []string {
	subgroups := make(map[string]bool)
	// Simple extraction - in practice, you'd parse assistant responses for subgroup mentions
	result := make([]string, 0, len(subgroups))
	for s := range subgroups {
		result = append(result, s)
	}
	return result
}
