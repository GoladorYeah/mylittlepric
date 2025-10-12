package services

import (
	"fmt"
	"strings"
	"time"
)

// GroundingStrategy determines when to use grounding (Google Search)
type GroundingStrategy struct {
	config *GroundingConfig
}

// GroundingConfig contains grounding configuration
type GroundingConfig struct {
	Mode                string
	MinWordsForProduct  int
	EnableTechnicalSpec bool
}

// GroundingDecision contains the decision about grounding usage
type GroundingDecision struct {
	ShouldUseGrounding bool
	Reason             string
	Confidence         float32
}

// GroundingStats tracks grounding usage statistics
type GroundingStats struct {
	TotalDecisions    int
	GroundingEnabled  int
	GroundingDisabled int
	ReasonCounts      map[string]int
	AverageConfidence float32
}

// NewGroundingStrategy creates a new grounding strategy
func NewGroundingStrategy(mode string) *GroundingStrategy {
	config := &GroundingConfig{
		Mode:                mode,
		MinWordsForProduct:  2,
		EnableTechnicalSpec: true,
	}

	switch mode {
	case "conservative":
		config.MinWordsForProduct = 3
		config.EnableTechnicalSpec = false
	case "aggressive":
		config.MinWordsForProduct = 1
		config.EnableTechnicalSpec = true
	default: // balanced
		config.MinWordsForProduct = 2
		config.EnableTechnicalSpec = true
	}

	return &GroundingStrategy{
		config: config,
	}
}

// Decide makes decision about whether to use grounding
func (gs *GroundingStrategy) Decide(userMessage string, history []map[string]string) GroundingDecision {
	messageLower := strings.ToLower(userMessage)
	wordCount := len(strings.Fields(userMessage))

	// 🆕 PRIORITY 0: Brand selection for latest models (НОВОЕ!)
	if gs.isBrandSelectionForLatestModels(messageLower, history) {
		return GroundingDecision{
			ShouldUseGrounding: true,
			Reason:             "brand_selection_latest_models",
			Confidence:         0.98,
		}
	}

	// PRIORITY 1: Specific product model (HIGHEST)
	if gs.hasSpecificModelPattern(messageLower, wordCount) {
		return GroundingDecision{
			ShouldUseGrounding: true,
			Reason:             "specific_product_verification",
			Confidence:         0.95,
		}
	}

	// PRIORITY 2: Explicit verification intent
	if gs.hasVerificationIntent(messageLower) {
		return GroundingDecision{
			ShouldUseGrounding: true,
			Reason:             "verification_intent",
			Confidence:         0.9,
		}
	}

	// PRIORITY 3: Recency/new products
	if gs.hasRecencyIntent(messageLower) {
		return GroundingDecision{
			ShouldUseGrounding: true,
			Reason:             "recency_intent",
			Confidence:         0.85,
		}
	}

	// Exclude simple dialogue
	if gs.isSimpleDialogue(messageLower, wordCount) {
		return GroundingDecision{
			ShouldUseGrounding: false,
			Reason:             "simple_dialogue",
			Confidence:         0.95,
		}
	}

	// PRIORITY 4: Advanced dialogue stage
	if gs.isAdvancedDialogueStage(history) {
		return GroundingDecision{
			ShouldUseGrounding: true,
			Reason:             "advanced_dialogue_stage",
			Confidence:         0.7,
		}
	}

	// PRIORITY 5: Technical specs (if enabled)
	if gs.config.EnableTechnicalSpec && gs.hasTechnicalSpecsIntent(messageLower) {
		return GroundingDecision{
			ShouldUseGrounding: true,
			Reason:             "technical_specs",
			Confidence:         0.6,
		}
	}

	// Default: no grounding for general queries
	return GroundingDecision{
		ShouldUseGrounding: false,
		Reason:             "general_query",
		Confidence:         0.75,
	}
}

// 🆕 NEW: Detect when user selected brand and we need latest models
func (gs *GroundingStrategy) isBrandSelectionForLatestModels(messageLower string, history []map[string]string) bool {
	// Check if this is a brand selection response
	majorBrands := []string{
		"apple", "iphone", "samsung", "google", "pixel",
		"xiaomi", "oneplus", "sony", "lg", "huawei",
	}

	isBrandMention := false
	for _, brand := range majorBrands {
		if strings.Contains(messageLower, brand) {
			isBrandMention = true
			break
		}
	}

	if !isBrandMention {
		return false
	}

	// Check conversation history for "which brand" pattern
	if len(history) > 0 {
		lastAssistantMessage := ""
		for i := len(history) - 1; i >= 0; i-- {
			if history[i]["role"] == "assistant" {
				lastAssistantMessage = strings.ToLower(history[i]["content"])
				break
			}
		}

		// If last question was about brand/model, use grounding for latest info
		brandQuestionPatterns := []string{
			"which brand", "what brand", "which iphone", "which model",
			"which samsung", "which google", "which phone",
		}

		for _, pattern := range brandQuestionPatterns {
			if strings.Contains(lastAssistantMessage, pattern) {
				return true
			}
		}
	}

	return false
}

