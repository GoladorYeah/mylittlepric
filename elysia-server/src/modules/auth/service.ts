/**
 * Auth Service - Business logic for authentication
 * Following Elysia Best Practices: abstract static classes
 */

import type { Redis } from 'ioredis';
import type { PrismaClient } from '@prisma/client';
import type { JWTService } from '../../utils/jwt';
import type { Config } from '../../config';
import type { User, GoogleUserInfo, AuthTokens } from '../../types';
import type { AuthModel } from './model';
import { randomUUID } from 'crypto';

// ═══════════════════════════════════════════════════════════
// GOOGLE OAUTH SERVICE
// ═══════════════════════════════════════════════════════════

/**
 * Google OAuth Service
 * Abstract class with static methods - no instance allocation needed
 */
export abstract class GoogleOAuth {
  /**
   * Get Google OAuth authorization URL
   */
  static getAuthUrl(config: Config, state: string): AuthModel.googleUrlResponse {
    const params = new URLSearchParams({
      client_id: config.googleClientId,
      redirect_uri: config.googleRedirectUrl,
      response_type: 'code',
      scope: 'openid email profile',
      state,
      access_type: 'offline',
      prompt: 'consent',
    });

    return {
      url: `https://accounts.google.com/o/oauth2/v2/auth?${params}`,
      state,
    };
  }

  /**
   * Exchange authorization code for tokens
   */
  static async exchangeCode(config: Config, code: string): Promise<any> {
    const response = await fetch('https://oauth2.googleapis.com/token', {
      method: 'POST',
      headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
      body: new URLSearchParams({
        code,
        client_id: config.googleClientId,
        client_secret: config.googleClientSecret,
        redirect_uri: config.googleRedirectUrl,
        grant_type: 'authorization_code',
      }),
    });

    if (!response.ok) {
      throw new Error('Failed to exchange code for tokens' satisfies AuthModel.authFailedError);
    }

    return response.json();
  }

  /**
   * Get user info from Google API
   */
  static async getUserInfo(accessToken: string): Promise<GoogleUserInfo> {
    const response = await fetch('https://www.googleapis.com/oauth2/v2/userinfo', {
      headers: { Authorization: `Bearer ${accessToken}` },
    });

    if (!response.ok) {
      throw new Error('Failed to get user info');
    }

    return response.json();
  }
}

// ═══════════════════════════════════════════════════════════
// AUTH SERVICE
// ═══════════════════════════════════════════════════════════

/**
 * Authentication Service
 * Abstract class with static methods - no instance allocation needed
 */
export abstract class Auth {
  /**
   * Create or update user from Google OAuth
   */
  static async createOrUpdateUser(
    prisma: PrismaClient,
    userInfo: GoogleUserInfo
  ): Promise<User> {
    const existingUser = await prisma.user.findFirst({
      where: {
        OR: [{ email: userInfo.email }, { provider_id: userInfo.id }],
      },
    });

    if (existingUser) {
      // Update existing user
      const updatedUser = await prisma.user.update({
        where: { id: existingUser.id },
        data: {
          name: userInfo.name,
          picture: userInfo.picture,
          updated_at: new Date(),
        },
      });

      return updatedUser;
    }

    // Create new user
    const newUser = await prisma.user.create({
      data: {
        id: randomUUID(),
        email: userInfo.email,
        name: userInfo.name,
        picture: userInfo.picture,
        provider: 'google',
        provider_id: userInfo.id,
        created_at: new Date(),
        updated_at: new Date(),
      },
    });

    return newUser;
  }

  /**
   * Generate JWT access and refresh tokens
   */
  static async generateTokens(
    jwtService: JWTService,
    redis: Redis,
    prisma: PrismaClient,
    user: User
  ): Promise<AuthTokens> {
    const accessToken = await jwtService.generateAccessToken({
      user_id: user.id,
      email: user.email,
    });

    const refreshToken = await jwtService.generateRefreshToken({
      user_id: user.id,
      email: user.email,
    });

    // Store refresh token in Redis
    const refreshTokenKey = `refresh_token:${user.id}`;
    await redis.set(refreshTokenKey, refreshToken, 'EX', 60 * 60 * 24 * 30); // 30 days

    // Store refresh token in database
    await prisma.refreshToken.create({
      data: {
        id: randomUUID(),
        user_id: user.id,
        token: refreshToken,
        expires_at: new Date(Date.now() + 60 * 60 * 24 * 30 * 1000), // 30 days
        created_at: new Date(),
      },
    });

    return {
      access_token: accessToken,
      refresh_token: refreshToken,
      expires_in: 3600, // 1 hour
    };
  }

  /**
   * Refresh access token using refresh token
   */
  static async refreshTokens(
    jwtService: JWTService,
    redis: Redis,
    prisma: PrismaClient,
    refreshToken: string
  ): Promise<AuthTokens | null> {
    // Verify refresh token
    const payload = await jwtService.verifyRefreshToken(refreshToken);

    if (!payload) {
      return null;
    }

    // Check if token exists in Redis
    const refreshTokenKey = `refresh_token:${payload.user_id}`;
    const storedToken = await redis.get(refreshTokenKey);

    if (!storedToken || storedToken !== refreshToken) {
      return null;
    }

    // Check if token exists in database
    const dbToken = await prisma.refreshToken.findFirst({
      where: {
        user_id: payload.user_id,
        token: refreshToken,
        expires_at: {
          gte: new Date(),
        },
      },
    });

    if (!dbToken) {
      return null;
    }

    // Get user
    const user = await prisma.user.findUnique({
      where: { id: payload.user_id },
    });

    if (!user) {
      return null;
    }

    // Generate new tokens
    return Auth.generateTokens(jwtService, redis, prisma, user);
  }

  /**
   * Verify access token and return user
   */
  static async verifyAccessToken(
    jwtService: JWTService,
    prisma: PrismaClient,
    token: string
  ): Promise<User | null> {
    const payload = await jwtService.verifyAccessToken(token);

    if (!payload) {
      return null;
    }

    const user = await prisma.user.findUnique({
      where: { id: payload.user_id },
    });

    return user;
  }

  /**
   * Logout user - invalidate refresh token
   */
  static async logout(
    redis: Redis,
    prisma: PrismaClient,
    userId: string
  ): Promise<void> {
    // Remove from Redis
    const refreshTokenKey = `refresh_token:${userId}`;
    await redis.del(refreshTokenKey);

    // Remove from database
    await prisma.refreshToken.deleteMany({
      where: { user_id: userId },
    });
  }
}
