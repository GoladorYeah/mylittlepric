/**
 * SearchHistoryService manages search history in PostgreSQL
 */

import type { SearchHistoryItem } from '../types';
import { randomUUID } from 'crypto';

export class SearchHistoryService {
  private sql: any;

  constructor(sql: any) {
    this.sql = sql;
  }

  /**
   * Create search history entry
   */
  async createEntry(
    userId: string,
    sessionId: string,
    query: string,
    category: string,
    searchType: string,
    country: string,
    language: string,
    resultCount: number
  ): Promise<SearchHistoryItem> {
    try {
      const [entry] = await this.sql`
        INSERT INTO search_history (
          id, user_id, session_id, query, category, search_type,
          country, language, result_count, created_at
        )
        VALUES (
          ${randomUUID()}, ${userId}, ${sessionId}, ${query}, ${category},
          ${searchType}, ${country}, ${language}, ${resultCount}, NOW()
        )
        RETURNING *
      `;

      return entry;
    } catch (error) {
      console.error('Error creating search history entry:', error);
      throw error;
    }
  }

  /**
   * Get user search history
   */
  async getUserHistory(
    userId: string,
    limit: number = 20
  ): Promise<SearchHistoryItem[]> {
    try {
      const entries = await this.sql`
        SELECT * FROM search_history
        WHERE user_id = ${userId}
        ORDER BY created_at DESC
        LIMIT ${limit}
      `;

      return entries;
    } catch (error) {
      console.error('Error getting user history:', error);
      return [];
    }
  }

  /**
   * Get session search history
   */
  async getSessionHistory(sessionId: string): Promise<SearchHistoryItem[]> {
    try {
      const entries = await this.sql`
        SELECT * FROM search_history
        WHERE session_id = ${sessionId}
        ORDER BY created_at DESC
      `;

      return entries;
    } catch (error) {
      console.error('Error getting session history:', error);
      return [];
    }
  }
}
