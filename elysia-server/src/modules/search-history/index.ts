/**
 * Search History module - handles search history endpoints
 * All endpoints require authentication
 */

import { Elysia, t } from 'elysia';
import type { Container } from '../../container';
import { extractAuth } from '../../middleware/auth';

export const searchHistoryModule = (container: Container) => {
  const authExtractor = extractAuth(container.jwtService);

  return new Elysia({ prefix: '/api/search-history' })
    // Get search history
    .get(
      '/',
      async ({ query, headers, set }) => {
        try {
          // Extract auth context
          const auth = await authExtractor(headers);

          if (!auth.userId) {
            set.status = 401;
            return {
              error: 'UNAUTHORIZED',
              message: 'Authentication required',
            };
          }

          const limit = parseInt(query.limit || '20', 10);
          const offset = parseInt(query.offset || '0', 10);
          const sessionId = query.session_id || undefined;

          // Get search history
          const history = await container.searchHistoryService.getUserSearchHistory(
            auth.userId,
            sessionId,
            limit,
            offset
          );

          return history;
        } catch (error: any) {
          console.error('Error getting search history:', error);
          set.status = 500;
          return {
            error: 'HISTORY_FETCH_ERROR',
            message: 'Failed to retrieve search history',
          };
        }
      },
      {
        query: t.Object({
          limit: t.Optional(t.String()),
          offset: t.Optional(t.String()),
          session_id: t.Optional(t.String()),
        }),
      }
    )

    // Delete specific search history entry
    .delete(
      '/:id',
      async ({ params, headers, set }) => {
        try {
          // Extract auth context
          const auth = await authExtractor(headers);

          if (!auth.userId) {
            set.status = 401;
            return {
              error: 'UNAUTHORIZED',
              message: 'Authentication required',
            };
          }

          const historyId = params.id;

          // Delete the history entry
          await container.searchHistoryService.deleteSearchHistory(
            historyId,
            auth.userId
          );

          return {
            success: true,
            message: 'Search history deleted successfully',
          };
        } catch (error: any) {
          console.error('Error deleting search history:', error);
          set.status = 500;
          return {
            error: 'DELETE_ERROR',
            message: error.message || 'Failed to delete search history',
          };
        }
      },
      {
        params: t.Object({
          id: t.String(),
        }),
      }
    )

    // Delete all search history for user
    .delete('/', async ({ headers, set }) => {
      try {
        // Extract auth context
        const auth = await authExtractor(headers);

        if (!auth.userId) {
          set.status = 401;
          return {
            error: 'UNAUTHORIZED',
            message: 'Authentication required to delete all history',
          };
        }

        // Delete all history
        await container.searchHistoryService.deleteAllUserSearchHistory(auth.userId);

        return {
          success: true,
          message: 'All search history deleted successfully',
        };
      } catch (error: any) {
        console.error('Error deleting all search history:', error);
        set.status = 500;
        return {
          error: 'DELETE_ERROR',
          message: 'Failed to delete search history',
        };
      }
    })

    // Track product click
    .post(
      '/:id/click',
      async ({ params, body, headers, set }) => {
        try {
          // Extract auth context
          const auth = await authExtractor(headers);

          if (!auth.userId) {
            set.status = 401;
            return {
              error: 'UNAUTHORIZED',
              message: 'Authentication required',
            };
          }

          const historyId = params.id;
          const { product_id } = body;

          // Update clicked product
          await container.searchHistoryService.updateClickedProduct(
            historyId,
            product_id
          );

          return {
            success: true,
            message: 'Product click tracked successfully',
          };
        } catch (error: any) {
          console.error('Error tracking product click:', error);
          set.status = 500;
          return {
            error: 'UPDATE_ERROR',
            message: 'Failed to track product click',
          };
        }
      },
      {
        params: t.Object({
          id: t.String(),
        }),
        body: t.Object({
          product_id: t.String(),
        }),
      }
    );
};
