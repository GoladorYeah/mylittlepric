/**
 * SearchHistoryService manages search history in PostgreSQL via Prisma
 */

import type { PrismaClient } from '@prisma/client';
import type { SearchHistory } from '@prisma/client';

export class SearchHistoryService {
  private prisma: PrismaClient;

  constructor(prisma: PrismaClient) {
    this.prisma = prisma;
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
  ): Promise<SearchHistory> {
    try {
      const entry = await this.prisma.searchHistory.create({
        data: {
          userId,
          sessionId,
          searchQuery: query,
          category,
          searchType,
          countryCode: country,
          languageCode: language,
          currency: 'CHF', // Default currency, could be passed as parameter
          resultCount,
        },
      });

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
  ): Promise<SearchHistory[]> {
    try {
      const entries = await this.prisma.searchHistory.findMany({
        where: {
          userId,
        },
        orderBy: {
          createdAt: 'desc',
        },
        take: limit,
      });

      return entries;
    } catch (error) {
      console.error('Error getting user history:', error);
      return [];
    }
  }

  /**
   * Get session search history
   */
  async getSessionHistory(sessionId: string): Promise<SearchHistory[]> {
    try {
      const entries = await this.prisma.searchHistory.findMany({
        where: {
          sessionId,
        },
        orderBy: {
          createdAt: 'desc',
        },
      });

      return entries;
    } catch (error) {
      console.error('Error getting session history:', error);
      return [];
    }
  }

  /**
   * Get search history with user and session details
   */
  async getHistoryWithDetails(
    userId?: string,
    sessionId?: string,
    limit: number = 20
  ): Promise<SearchHistory[]> {
    try {
      const where: any = {};
      if (userId) where.userId = userId;
      if (sessionId) where.sessionId = sessionId;

      const entries = await this.prisma.searchHistory.findMany({
        where,
        include: {
          user: {
            select: {
              id: true,
              email: true,
              fullName: true,
            },
          },
          session: {
            select: {
              id: true,
              sessionId: true,
              countryCode: true,
              languageCode: true,
            },
          },
        },
        orderBy: {
          createdAt: 'desc',
        },
        take: limit,
      });

      return entries as any;
    } catch (error) {
      console.error('Error getting search history with details:', error);
      return [];
    }
  }

  /**
   * Update clicked product for analytics
   */
  async updateClickedProduct(
    historyId: string,
    productId: string
  ): Promise<void> {
    try {
      await this.prisma.searchHistory.update({
        where: { id: historyId },
        data: { clickedProductId: productId },
      });
    } catch (error) {
      console.error('Error updating clicked product:', error);
    }
  }

  /**
   * Get search statistics for a user
   */
  async getUserSearchStats(userId: string): Promise<{
    totalSearches: number;
    searchesByType: Record<string, number>;
    searchesByCategory: Record<string, number>;
  }> {
    try {
      const searches = await this.prisma.searchHistory.findMany({
        where: { userId },
        select: {
          searchType: true,
          category: true,
        },
      });

      const searchesByType: Record<string, number> = {};
      const searchesByCategory: Record<string, number> = {};

      searches.forEach((search) => {
        searchesByType[search.searchType] = (searchesByType[search.searchType] || 0) + 1;
        if (search.category) {
          searchesByCategory[search.category] = (searchesByCategory[search.category] || 0) + 1;
        }
      });

      return {
        totalSearches: searches.length,
        searchesByType,
        searchesByCategory,
      };
    } catch (error) {
      console.error('Error getting user search stats:', error);
      return {
        totalSearches: 0,
        searchesByType: {},
        searchesByCategory: {},
      };
    }
  }

  /**
   * Clean up expired anonymous search history
   */
  async cleanupExpiredHistory(): Promise<number> {
    try {
      const result = await this.prisma.searchHistory.deleteMany({
        where: {
          expiresAt: {
            lt: new Date(),
          },
          userId: null, // Only delete anonymous searches
        },
      });

      return result.count;
    } catch (error) {
      console.error('Error cleaning up expired history:', error);
      return 0;
    }
  }
}
