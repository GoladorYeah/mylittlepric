/**
 * Stats Model - TypeBox schemas and type definitions
 * Following Elysia Best Practices: namespace for DTOs
 */

import { t } from 'elysia';

export namespace StatsModel {
  // ═══════════════════════════════════════════════════════════
  // HEALTH CHECK
  // ═══════════════════════════════════════════════════════════

  export const healthResponse = t.Object({
    status: t.String(),
    redis: t.String(),
    gemini_keys: t.Object({
      count: t.Number(),
      status: t.String(),
    }),
    serp_keys: t.Object({
      count: t.Number(),
      status: t.String(),
    }),
    grounding: t.Object({
      mode: t.String(),
      enabled: t.Boolean(),
    }),
    embedding: t.Object({
      status: t.String(),
    }),
  });

  export type healthResponse = typeof healthResponse.static;

  // ═══════════════════════════════════════════════════════════
  // KEY STATISTICS
  // ═══════════════════════════════════════════════════════════

  export const keyStats = t.Object({
    key_index: t.Number(),
    total_usage: t.String(),
    success_count: t.String(),
    failure_count: t.String(),
    avg_response_time_ms: t.Optional(t.Number()),
  });

  export type keyStats = typeof keyStats.static;

  export const keyStatsResponse = t.Object({
    gemini: t.Array(keyStats),
    serp: t.Array(keyStats),
  });

  export type keyStatsResponse = typeof keyStatsResponse.static;

  // ═══════════════════════════════════════════════════════════
  // GROUNDING STATISTICS
  // ═══════════════════════════════════════════════════════════

  export const groundingStatsResponse = t.Object({
    total_decisions: t.Number(),
    grounding_enabled: t.Number(),
    grounding_disabled: t.Number(),
    grounding_percentage: t.String(),
    reason_breakdown: t.Record(t.String(), t.Number()),
    average_confidence: t.String(),
    mode: t.String(),
    config: t.Object({
      enabled: t.Boolean(),
      min_words: t.Number(),
    }),
  });

  export type groundingStatsResponse = typeof groundingStatsResponse.static;

  // ═══════════════════════════════════════════════════════════
  // TOKEN STATISTICS
  // ═══════════════════════════════════════════════════════════

  export const tokenStatsResponse = t.Object({
    token_usage: t.Object({
      total_requests: t.Number(),
      total_input_tokens: t.Number(),
      total_output_tokens: t.Number(),
      total_tokens: t.Number(),
      requests_with_grounding: t.Number(),
      average_input_tokens: t.String(),
      average_output_tokens: t.String(),
    }),
    timestamp: t.Date(),
  });

  export type tokenStatsResponse = typeof tokenStatsResponse.static;

  // ═══════════════════════════════════════════════════════════
  // ALL STATISTICS
  // ═══════════════════════════════════════════════════════════

  export const allStatsResponse = t.Object({
    health: healthResponse,
    api_keys: keyStatsResponse,
    grounding: t.Omit(groundingStatsResponse, ['config']),
    tokens: t.Object({
      total_requests: t.Number(),
      total_input_tokens: t.Number(),
      total_output_tokens: t.Number(),
      total_tokens: t.Number(),
      requests_with_grounding: t.Number(),
      average_input_tokens: t.String(),
      average_output_tokens: t.String(),
    }),
    timestamp: t.Date(),
  });

  export type allStatsResponse = typeof allStatsResponse.static;

  // ═══════════════════════════════════════════════════════════
  // ERROR RESPONSES
  // ═══════════════════════════════════════════════════════════

  export const errorResponse = t.Object({
    error: t.Boolean(),
    message: t.String(),
  });

  export type errorResponse = typeof errorResponse.static;
}
