/**
 * SessionService manages chat sessions with Redis persistence
 */

import type { Redis } from 'ioredis';
import type { ChatSession, Message, SearchState } from '../types';
import { randomUUID } from 'crypto';

export class SessionService {
  private redis: Redis;
  private sessionTTL: number;
  private maxMessages: number;

  constructor(redis: Redis, sessionTTL: number, maxMessages: number) {
    this.redis = redis;
    this.sessionTTL = sessionTTL;
    this.maxMessages = maxMessages;
  }

  /**
   * Get or create session
   */
  async getOrCreateSession(
    sessionId: string,
    country: string,
    language: string,
    currency: string
  ): Promise<ChatSession> {
    const key = `session:${sessionId}`;

    try {
      const data = await this.redis.get(key);

      if (data) {
        const session = JSON.parse(data);
        // Update TTL
        await this.redis.expire(key, this.sessionTTL);
        return session;
      }

      // Create new session
      const newSession: ChatSession = {
        id: randomUUID(),
        session_id: sessionId,
        country_code: country,
        language_code: language,
        currency: currency,
        message_count: 0,
        search_state: {
          status: 'idle',
          category: '',
          search_count: 0,
        },
        cycle_state: {
          cycle_id: 1,
          iteration: 1,
          cycle_history: [],
          last_defined: [],
          prompt_id: 'v1',
          prompt_hash: '',
        },
        created_at: new Date(),
        updated_at: new Date(),
        expires_at: new Date(Date.now() + this.sessionTTL * 1000),
      };

      await this.redis.setex(key, this.sessionTTL, JSON.stringify(newSession));

      return newSession;
    } catch (error) {
      console.error('Error getting/creating session:', error);
      throw error;
    }
  }

  /**
   * Update session
   */
  async updateSession(session: ChatSession): Promise<void> {
    const key = `session:${session.session_id}`;
    session.updated_at = new Date();

    try {
      await this.redis.setex(key, this.sessionTTL, JSON.stringify(session));
    } catch (error) {
      console.error('Error updating session:', error);
      throw error;
    }
  }

  /**
   * Get conversation history
   */
  async getMessages(sessionId: string): Promise<Message[]> {
    const key = `messages:${sessionId}`;

    try {
      const data = await this.redis.get(key);
      return data ? JSON.parse(data) : [];
    } catch (error) {
      console.error('Error getting messages:', error);
      return [];
    }
  }

  /**
   * Add message to conversation history
   */
  async addMessage(sessionId: string, message: Message): Promise<void> {
    const key = `messages:${sessionId}`;

    try {
      const messages = await this.getMessages(sessionId);
      messages.push(message);

      // Keep only last N messages
      const trimmed = messages.slice(-this.maxMessages);

      await this.redis.setex(key, this.sessionTTL, JSON.stringify(trimmed));
    } catch (error) {
      console.error('Error adding message:', error);
      throw error;
    }
  }

  /**
   * Update search state
   */
  async updateSearchState(
    sessionId: string,
    searchState: Partial<SearchState>
  ): Promise<void> {
    const session = await this.getOrCreateSession(sessionId, '', '', '');
    session.search_state = { ...session.search_state, ...searchState };
    await this.updateSession(session);
  }

  /**
   * Clear session
   */
  async clearSession(sessionId: string): Promise<void> {
    try {
      await this.redis.del(`session:${sessionId}`, `messages:${sessionId}`);
    } catch (error) {
      console.error('Error clearing session:', error);
      throw error;
    }
  }
}
