/**
 * GroundingStrategy determines when to use Google Search grounding
 * Based on query analysis, conversation context, and category
 */

import type { Config } from '../config';
import type { EmbeddingService } from './embedding.service';
import type { Message } from '../types';
import { cosineSimilarity } from '../utils/math';

export interface GroundingDecision {
  useGrounding: boolean;
  confidence: number;
  reason: string;
}

export class GroundingStrategy {
  private embedding: EmbeddingService;
  private config: Config;

  constructor(embedding: EmbeddingService, config: Config) {
    this.embedding = embedding;
    this.config = config;
  }

  async shouldUseGrounding(
    userMessage: string,
    history: Message[],
    category: string
  ): Promise<GroundingDecision> {
    const queryEmbedding = await this.embedding.getQueryEmbedding(userMessage);

    if (!queryEmbedding) {
      return {
        useGrounding: false,
        confidence: 0,
        reason: 'embedding_failed',
      };
    }

    // Check if brand-only query
    if (await this.isBrandOnlyQuery(userMessage)) {
      return {
        useGrounding: true,
        confidence: this.config.geminiBrandQueryConfidence,
        reason: 'brand_only_query_vector',
      };
    }

    // Calculate various scores
    const freshInfoScore = await this.calculateFreshInfoSimilarity(
      queryEmbedding
    );
    const specificProductScore = await this.calculateSpecificProductSimilarity(
      queryEmbedding
    );
    const dialogueDriftScore = await this.calculateDialogueDrift(
      queryEmbedding,
      history
    );
    const electronicsScore = await this.calculateCategorySimilarity(
      queryEmbedding,
      'electronics'
    );

    // Weighted total score
    const totalScore =
      freshInfoScore * this.config.geminiGroundingWeightFreshInfo +
      specificProductScore * this.config.geminiGroundingWeightSpecific +
      dialogueDriftScore * this.config.geminiGroundingWeightDrift +
      electronicsScore * this.config.geminiGroundingWeightElectron;

    const useGrounding = totalScore > this.config.geminiGroundingDecisionThresh;
    const reason = this.determineReason(
      freshInfoScore,
      specificProductScore,
      dialogueDriftScore,
      electronicsScore
    );

    return {
      useGrounding,
      confidence: totalScore,
      reason,
    };
  }

  private async isBrandOnlyQuery(userMessage: string): Promise<boolean> {
    const msgLower = userMessage.toLowerCase().trim();
    const words = msgLower.split(/\s+/);

    if (words.length > this.config.geminiBrandQueryMaxWords) {
      return false;
    }

    const brandConcept = await this.embedding.getQueryEmbedding(
      'samsung apple sony xiaomi lg dell hp lenovo asus oppo oneplus realme vivo popular electronics brand manufacturer company'
    );

    const productConcept = await this.embedding.getQueryEmbedding(
      'laptop phone tv computer monitor headphones product type category general question need want'
    );

    const queryEmbedding = await this.embedding.getQueryEmbedding(userMessage);

    if (!brandConcept || !productConcept || !queryEmbedding) {
      return false;
    }

    const brandSimilarity = cosineSimilarity(queryEmbedding, brandConcept);
    const productSimilarity = cosineSimilarity(queryEmbedding, productConcept);

    return (
      brandSimilarity > this.config.geminiBrandSimilarityThresh &&
      brandSimilarity > productSimilarity
    );
  }

  private async calculateFreshInfoSimilarity(
    queryEmbedding: number[]
  ): Promise<number> {
    const freshInfoPatterns = [
      'latest newest current recent 2024 2025 model new release updated',
      'последний новый актуальный свежий модель релиз обновленный',
      'what is the newest what is the latest show me current',
    ];

    let maxSimilarity = 0;

    for (const pattern of freshInfoPatterns) {
      const patternEmbedding = await this.embedding.getQueryEmbedding(pattern);
      if (patternEmbedding) {
        const similarity = cosineSimilarity(queryEmbedding, patternEmbedding);
        if (similarity > maxSimilarity) {
          maxSimilarity = similarity;
        }
      }
    }

    return maxSimilarity;
  }

  private async calculateSpecificProductSimilarity(
    queryEmbedding: number[]
  ): Promise<number> {
    const specificPatterns = [
      'Samsung Galaxy S24 Ultra Apple iPhone 16 Pro Dell XPS 13 Sony TV LG OLED',
      'specific model number brand name exact product full name',
      'конкретная модель номер бренд точный продукт полное название',
    ];

    let maxSimilarity = 0;

    for (const pattern of specificPatterns) {
      const patternEmbedding = await this.embedding.getQueryEmbedding(pattern);
      if (patternEmbedding) {
        const similarity = cosineSimilarity(queryEmbedding, patternEmbedding);
        if (similarity > maxSimilarity) {
          maxSimilarity = similarity;
        }
      }
    }

    return maxSimilarity;
  }

  private async calculateDialogueDrift(
    queryEmbedding: number[],
    history: Message[]
  ): Promise<number> {
    if (history.length < this.config.geminiDialogueHistoryWindow) {
      return 0;
    }

    const windowSize = this.config.geminiDialogueHistoryWindow;
    const recentMessages = history
      .slice(-windowSize)
      .filter((msg) => msg.role === 'user')
      .map((msg) => msg.content);

    if (recentMessages.length === 0) {
      return 0;
    }

    const combinedHistory = recentMessages.join(' ');
    const historyEmbedding = await this.embedding.getQueryEmbedding(
      combinedHistory
    );

    if (!historyEmbedding) {
      return 0;
    }

    const similarity = cosineSimilarity(queryEmbedding, historyEmbedding);
    const drift = 1.0 - similarity;

    if (drift > this.config.geminiDialogueDriftThresh) {
      return this.config.geminiDriftScoreBonus;
    }

    return 0;
  }

  private async calculateCategorySimilarity(
    queryEmbedding: number[],
    category: string
  ): Promise<number> {
    // This assumes embedding service has category embeddings loaded
    const categoryEmbedding = await this.embedding.getQueryEmbedding(category);

    if (!categoryEmbedding) {
      return 0;
    }

    const similarity = cosineSimilarity(queryEmbedding, categoryEmbedding);

    if (
      category === 'electronics' &&
      similarity > this.config.geminiElectronicsThreshHigh
    ) {
      return this.config.geminiElectronicsScoreHigh;
    }

    if (similarity > this.config.geminiCategorySimilarityThresh) {
      return this.config.geminiCategoryScore;
    }

    return 0;
  }

  private determineReason(
    fresh: number,
    specific: number,
    drift: number,
    electronics: number
  ): string {
    const scores: Record<string, number> = {
      fresh_info_semantic: fresh,
      specific_product_semantic: specific,
      dialogue_drift_detected: drift,
      electronics_category: electronics,
    };

    let maxScore = -1;
    let reason = 'unknown';

    for (const [key, value] of Object.entries(scores)) {
      if (value > maxScore) {
        maxScore = value;
        reason = key;
      }
    }

    return reason;
  }
}
