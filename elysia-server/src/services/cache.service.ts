/**
 * CacheService handles caching for search results, products, and Gemini responses
 */

import type { Redis } from 'ioredis';
import type { Config } from '../config';
import type { ProductCard, GeminiResponse } from '../types';
import type { EmbeddingService } from './embedding.service';
import { cosineSimilarity } from '../utils/math';

export class CacheService {
  private redis: Redis;
  private config: Config;
  private embedding: EmbeddingService;

  constructor(redis: Redis, config: Config, embedding: EmbeddingService) {
    this.redis = redis;
    this.config = config;
    this.embedding = embedding;
  }

  /**
   * Get search results from cache
   */
  async getSearchResults(cacheKey: string): Promise<ProductCard[] | null> {
    try {
      const data = await this.redis.get(cacheKey);

      if (data) {
        return JSON.parse(data);
      }

      // Try to find similar cached query
      const similarKey = await this.embedding.findSimilarCachedQuery(
        cacheKey,
        0.92
      );

      if (similarKey) {
        const similarData = await this.redis.get(similarKey);
        if (similarData) {
          return JSON.parse(similarData);
        }
      }

      return null;
    } catch (error) {
      console.error('Error getting search results from cache:', error);
      return null;
    }
  }

  /**
   * Set search results in cache
   */
  async setSearchResults(
    cacheKey: string,
    cards: ProductCard[],
    ttl: number
  ): Promise<void> {
    try {
      const dedupedCards = await this.deduplicateProducts(cards);
      await this.redis.setex(cacheKey, ttl, JSON.stringify(dedupedCards));
    } catch (error) {
      console.error('Error setting search results in cache:', error);
    }
  }

  /**
   * Deduplicate products using embeddings
   */
  private async deduplicateProducts(
    cards: ProductCard[]
  ): Promise<ProductCard[]> {
    if (cards.length <= 1) {
      return cards;
    }

    interface CardWithEmbedding {
      card: ProductCard;
      embedding: number[];
    }

    // Pre-compute embeddings for all products
    const cardsWithEmbeddings: CardWithEmbedding[] = [];

    for (const card of cards) {
      const embedding = await this.embedding.getQueryEmbedding(card.name);
      if (embedding) {
        cardsWithEmbeddings.push({ card, embedding });
      }
    }

    if (cardsWithEmbeddings.length === 0) {
      return cards; // Fallback: return all if embeddings failed
    }

    // Compare embeddings to find duplicates
    const unique: CardWithEmbedding[] = [cardsWithEmbeddings[0]];

    for (let i = 1; i < cardsWithEmbeddings.length; i++) {
      let isDuplicate = false;

      for (const uniqueCard of unique) {
        const similarity = cosineSimilarity(
          cardsWithEmbeddings[i].embedding,
          uniqueCard.embedding
        );

        if (similarity >= 0.95) {
          isDuplicate = true;
          break;
        }
      }

      if (!isDuplicate) {
        unique.push(cardsWithEmbeddings[i]);
      }
    }

    return unique.map((item) => item.card);
  }

  /**
   * Get product details by page token
   */
  async getProductByToken(pageToken: string): Promise<any | null> {
    const cacheKey = `product:${pageToken}`;

    try {
      const data = await this.redis.get(cacheKey);
      return data ? JSON.parse(data) : null;
    } catch (error) {
      console.error('Error getting product by token:', error);
      return null;
    }
  }

  /**
   * Set product details by page token
   */
  async setProductByToken(
    pageToken: string,
    product: any,
    ttl: number
  ): Promise<void> {
    const cacheKey = `product:${pageToken}`;

    try {
      await this.redis.setex(cacheKey, ttl, JSON.stringify(product));
    } catch (error) {
      console.error('Error setting product by token:', error);
    }
  }

  /**
   * Get Gemini response from cache
   */
  async getGeminiResponse(cacheKey: string): Promise<GeminiResponse | null> {
    try {
      const data = await this.redis.get(cacheKey);
      return data ? JSON.parse(data) : null;
    } catch (error) {
      console.error('Error getting Gemini response from cache:', error);
      return null;
    }
  }

  /**
   * Set Gemini response in cache
   */
  async setGeminiResponse(
    cacheKey: string,
    response: GeminiResponse
  ): Promise<void> {
    try {
      await this.redis.setex(
        cacheKey,
        this.config.cacheGeminiTTL,
        JSON.stringify(response)
      );
    } catch (error) {
      console.error('Error setting Gemini response in cache:', error);
    }
  }
}
