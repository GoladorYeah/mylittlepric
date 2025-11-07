/**
 * EmbeddingService handles text embeddings using Google's Gemini API
 * Used for category detection and semantic similarity matching
 */

import { GoogleGenerativeAI } from '@google/generative-ai';
import type { Redis } from 'ioredis';
import type { Config } from '../config';
import { cosineSimilarity } from '../utils/math';

export class EmbeddingService {
  private genai: GoogleGenerativeAI;
  private redis: Redis;
  private config: Config;
  private categoryEmbeddings: Map<string, number[]> = new Map();

  constructor(apiKey: string, redis: Redis, config: Config) {
    this.genai = new GoogleGenerativeAI(apiKey);
    this.redis = redis;
    this.config = config;
    this.loadCategoryEmbeddings();
  }

  /**
   * Load category embeddings from Redis or generate new ones
   */
  private async loadCategoryEmbeddings(): Promise<void> {
    const key = 'embeddings:categories:v1';

    try {
      const data = await this.redis.get(key);

      if (data) {
        const parsed = JSON.parse(data);
        this.categoryEmbeddings = new Map(Object.entries(parsed));
        console.log('✅ Loaded category embeddings from cache');
      } else {
        await this.generateCategoryEmbeddings();

        const serialized = JSON.stringify(
          Object.fromEntries(this.categoryEmbeddings)
        );
        await this.redis.set(key, serialized);
        console.log('✅ Generated and cached category embeddings');
      }
    } catch (error) {
      console.error('⚠️ Failed to load category embeddings:', error);
      await this.generateCategoryEmbeddings();
    }
  }

  /**
   * Generate embeddings for predefined product categories
   */
  private async generateCategoryEmbeddings(): Promise<void> {
    const categories: Record<string, string> = {
      electronics: 'laptop computer phone tablet tv monitor camera headphones speaker gadget electronics device',
      clothing: 'shirt pants dress shoes jacket coat sweater jeans clothing fashion apparel wear',
      furniture: 'chair table bed sofa desk cabinet shelf bookcase furniture home decor',
      kitchen: 'pan pot knife plate cup dish spoon fork cookware kitchen utensil appliance',
      sports: 'bicycle ball racket fitness gym equipment sports workout training exercise',
      tools: 'drill hammer screwdriver wrench saw power tool hand tool equipment',
      decor: 'lamp vase picture frame mirror decoration ornament home decor',
      textiles: 'pillow blanket towel sheet curtain textile fabric bedding linen',
    };

    for (const [category, text] of Object.entries(categories)) {
      const embedding = await this.getEmbedding(text);
      if (embedding) {
        this.categoryEmbeddings.set(category, embedding);
      }
    }
  }

  /**
   * Get embedding for a text using Gemini API
   */
  private async getEmbedding(text: string): Promise<number[] | null> {
    try {
      const model = this.genai.getGenerativeModel({
        model: this.config.geminiEmbeddingModel,
      });

      const result = await model.embedContent(text);
      return result.embedding.values;
    } catch (error) {
      console.error('Error getting embedding:', error);
      return null;
    }
  }

  /**
   * Get embedding for a query (with caching)
   */
  async getQueryEmbedding(query: string): Promise<number[] | null> {
    const cacheKey = `embeddings:query:${query}`;

    try {
      // Try to get from cache
      const cached = await this.redis.get(cacheKey);
      if (cached) {
        return JSON.parse(cached);
      }

      // Generate new embedding
      const embedding = await this.getEmbedding(query);
      if (embedding) {
        // Cache the embedding
        await this.redis.setex(
          cacheKey,
          this.config.cacheQueryEmbeddingTTL,
          JSON.stringify(embedding)
        );
      }

      return embedding;
    } catch (error) {
      console.error('Error getting query embedding:', error);
      return null;
    }
  }

  /**
   * Detect product category based on user message
   */
  async detectCategory(userMessage: string): Promise<string> {
    const queryEmbedding = await this.getQueryEmbedding(userMessage);
    if (!queryEmbedding) {
      return '';
    }

    let maxSimilarity = -1;
    let bestCategory = '';

    for (const [category, categoryEmbedding] of this.categoryEmbeddings) {
      const similarity = cosineSimilarity(queryEmbedding, categoryEmbedding);
      if (similarity > maxSimilarity) {
        maxSimilarity = similarity;
        bestCategory = category;
      }
    }

    if (maxSimilarity > this.config.embeddingCategoryDetectionThresh) {
      return bestCategory;
    }

    return '';
  }

  /**
   * Find similar cached query
   */
  async findSimilarCachedQuery(
    query: string,
    threshold: number
  ): Promise<string> {
    const queryEmbedding = await this.getQueryEmbedding(query);
    if (!queryEmbedding) {
      return '';
    }

    try {
      const pattern = 'cache:search:*';
      const maxKeysToCheck = 100;
      let keysChecked = 0;

      // Get keys matching pattern
      const keys = await this.redis.keys(pattern);

      for (const cacheKey of keys) {
        if (keysChecked >= maxKeysToCheck) {
          console.log(
            `⚠️ FindSimilarCachedQuery: Reached max keys limit (${maxKeysToCheck}), stopping scan`
          );
          break;
        }

        const cachedQuery = cacheKey.substring('cache:search:'.length);
        const cachedEmbedding = await this.getQueryEmbedding(cachedQuery);

        if (cachedEmbedding) {
          const similarity = cosineSimilarity(queryEmbedding, cachedEmbedding);
          if (similarity >= threshold) {
            return cacheKey;
          }
        }

        keysChecked++;
      }
    } catch (error) {
      console.error('⚠️ FindSimilarCachedQuery: Error:', error);
    }

    return '';
  }

  /**
   * Check if two product names are duplicates based on embedding similarity
   */
  async areDuplicateProducts(
    name1: string,
    name2: string,
    threshold: number
  ): Promise<boolean> {
    const emb1 = await this.getQueryEmbedding(name1);
    const emb2 = await this.getQueryEmbedding(name2);

    if (!emb1 || !emb2) {
      return false;
    }

    const similarity = cosineSimilarity(emb1, emb2);
    return similarity >= threshold;
  }
}
