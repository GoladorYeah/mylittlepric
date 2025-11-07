// backend/internal/services/cache.go
package services

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/redis/go-redis/v9"

	"mylittleprice/internal/config"
	"mylittleprice/internal/models"
)

type CacheService struct {
	redis     *redis.Client
	config    *config.Config
	embedding *EmbeddingService
	ctx       context.Context
}

// NewCacheService creates a new CacheService with injected dependencies
// Following the Dependency Injection pattern used throughout the application
func NewCacheService(redisClient *redis.Client, cfg *config.Config, embedding *EmbeddingService) *CacheService {
	return &CacheService{
		redis:     redisClient,
		config:    cfg,
		embedding: embedding,
		ctx:       context.Background(),
	}
}

func (c *CacheService) GetSearchResults(cacheKey string) ([]models.ProductCard, error) {
	data, err := c.redis.Get(c.ctx, cacheKey).Bytes()
	if err == redis.Nil {
		similarKey := c.embedding.FindSimilarCachedQuery(cacheKey, 0.92)
		if similarKey != "" {
			data, err = c.redis.Get(c.ctx, similarKey).Bytes()
			if err == nil {
				var cards []models.ProductCard
				if err := json.Unmarshal(data, &cards); err == nil {
					return cards, nil
				}
			}
		}
		return nil, fmt.Errorf("cache miss")
	}
	if err != nil {
		return nil, fmt.Errorf("redis error: %w", err)
	}

	var cards []models.ProductCard
	if err := json.Unmarshal(data, &cards); err != nil {
		return nil, fmt.Errorf("unmarshal error: %w", err)
	}

	return cards, nil
}

func (c *CacheService) SetSearchResults(cacheKey string, cards []models.ProductCard, ttl time.Duration) error {
	dedupedCards := c.deduplicateProducts(cards)

	data, err := json.Marshal(dedupedCards)
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}

	return c.redis.Set(c.ctx, cacheKey, data, ttl).Err()
}

func (c *CacheService) deduplicateProducts(cards []models.ProductCard) []models.ProductCard {
	if len(cards) <= 1 {
		return cards
	}

	// Optimized deduplication: compute embeddings once, then compare
	// This avoids O(n²) with expensive API calls inside the loop
	type cardWithEmbedding struct {
		card      models.ProductCard
		embedding []float32
	}

	// Pre-compute embeddings for all products (O(n) API calls)
	cardsWithEmbeddings := make([]cardWithEmbedding, 0, len(cards))
	for _, card := range cards {
		emb := c.embedding.GetQueryEmbedding(card.Name)
		if emb != nil {
			cardsWithEmbeddings = append(cardsWithEmbeddings, cardWithEmbedding{
				card:      card,
				embedding: emb,
			})
		}
	}

	if len(cardsWithEmbeddings) == 0 {
		return cards // Fallback: return all if embeddings failed
	}

	// Now compare embeddings (O(n²) comparisons but no API calls)
	unique := []cardWithEmbedding{cardsWithEmbeddings[0]}

	for i := 1; i < len(cardsWithEmbeddings); i++ {
		isDuplicate := false
		for j := range unique {
			similarity := cosineSimilarity(cardsWithEmbeddings[i].embedding, unique[j].embedding)
			if similarity >= 0.95 {
				isDuplicate = true
				break
			}
		}
		if !isDuplicate {
			unique = append(unique, cardsWithEmbeddings[i])
		}
	}

	// Extract cards from unique list
	result := make([]models.ProductCard, len(unique))
	for i, item := range unique {
		result[i] = item.card
	}

	return result
}

func (c *CacheService) GetProductByToken(pageToken string) (map[string]interface{}, error) {
	cacheKey := fmt.Sprintf("product:%s", pageToken)

	data, err := c.redis.Get(c.ctx, cacheKey).Bytes()
	if err == redis.Nil {
		return nil, fmt.Errorf("cache miss")
	}
	if err != nil {
		return nil, fmt.Errorf("redis error: %w", err)
	}

	var product map[string]interface{}
	if err := json.Unmarshal(data, &product); err != nil {
		return nil, fmt.Errorf("unmarshal error: %w", err)
	}

	return product, nil
}

func (c *CacheService) SetProductByToken(pageToken string, product map[string]interface{}, ttl int) error {
	cacheKey := fmt.Sprintf("product:%s", pageToken)

	data, err := json.Marshal(product)
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}

	duration := time.Duration(ttl) * time.Second
	return c.redis.Set(c.ctx, cacheKey, data, duration).Err()
}

func (c *CacheService) GetGeminiResponse(cacheKey string) (*models.GeminiResponse, error) {
	data, err := c.redis.Get(c.ctx, cacheKey).Bytes()
	if err == redis.Nil {
		return nil, fmt.Errorf("cache miss")
	}
	if err != nil {
		return nil, fmt.Errorf("redis error: %w", err)
	}

	var response models.GeminiResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("unmarshal error: %w", err)
	}

	return &response, nil
}

func (c *CacheService) SetGeminiResponse(cacheKey string, response *models.GeminiResponse) error {
	data, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}

	ttl := time.Duration(c.config.CacheGeminiTTL) * time.Second
	return c.redis.Set(c.ctx, cacheKey, data, ttl).Err()
}

// cosineSimilarity calculates cosine similarity between two embedding vectors
func cosineSimilarity(a, b []float32) float32 {
	if len(a) != len(b) {
		return 0
	}

	var dotProduct, normA, normB float32
	for i := range a {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	if normA == 0 || normB == 0 {
		return 0
	}

	return dotProduct / (float32(math.Sqrt(float64(normA))) * float32(math.Sqrt(float64(normB))))
}
