/**
 * KeyRotator manages API key rotation using Redis
 * Tracks usage statistics and implements round-robin rotation
 */

import type { Redis } from 'ioredis';
import type { KeyRotatorStats } from '../types';

export class KeyRotator {
  private keys: string[];
  private serviceName: string;
  private redis: Redis;

  constructor(serviceName: string, keys: string[], redis: Redis) {
    this.keys = keys;
    this.serviceName = serviceName;
    this.redis = redis;
  }

  /**
   * Get the next API key in rotation
   */
  async getNextKey(): Promise<{ key: string; index: number }> {
    if (this.keys.length === 0) {
      throw new Error(`No API keys available for ${this.serviceName}`);
    }

    // Single key - no rotation needed
    if (this.keys.length === 1) {
      return { key: this.keys[0], index: 0 };
    }

    try {
      // Get current index from Redis
      const counterKey = `keyrotator:${this.serviceName}:counter`;

      // Increment and get the counter (atomic operation)
      const counter = await this.redis.incr(counterKey);

      // Calculate index using modulo
      const index = (counter - 1) % this.keys.length;

      return { key: this.keys[index], index };
    } catch (error) {
      // Fallback to first key if Redis fails
      console.error('Redis error in key rotation, using first key:', error);
      return { key: this.keys[0], index: 0 };
    }
  }

  /**
   * Get a specific key by index
   */
  getKeyByIndex(index: number): string {
    if (index < 0 || index >= this.keys.length) {
      throw new Error(`Invalid key index: ${index}`);
    }
    return this.keys[index];
  }

  /**
   * Record API key usage for analytics
   */
  async recordUsage(
    keyIndex: number,
    success: boolean,
    responseTime: number
  ): Promise<void> {
    const usageKey = `keyrotator:${this.serviceName}:usage:${keyIndex}`;

    try {
      const pipeline = this.redis.pipeline();

      // Increment usage counter
      pipeline.incr(usageKey);

      // Record success/failure
      if (success) {
        pipeline.incr(`${usageKey}:success`);
      } else {
        pipeline.incr(`${usageKey}:failures`);
      }

      // Record response time (milliseconds)
      pipeline.hincrby(`${usageKey}:response_times`, 'total', responseTime);
      pipeline.hincrby(`${usageKey}:response_times`, 'count', 1);

      await pipeline.exec();
    } catch (error) {
      console.error('Error recording key usage:', error);
    }
  }

  /**
   * Get usage statistics for a specific key
   */
  async getKeyStats(keyIndex: number): Promise<KeyRotatorStats> {
    const usageKey = `keyrotator:${this.serviceName}:usage:${keyIndex}`;

    try {
      const pipeline = this.redis.pipeline();
      pipeline.get(usageKey);
      pipeline.get(`${usageKey}:success`);
      pipeline.get(`${usageKey}:failures`);
      pipeline.hgetall(`${usageKey}:response_times`);

      const results = await pipeline.exec();

      if (!results) {
        return {
          key_index: keyIndex,
          total_usage: '0',
          success_count: '0',
          failure_count: '0',
        };
      }

      const totalUsage = (results[0][1] as string) || '0';
      const successCount = (results[1][1] as string) || '0';
      const failureCount = (results[2][1] as string) || '0';
      const responseTimes = (results[3][1] as Record<string, string>) || {};

      const stats: KeyRotatorStats = {
        key_index: keyIndex,
        total_usage: totalUsage,
        success_count: successCount,
        failure_count: failureCount,
      };

      // Calculate average response time
      if (responseTimes.total && responseTimes.count) {
        const total = parseInt(responseTimes.total, 10);
        const count = parseInt(responseTimes.count, 10);
        if (count > 0) {
          stats.avg_response_time_ms = Math.round(total / count);
        }
      }

      return stats;
    } catch (error) {
      console.error('Error getting key stats:', error);
      return {
        key_index: keyIndex,
        total_usage: '0',
        success_count: '0',
        failure_count: '0',
      };
    }
  }

  /**
   * Get statistics for all keys
   */
  async getAllStats(): Promise<KeyRotatorStats[]> {
    const stats: KeyRotatorStats[] = [];

    for (let i = 0; i < this.keys.length; i++) {
      const keyStats = await this.getKeyStats(i);
      stats.push(keyStats);
    }

    return stats;
  }

  /**
   * Reset the rotation counter (useful for testing)
   */
  async resetCounter(): Promise<void> {
    const counterKey = `keyrotator:${this.serviceName}:counter`;
    await this.redis.del(counterKey);
  }

  /**
   * Get the total number of available keys
   */
  getTotalKeys(): number {
    return this.keys.length;
  }
}
