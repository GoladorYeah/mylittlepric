// backend/internal/services/embedding.go
package services

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"google.golang.org/genai"

	"mylittleprice/internal/config"
)

type EmbeddingService struct {
	client             *genai.Client
	redis              *redis.Client
	config             *config.Config
	ctx                context.Context
	categoryEmbeddings map[string][]float32
	mu                 sync.RWMutex
}

func NewEmbeddingService(client *genai.Client, redis *redis.Client, cfg *config.Config) *EmbeddingService {
	s := &EmbeddingService{
		client:             client,
		redis:              redis,
		config:             cfg,
		ctx:                context.Background(),
		categoryEmbeddings: make(map[string][]float32),
	}
	s.loadCategoryEmbeddings()
	return s
}

func (e *EmbeddingService) loadCategoryEmbeddings() {
	key := "embeddings:categories:v1"
	data, err := e.redis.Get(e.ctx, key).Bytes()

	if err == redis.Nil {
		e.generateCategoryEmbeddings()
		jsonData, _ := json.Marshal(e.categoryEmbeddings)
		e.redis.Set(e.ctx, key, jsonData, 0)
	} else {
		json.Unmarshal(data, &e.categoryEmbeddings)
	}
}

func (e *EmbeddingService) generateCategoryEmbeddings() {
	categories := map[string]string{
		"electronics": "laptop computer phone tablet tv monitor camera headphones speaker gadget electronics device",
		"clothing":    "shirt pants dress shoes jacket coat sweater jeans clothing fashion apparel wear",
		"furniture":   "chair table bed sofa desk cabinet shelf bookcase furniture home decor",
		"kitchen":     "pan pot knife plate cup dish spoon fork cookware kitchen utensil appliance",
		"sports":      "bicycle ball racket fitness gym equipment sports workout training exercise",
		"tools":       "drill hammer screwdriver wrench saw power tool hand tool equipment",
		"decor":       "lamp vase picture frame mirror decoration ornament home decor",
		"textiles":    "pillow blanket towel sheet curtain textile fabric bedding linen",
	}

	for category, text := range categories {
		embedding := e.getEmbedding(text)
		if embedding != nil {
			e.categoryEmbeddings[category] = embedding
		}
	}
}

func (e *EmbeddingService) getEmbedding(text string) []float32 {
	resp, err := e.client.Models.EmbedContent(
		e.ctx,
		e.config.GeminiEmbeddingModel,
		genai.Text(text),
		nil,
	)
	if err != nil || resp == nil || len(resp.Embeddings) == 0 {
		return nil
	}
	return resp.Embeddings[0].Values
}

func (e *EmbeddingService) GetQueryEmbedding(query string) []float32 {
	cacheKey := fmt.Sprintf("embeddings:query:%s", query)
	cached, err := e.redis.Get(e.ctx, cacheKey).Bytes()

	if err == nil {
		var embedding []float32
		json.Unmarshal(cached, &embedding)
		return embedding
	}

	embedding := e.getEmbedding(query)
	if embedding != nil {
		jsonData, _ := json.Marshal(embedding)
		ttl := time.Duration(e.config.CacheQueryEmbeddingTTL) * time.Second
		e.redis.Set(e.ctx, cacheKey, jsonData, ttl)
	}
	return embedding
}

func (e *EmbeddingService) DetectCategory(userMessage string) string {
	queryEmbedding := e.GetQueryEmbedding(userMessage)
	if queryEmbedding == nil {
		return ""
	}

	e.mu.RLock()
	defer e.mu.RUnlock()

	maxSimilarity := float32(-1)
	bestCategory := ""

	for category, categoryEmbedding := range e.categoryEmbeddings {
		similarity := cosineSimilarity(queryEmbedding, categoryEmbedding)
		if similarity > maxSimilarity {
			maxSimilarity = similarity
			bestCategory = category
		}
	}

	if maxSimilarity > float32(e.config.EmbeddingCategoryDetectionThresh) {
		return bestCategory
	}
	return ""
}

func (e *EmbeddingService) FindSimilarCachedQuery(query string, threshold float32) string {
	queryEmbedding := e.GetQueryEmbedding(query)
	if queryEmbedding == nil {
		return ""
	}

	pattern := "cache:search:*"
	iter := e.redis.Scan(e.ctx, 0, pattern, 100).Iterator()

	for iter.Next(e.ctx) {
		cacheKey := iter.Val()

		cachedQuery := cacheKey[len("cache:search:"):]
		cachedEmbedding := e.GetQueryEmbedding(cachedQuery)

		if cachedEmbedding != nil {
			similarity := cosineSimilarity(queryEmbedding, cachedEmbedding)
			if similarity >= threshold {
				return cacheKey
			}
		}
	}

	return ""
}

func (e *EmbeddingService) AreDuplicateProducts(name1, name2 string, threshold float32) bool {
	emb1 := e.GetQueryEmbedding(name1)
	emb2 := e.GetQueryEmbedding(name2)

	if emb1 == nil || emb2 == nil {
		return false
	}

	similarity := cosineSimilarity(emb1, emb2)
	return similarity >= threshold
}

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
