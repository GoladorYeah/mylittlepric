/**
 * Search History Model - TypeBox schemas and type definitions
 * Following Elysia Best Practices: namespace for DTOs
 */

import { t } from 'elysia';

export namespace SearchHistoryModel {
  // ═══════════════════════════════════════════════════════════
  // GET SEARCH HISTORY
  // ═══════════════════════════════════════════════════════════

  // GET /api/search-history - Query params
  export const getHistoryQuery = t.Object({
    limit: t.Optional(t.String()),
    offset: t.Optional(t.String()),
    session_id: t.Optional(t.String()),
  });

  export type getHistoryQuery = typeof getHistoryQuery.static;

  // Search history item
  export const historyItem = t.Object({
    id: t.String(),
    user_id: t.String(),
    session_id: t.String(),
    query: t.String(),
    category: t.String(),
    search_type: t.String(),
    country: t.String(),
    language: t.String(),
    result_count: t.Number(),
    created_at: t.Date(),
  });

  export type historyItem = typeof historyItem.static;

  // GET /api/search-history - Response
  export const getHistoryResponse = t.Array(historyItem);

  export type getHistoryResponse = typeof getHistoryResponse.static;

  // ═══════════════════════════════════════════════════════════
  // DELETE HISTORY
  // ═══════════════════════════════════════════════════════════

  // DELETE /api/search-history/:id - Params
  export const deleteHistoryParams = t.Object({
    id: t.String(),
  });

  export type deleteHistoryParams = typeof deleteHistoryParams.static;

  // Success response
  export const successResponse = t.Object({
    success: t.Boolean(),
    message: t.String(),
  });

  export type successResponse = typeof successResponse.static;

  // ═══════════════════════════════════════════════════════════
  // TRACK PRODUCT CLICK
  // ═══════════════════════════════════════════════════════════

  // POST /api/search-history/:id/click - Body
  export const trackClickBody = t.Object({
    product_id: t.String(),
  });

  export type trackClickBody = typeof trackClickBody.static;

  // ═══════════════════════════════════════════════════════════
  // ERROR RESPONSES
  // ═══════════════════════════════════════════════════════════

  export const errorResponse = t.Object({
    error: t.String(),
    message: t.String(),
  });

  export type errorResponse = typeof errorResponse.static;

  export const unauthorizedError = t.Literal('UNAUTHORIZED');
  export type unauthorizedError = typeof unauthorizedError.static;

  export const fetchError = t.Literal('HISTORY_FETCH_ERROR');
  export type fetchError = typeof fetchError.static;

  export const deleteError = t.Literal('DELETE_ERROR');
  export type deleteError = typeof deleteError.static;

  export const updateError = t.Literal('UPDATE_ERROR');
  export type updateError = typeof updateError.static;
}
