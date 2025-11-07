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

    // All statistics
    .get('/stats/all', async () => {
      try {
        const health = await container.healthCheck();
        const geminiStats = await container.geminiRotator.getAllStats();
        const serpStats = await container.serpRotator.getAllStats();

        return {
          health,
          keys: {
            gemini: geminiStats,
            serp: serpStats,
          },
        };
      } catch (error: any) {
        console.error('Stats error:', error);
        return {
          error: true,
          message: 'Failed to get statistics',
        };
      }
    });
