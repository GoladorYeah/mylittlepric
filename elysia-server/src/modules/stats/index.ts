/**
 * Stats module - handles statistics and health endpoints
 */

import { Elysia } from 'elysia';
import type { Container } from '../../container';

export const statsModule = (container: Container) =>
  new Elysia({ prefix: '/api' })
    // Health check
    .get('/health', async () => {
      return container.healthCheck();
    })

    // Key statistics
    .get('/stats/keys', async () => {
      try {
        const geminiStats = await container.geminiRotator.getAllStats();
        const serpStats = await container.serpRotator.getAllStats();

        return {
          gemini: geminiStats,
          serp: serpStats,
        };
      } catch (error: any) {
        console.error('Stats error:', error);
        return {
          error: true,
          message: 'Failed to get key statistics',
        };
      }
    })

    // Grounding statistics
    .get('/stats/grounding', () => {
      try {
        const stats = container.geminiService.getGroundingStats();

        const groundingPercentage =
          stats.totalDecisions > 0
            ? ((stats.groundingEnabled / stats.totalDecisions) * 100).toFixed(1)
            : '0.0';

        return {
          total_decisions: stats.totalDecisions,
          grounding_enabled: stats.groundingEnabled,
          grounding_disabled: stats.groundingDisabled,
          grounding_percentage: `${groundingPercentage}%`,
          reason_breakdown: stats.reasonCounts,
          average_confidence: stats.averageConfidence.toFixed(2),
          mode: container.config.geminiGroundingMode,
          config: {
            enabled: container.config.geminiUseGrounding,
            min_words: container.config.geminiGroundingMinWords,
          },
        };
      } catch (error: any) {
        console.error('Grounding stats error:', error);
        return {
          error: true,
          message: 'Failed to get grounding statistics',
        };
      }
    })

    // Token statistics
    .get('/stats/tokens', () => {
      try {
        const tokenStats = container.geminiService.getTokenStats();

        return {
          token_usage: {
            total_requests: tokenStats.totalRequests,
            total_input_tokens: tokenStats.totalInputTokens,
            total_output_tokens: tokenStats.totalOutputTokens,
            total_tokens: tokenStats.totalTokens,
            requests_with_grounding: tokenStats.requestsWithGrounding,
            average_input_tokens: tokenStats.averageInputTokens.toFixed(2),
            average_output_tokens: tokenStats.averageOutputTokens.toFixed(2),
          },
          timestamp: new Date(),
        };
      } catch (error: any) {
        console.error('Token stats error:', error);
        return {
          error: true,
          message: 'Failed to get token statistics',
        };
      }
    })

    // All statistics
    .get('/stats/all', async () => {
      try {
        const health = await container.healthCheck();
        const geminiStats = await container.geminiRotator.getAllStats();
        const serpStats = await container.serpRotator.getAllStats();
        const groundingStats = container.geminiService.getGroundingStats();
        const tokenStats = container.geminiService.getTokenStats();

        const groundingPercentage =
          groundingStats.totalDecisions > 0
            ? (
                (groundingStats.groundingEnabled / groundingStats.totalDecisions) *
                100
              ).toFixed(1)
            : '0.0';

        return {
          health,
          api_keys: {
            gemini: geminiStats,
            serp: serpStats,
          },
          grounding: {
            total_decisions: groundingStats.totalDecisions,
            grounding_enabled: groundingStats.groundingEnabled,
            grounding_disabled: groundingStats.groundingDisabled,
            grounding_percentage: `${groundingPercentage}%`,
            reason_breakdown: groundingStats.reasonCounts,
            average_confidence: groundingStats.averageConfidence.toFixed(2),
            mode: container.config.geminiGroundingMode,
          },
          tokens: {
            total_requests: tokenStats.totalRequests,
            total_input_tokens: tokenStats.totalInputTokens,
            total_output_tokens: tokenStats.totalOutputTokens,
            total_tokens: tokenStats.totalTokens,
            requests_with_grounding: tokenStats.requestsWithGrounding,
            average_input_tokens: tokenStats.averageInputTokens.toFixed(2),
            average_output_tokens: tokenStats.averageOutputTokens.toFixed(2),
          },
          timestamp: new Date(),
        };
      } catch (error: any) {
        console.error('Stats error:', error);
        return {
          error: true,
          message: 'Failed to get statistics',
        };
      }
    });
