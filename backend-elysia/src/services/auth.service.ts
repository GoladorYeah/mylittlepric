/**
 * AuthService handles user authentication via Google OAuth
 */

import type { Redis } from 'ioredis';
import type { JWTService } from '../utils/jwt';
import type { Config } from '../config';
import type { User, GoogleUserInfo, AuthTokens } from '../types';
import { randomUUID } from 'crypto';

export class GoogleOAuthService {
  private config: Config;

  constructor(config: Config) {
    this.config = config;
  }

  /**
   * Get Google OAuth URL
   */
  getAuthUrl(state: string): string {
    const params = new URLSearchParams({
      client_id: this.config.googleClientId,
      redirect_uri: this.config.googleRedirectUrl,
      response_type: 'code',
      scope: 'openid email profile',
      state,
      access_type: 'offline',
      prompt: 'consent',
    });

    return `https://accounts.google.com/o/oauth2/v2/auth?${params}`;
  }

  /**
   * Exchange code for tokens
   */
  async exchangeCode(code: string): Promise<any> {
    const response = await fetch('https://oauth2.googleapis.com/token', {
      method: 'POST',
      headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
      body: new URLSearchParams({
        code,
        client_id: this.config.googleClientId,
        client_secret: this.config.googleClientSecret,
        redirect_uri: this.config.googleRedirectUrl,
        grant_type: 'authorization_code',
      }),
    });

    if (!response.ok) {
      throw new Error('Failed to exchange code for tokens');
    }

    return response.json();
  }

  /**
   * Get user info from Google
   */
  async getUserInfo(accessToken: string): Promise<GoogleUserInfo> {
    const response = await fetch(
      'https://www.googleapis.com/oauth2/v2/userinfo',
      {
        headers: { Authorization: `Bearer ${accessToken}` },
      }
    );

    if (!response.ok) {
      throw new Error('Failed to get user info');
    }

    return response.json();
  }
}

export class AuthService {
  private redis: Redis;
  private jwtService: JWTService;
  private googleOAuth: GoogleOAuthService;

  constructor(
    redis: Redis,
    jwtService: JWTService,
    googleOAuth: GoogleOAuthService
  ) {
    this.redis = redis;
    this.jwtService = jwtService;
    this.googleOAuth = googleOAuth;
  }

  /**
   * Create or update user
   */
  async createOrUpdateUser(userInfo: GoogleUserInfo): Promise<User> {
    const userId = `google:${userInfo.id}`;
    const key = `user:${userId}`;

    const user: User = {
      id: userId,
      email: userInfo.email,
      name: userInfo.name,
      picture: userInfo.picture,
      provider: 'google',
      provider_id: userInfo.id,
      created_at: new Date(),
      updated_at: new Date(),
    };

    // Save user to Redis
    await this.redis.setex(key, 86400 * 30, JSON.stringify(user)); // 30 days

    return user;
  }

  /**
   * Generate authentication tokens
   */
  async generateTokens(user: User): Promise<AuthTokens> {
    const accessToken = await this.jwtService.generateAccessToken(
      user.id,
      user.email
    );
    const refreshToken = await this.jwtService.generateRefreshToken(
      user.id,
      user.email
    );

    // Store refresh token in Redis
    const refreshKey = `refresh_token:${user.id}`;
    await this.redis.setex(
      refreshKey,
      this.jwtService.getRefreshTTL(),
      refreshToken
    );

    return {
      access_token: accessToken,
      refresh_token: refreshToken,
      expires_in: this.jwtService.getAccessTTL(),
    };
  }

  /**
   * Verify access token
   */
  async verifyAccessToken(token: string): Promise<User | null> {
    const payload = await this.jwtService.verifyAccessToken(token);
    if (!payload) return null;

    const key = `user:${payload.user_id}`;
    const data = await this.redis.get(key);

    return data ? JSON.parse(data) : null;
  }

  /**
   * Refresh tokens
   */
  async refreshTokens(refreshToken: string): Promise<AuthTokens | null> {
    const payload = await this.jwtService.verifyRefreshToken(refreshToken);
    if (!payload) return null;

    // Verify refresh token exists in Redis
    const refreshKey = `refresh_token:${payload.user_id}`;
    const storedToken = await this.redis.get(refreshKey);

    if (storedToken !== refreshToken) {
      return null;
    }

    // Get user
    const userKey = `user:${payload.user_id}`;
    const userData = await this.redis.get(userKey);
    if (!userData) return null;

    const user: User = JSON.parse(userData);

    // Generate new tokens
    return this.generateTokens(user);
  }

  /**
   * Logout user
   */
  async logout(userId: string): Promise<void> {
    const refreshKey = `refresh_token:${userId}`;
    await this.redis.del(refreshKey);
  }
}