// hasSpecificModelPattern checks for specific product models
func (gs *GroundingStrategy) hasSpecificModelPattern(messageLower string, wordCount int) bool {
	if wordCount < gs.config.MinWordsForProduct {
		return false
	}

	// Known brand + model patterns
	knownBrandModels := []string{
		// Apple
		"iphone 15", "iphone 16", "iphone 17", "iphone 14", "iphone 13",
		"macbook pro", "macbook air", "ipad pro", "ipad air",
		"apple watch", "airpods pro", "airpods max",

		// Samsung
		"galaxy s24", "galaxy s25", "galaxy s26", "galaxy s23", "galaxy a54",
		"galaxy z fold", "galaxy z flip", "galaxy tab", "galaxy watch",

		// Google
		"pixel 8", "pixel 9", "pixel 10", "pixel 7", "pixel fold",

		// Other brands
		"xiaomi 13", "xiaomi 14", "xiaomi 15", "redmi note",
		"oneplus 11", "oneplus 12", "oneplus 13", "oneplus nord",
		"dell xps", "hp spectre", "lenovo thinkpad",
		"asus rog", "acer predator", "msi gaming",
		"sony wh-", "bose quietcomfort",
	}

	for _, pattern := range knownBrandModels {
		if strings.Contains(messageLower, pattern) {
			return true
		}
	}

	// Check for digits (model numbers)
	hasDigits := false
	for _, char := range messageLower {
		if char >= '0' && char <= '9' {
			hasDigits = true
			break
		}
	}

	if hasDigits && wordCount >= gs.config.MinWordsForProduct && wordCount <= 10 {
		brandIndicators := []string{
			"apple", "iphone", "ipad", "macbook",
			"samsung", "galaxy", "google", "pixel",
			"xiaomi", "redmi", "oneplus",
			"sony", "dell", "hp", "lenovo", "asus",
		}

		for _, brand := range brandIndicators {
			if strings.Contains(messageLower, brand) {
				return true
			}
		}
	}

	// Model suffixes
	modelSuffixes := []string{
		"pro max", "pro", "ultra", "plus", "mini", "lite", "air",
		"edge", "note", "fold", "flip",
	}

	for _, suffix := range modelSuffixes {
		if strings.Contains(messageLower, suffix) && wordCount >= 2 {
			return true
		}
	}

	return false
}

// isSimpleDialogue checks if message is simple dialogue
func (gs *GroundingStrategy) isSimpleDialogue(messageLower string, wordCount int) bool {
	greetings := []string{
		"hello", "hi", "hey", "thanks", "thank you",
		"bye", "goodbye",
	}

	for _, greeting := range greetings {
		if messageLower == greeting || strings.HasPrefix(messageLower, greeting+" ") {
			return true
		}
	}

	// Simple answers
	if wordCount <= 2 {
		simpleAnswers := []string{"yes", "no", "ok", "okay"}
		for _, answer := range simpleAnswers {
			if messageLower == answer {
				return true
			}
		}
	}

	// General questions WITHOUT specifics
	if wordCount <= 6 {
		generalQuestions := []string{
			"want to buy", "looking for", "need", "want",
			"help me choose", "what do you recommend",
		}

		for _, phrase := range generalQuestions {
			if strings.Contains(messageLower, phrase) {
				if !gs.hasSpecificModelPattern(messageLower, wordCount) {
					return true
				}
			}
		}
	}

	return false
}

// hasVerificationIntent checks for verification intent
func (gs *GroundingStrategy) hasVerificationIntent(messageLower string) bool {
	verificationPhrases := []string{
		"is there", "does exist", "is available", "can i buy",
		"has been released", "is out", "already out",
		"in stock", "available now", "latest version",
	}

	for _, phrase := range verificationPhrases {
		if strings.Contains(messageLower, phrase) {
			return true
		}
	}

	return false
}

// hasRecencyIntent checks for requests about new products
func (gs *GroundingStrategy) hasRecencyIntent(messageLower string) bool {
	currentYear := time.Now().Year()

	recencyKeywords := []string{
		"new", "newest", "latest", "recent", "current",
		"2024", "2025", "2026", "2027",
	}

	// Add current year and next year
	recencyKeywords = append(recencyKeywords,
		fmt.Sprintf("%d", currentYear),
		fmt.Sprintf("%d", currentYear+1),
	)

	for _, keyword := range recencyKeywords {
		if strings.Contains(messageLower, keyword) {
			return true
		}
	}

	return false
}

// isAdvancedDialogueStage checks if dialogue is in advanced stage
func (gs *GroundingStrategy) isAdvancedDialogueStage(history []map[string]string) bool {
	minMessages := 4
	if gs.config.Mode == "conservative" {
		minMessages = 6
	} else if gs.config.Mode == "aggressive" {
		minMessages = 3
	}

	if len(history) >= minMessages {
		lookbackCount := 4
		if len(history) < lookbackCount {
			lookbackCount = len(history)
		}

		recentMessages := history[len(history)-lookbackCount:]

		for _, msg := range recentMessages {
			content := strings.ToLower(msg["content"])
			wordCount := len(strings.Fields(content))
			if gs.hasSpecificModelPattern(content, wordCount) {
				return true
			}
		}
	}

	return false
}

// hasTechnicalSpecsIntent checks for technical specification requests
func (gs *GroundingStrategy) hasTechnicalSpecsIntent(messageLower string) bool {
	technicalKeywords := []string{
		"specifications", "specs", "features",
		"processor", "cpu", "memory", "ram", "storage",
		"battery", "display", "screen", "camera",
	}

	for _, keyword := range technicalKeywords {
		if strings.Contains(messageLower, keyword) {
			return true
		}
	}

	return false
}

// RecordDecision records a grounding decision for statistics
func (stats *GroundingStats) RecordDecision(decision GroundingDecision) {
	stats.TotalDecisions++

	if decision.ShouldUseGrounding {
		stats.GroundingEnabled++
	} else {
		stats.GroundingDisabled++
	}

	stats.ReasonCounts[decision.Reason]++

	// Update average confidence
	if stats.TotalDecisions > 0 {
		stats.AverageConfidence = (stats.AverageConfidence*float32(stats.TotalDecisions-1) + decision.Confidence) / float32(stats.TotalDecisions)
	}
}
