// backend/internal/services/grounding_strategy.go
package services

import (
	"strings"

	"mylittleprice/internal/config"
)

type GroundingDecision struct {
	UseGrounding bool
	Confidence   float32
	Reason       string
}

type GroundingStrategy struct {
	embedding *EmbeddingService
	config    *config.Config
}

func NewGroundingStrategy(embedding *EmbeddingService, cfg *config.Config) *GroundingStrategy {
	return &GroundingStrategy{
		embedding: embedding,
		config:    cfg,
	}
}

func (gs *GroundingStrategy) ShouldUseGrounding(
	userMessage string,
	history []map[string]string,
	category string,
) GroundingDecision {

	queryEmbedding := gs.embedding.GetQueryEmbedding(userMessage)
	if queryEmbedding == nil {
		return GroundingDecision{
			UseGrounding: false,
			Confidence:   0.0,
			Reason:       "embedding_failed",
		}
	}

	if gs.isBrandOnlyQuery(userMessage) {
		return GroundingDecision{
			UseGrounding: true,
			Confidence:   float32(gs.config.GeminiBrandQueryConfidence),
			Reason:       "brand_only_query_vector",
		}
	}

	freshInfoScore := gs.calculateFreshInfoSimilarity(queryEmbedding)
	specificProductScore := gs.calculateSpecificProductSimilarity(queryEmbedding)
	dialogueDriftScore := gs.calculateDialogueDrift(queryEmbedding, history)
	electronicsScore := gs.calculateCategorySimilarity(queryEmbedding, "electronics")

	totalScore := (freshInfoScore * float32(gs.config.GeminiGroundingWeightFreshInfo)) +
		(specificProductScore * float32(gs.config.GeminiGroundingWeightSpecific)) +
		(dialogueDriftScore * float32(gs.config.GeminiGroundingWeightDrift)) +
		(electronicsScore * float32(gs.config.GeminiGroundingWeightElectron))

	useGrounding := totalScore > float32(gs.config.GeminiGroundingDecisionThresh)
	reason := gs.determineReason(freshInfoScore, specificProductScore, dialogueDriftScore, electronicsScore)

	return GroundingDecision{
		UseGrounding: useGrounding,
		Confidence:   totalScore,
		Reason:       reason,
	}
}

func (gs *GroundingStrategy) isBrandOnlyQuery(userMessage string) bool {
	msgLower := strings.ToLower(strings.TrimSpace(userMessage))
	words := strings.Fields(msgLower)

	if len(words) > gs.config.GeminiBrandQueryMaxWords {
		return false
	}

	brandConcept := gs.embedding.GetQueryEmbedding(
		"samsung apple sony xiaomi lg dell hp lenovo asus oppo oneplus realme vivo popular electronics brand manufacturer company",
	)

	productConcept := gs.embedding.GetQueryEmbedding(
		"laptop phone tv computer monitor headphones product type category general question need want",
	)

	queryEmbedding := gs.embedding.GetQueryEmbedding(userMessage)

	if brandConcept == nil || productConcept == nil || queryEmbedding == nil {
		return false
	}

	brandSimilarity := cosineSimilarity(queryEmbedding, brandConcept)
	productSimilarity := cosineSimilarity(queryEmbedding, productConcept)

	return brandSimilarity > float32(gs.config.GeminiBrandSimilarityThresh) && brandSimilarity > productSimilarity
}

func (gs *GroundingStrategy) calculateFreshInfoSimilarity(queryEmbedding []float32) float32 {
	freshInfoPatterns := []string{
		"latest newest current recent 2024 2025 model new release updated",
		"последний новый актуальный свежий модель релиз обновленный",
		"what is the newest what is the latest show me current",
	}

	maxSimilarity := float32(0.0)
	for _, pattern := range freshInfoPatterns {
		patternEmbedding := gs.embedding.GetQueryEmbedding(pattern)
		if patternEmbedding != nil {
			similarity := cosineSimilarity(queryEmbedding, patternEmbedding)
			if similarity > maxSimilarity {
				maxSimilarity = similarity
			}
		}
	}

	return maxSimilarity
}

func (gs *GroundingStrategy) calculateSpecificProductSimilarity(queryEmbedding []float32) float32 {
	specificPatterns := []string{
		"Samsung Galaxy S24 Ultra Apple iPhone 16 Pro Dell XPS 13 Sony TV LG OLED",
		"specific model number brand name exact product full name",
		"конкретная модель номер бренд точный продукт полное название",
	}

	maxSimilarity := float32(0.0)
	for _, pattern := range specificPatterns {
		patternEmbedding := gs.embedding.GetQueryEmbedding(pattern)
		if patternEmbedding != nil {
			similarity := cosineSimilarity(queryEmbedding, patternEmbedding)
			if similarity > maxSimilarity {
				maxSimilarity = similarity
			}
		}
	}

	return maxSimilarity
}

func (gs *GroundingStrategy) calculateDialogueDrift(queryEmbedding []float32, history []map[string]string) float32 {
	if len(history) < gs.config.GeminiDialogueHistoryWindow {
		return 0.0
	}

	recentMessages := []string{}
	windowSize := gs.config.GeminiDialogueHistoryWindow
	for i := len(history) - windowSize; i < len(history); i++ {
		if history[i]["role"] == "user" {
			recentMessages = append(recentMessages, history[i]["content"])
		}
	}

	if len(recentMessages) == 0 {
		return 0.0
	}

	combinedHistory := strings.Join(recentMessages, " ")
	historyEmbedding := gs.embedding.GetQueryEmbedding(combinedHistory)

	if historyEmbedding == nil {
		return 0.0
	}

	similarity := cosineSimilarity(queryEmbedding, historyEmbedding)
	drift := 1.0 - similarity

	if drift > float32(gs.config.GeminiDialogueDriftThresh) {
		return float32(gs.config.GeminiDriftScoreBonus)
	}

	return 0.0
}

func (gs *GroundingStrategy) calculateCategorySimilarity(queryEmbedding []float32, category string) float32 {
	categoryEmbedding := gs.embedding.categoryEmbeddings[category]
	if categoryEmbedding == nil {
		return 0.0
	}

	similarity := cosineSimilarity(queryEmbedding, categoryEmbedding)

	if category == "electronics" && similarity > float32(gs.config.GeminiElectronicsThreshHigh) {
		return float32(gs.config.GeminiElectronicsScoreHigh)
	}

	if similarity > float32(gs.config.GeminiCategorySimilarityThresh) {
		return float32(gs.config.GeminiCategoryScore)
	}

	return 0.0
}

func (gs *GroundingStrategy) determineReason(fresh, specific, drift, electronics float32) string {
	scores := map[string]float32{
		"fresh_info_semantic":       fresh,
		"specific_product_semantic": specific,
		"dialogue_drift_detected":   drift,
		"electronics_category":      electronics,
	}

	maxReason := ""
	maxScore := float32(0.0)

	for reason, score := range scores {
		if score > maxScore {
			maxScore = score
			maxReason = reason
		}
	}

	if maxReason == "" {
		return "vector_threshold_not_met"
	}

	return maxReason
}
