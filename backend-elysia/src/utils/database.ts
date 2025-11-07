/**
 * Database connection utilities for Redis and PostgreSQL
 */

import Redis from 'ioredis';
import postgres from 'postgres';
import type { Config } from '../config';

/**
 * Initialize Redis connection
 */
export function initRedis(config: Config): Redis {
  const [host, port] = config.redisUrl.split(':');

  const redis = new Redis({
    host,
    port: parseInt(port, 10),
    password: config.redisPassword || undefined,
    db: config.redisDb,
    retryStrategy(times) {
      const delay = Math.min(times * 50, 2000);
      return delay;
    },
    maxRetriesPerRequest: 3,
  });

  redis.on('connect', () => {
    console.log('✅ Connected to Redis');
  });

  redis.on('error', (err) => {
    console.error('❌ Redis connection error:', err);
  });

  return redis;
}

/**
 * Initialize PostgreSQL connection
 */
export function initPostgres(config: Config) {
  const sql = postgres(config.databaseUrl, {
    max: 25,
    idle_timeout: 20,
    connect_timeout: 10,
    onnotice: () => {}, // Suppress notices
  });

  console.log('✅ Connected to PostgreSQL');
  return sql;
}

/**
 * Health check for Redis
 */
export async function checkRedisHealth(redis: Redis): Promise<boolean> {
  try {
    const result = await redis.ping();
    return result === 'PONG';
  } catch (error) {
    console.error('Redis health check failed:', error);
    return false;
  }
}

/**
 * Health check for PostgreSQL
 */
export async function checkPostgresHealth(sql: any): Promise<boolean> {
  try {
    await sql`SELECT 1`;
    return true;
  } catch (error) {
    console.error('PostgreSQL health check failed:', error);
    return false;
  }
}
