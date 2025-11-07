/**
 * Database connection utilities for Redis and PostgreSQL (with Prisma)
 */

import Redis from 'ioredis';
import { PrismaClient } from '@prisma/client';
import type { Config } from '../config';

/**
 * Global Prisma client instance (singleton pattern for Bun)
 */
declare global {
  var prisma: PrismaClient | undefined;
}

/**
 * Initialize Prisma client with singleton pattern
 * This prevents multiple instances in development with hot reload
 */
export function initPrisma(): PrismaClient {
  if (global.prisma) {
    return global.prisma;
  }

  const prisma = new PrismaClient({
    log: process.env.NODE_ENV === 'development'
      ? ['query', 'error', 'warn']
      : ['error'],
  });

  // Store in global for hot reload
  if (process.env.NODE_ENV === 'development') {
    global.prisma = prisma;
  }

  console.log('✅ Connected to PostgreSQL via Prisma');
  return prisma;
}

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
 * Health check for PostgreSQL via Prisma
 */
export async function checkPrismaHealth(prisma: PrismaClient): Promise<boolean> {
  try {
    await prisma.$queryRaw`SELECT 1`;
    return true;
  } catch (error) {
    console.error('Prisma health check failed:', error);
    return false;
  }
}
