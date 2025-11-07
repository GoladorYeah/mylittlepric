/**
 * JWT utility for token generation and validation
 */

import { sign, verify } from '@elysiajs/jwt';

export interface JWTPayload {
  user_id: string;
  email: string;
  type: 'access' | 'refresh';
  iat?: number;
  exp?: number;
}

export class JWTService {
  private accessSecret: string;
  private refreshSecret: string;
  private accessTTL: number; // seconds
  private refreshTTL: number; // seconds

  constructor(
    accessSecret: string,
    refreshSecret: string,
    accessTTL: number,
    refreshTTL: number
  ) {
    this.accessSecret = accessSecret;
    this.refreshSecret = refreshSecret;
    this.accessTTL = accessTTL;
    this.refreshTTL = refreshTTL;
  }

  /**
   * Generate an access token
   */
  async generateAccessToken(userId: string, email: string): Promise<string> {
    const payload: JWTPayload = {
      user_id: userId,
      email,
      type: 'access',
      iat: Math.floor(Date.now() / 1000),
      exp: Math.floor(Date.now() / 1000) + this.accessTTL,
    };

    return await this.signToken(payload, this.accessSecret);
  }

  /**
   * Generate a refresh token
   */
  async generateRefreshToken(userId: string, email: string): Promise<string> {
    const payload: JWTPayload = {
      user_id: userId,
      email,
      type: 'refresh',
      iat: Math.floor(Date.now() / 1000),
      exp: Math.floor(Date.now() / 1000) + this.refreshTTL,
    };

    return await this.signToken(payload, this.refreshSecret);
  }

  /**
   * Verify an access token
   */
  async verifyAccessToken(token: string): Promise<JWTPayload | null> {
    return await this.verifyToken(token, this.accessSecret, 'access');
  }

  /**
   * Verify a refresh token
   */
  async verifyRefreshToken(token: string): Promise<JWTPayload | null> {
    return await this.verifyToken(token, this.refreshSecret, 'refresh');
  }

  /**
   * Sign a JWT token
   */
  private async signToken(
    payload: JWTPayload,
    secret: string
  ): Promise<string> {
    // Using native JWT signing with Bun
    const encoder = new TextEncoder();
    const data = encoder.encode(JSON.stringify(payload));
    const key = encoder.encode(secret);

    // For simplicity, using base64 encoding
    // In production, consider using proper JWT library
    const header = { alg: 'HS256', typ: 'JWT' };
    const encodedHeader = btoa(JSON.stringify(header));
    const encodedPayload = btoa(JSON.stringify(payload));

    const signature = await this.createSignature(
      `${encodedHeader}.${encodedPayload}`,
      secret
    );

    return `${encodedHeader}.${encodedPayload}.${signature}`;
  }

  /**
   * Verify a JWT token
   */
  private async verifyToken(
    token: string,
    secret: string,
    expectedType: 'access' | 'refresh'
  ): Promise<JWTPayload | null> {
    try {
      const parts = token.split('.');
      if (parts.length !== 3) {
        return null;
      }

      const [encodedHeader, encodedPayload, signature] = parts;

      // Verify signature
      const expectedSignature = await this.createSignature(
        `${encodedHeader}.${encodedPayload}`,
        secret
      );

      if (signature !== expectedSignature) {
        return null;
      }

      // Decode payload
      const payload: JWTPayload = JSON.parse(atob(encodedPayload));

      // Check expiration
      if (payload.exp && payload.exp < Math.floor(Date.now() / 1000)) {
        return null;
      }

      // Check type
      if (payload.type !== expectedType) {
        return null;
      }

      return payload;
    } catch (error) {
      console.error('JWT verification error:', error);
      return null;
    }
  }

  /**
   * Create HMAC signature
   */
  private async createSignature(
    data: string,
    secret: string
  ): Promise<string> {
    const encoder = new TextEncoder();
    const keyData = encoder.encode(secret);
    const messageData = encoder.encode(data);

    const key = await crypto.subtle.importKey(
      'raw',
      keyData,
      { name: 'HMAC', hash: 'SHA-256' },
      false,
      ['sign']
    );

    const signature = await crypto.subtle.sign('HMAC', key, messageData);

    // Convert to base64url
    return btoa(String.fromCharCode(...new Uint8Array(signature)))
      .replace(/\+/g, '-')
      .replace(/\//g, '_')
      .replace(/=/g, '');
  }

  /**
   * Get access token TTL in seconds
   */
  getAccessTTL(): number {
    return this.accessTTL;
  }

  /**
   * Get refresh token TTL in seconds
   */
  getRefreshTTL(): number {
    return this.refreshTTL;
  }
}
