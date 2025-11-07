/**
 * Auth Module - Controller (HTTP routing & validation)
 * Following Elysia Best Practices: Elysia instance = controller
 */

import { Elysia } from 'elysia';
import type { Container } from '../../container';
import { AuthModel } from './model';
import { Auth, GoogleOAuth } from './service';

/**
 * Auth Module
 * Handles Google OAuth authentication, token refresh, and logout
 */
export const authModule = (container: Container) =>
  new Elysia({ prefix: '/api/auth' })
    // ═══════════════════════════════════════════════════════════
    // GET /api/auth/google/url
    // Get Google OAuth authorization URL
    // ═══════════════════════════════════════════════════════════
    .get(
      '/google/url',
      ({ query }) => {
        const state = query.state || Math.random().toString(36).substring(7);
        return GoogleOAuth.getAuthUrl(container.config, state);
      },
      {
        query: AuthModel.googleUrlQuery,
        response: {
          200: AuthModel.googleUrlResponse,
        },
      }
    )

    // ═══════════════════════════════════════════════════════════
    // GET /api/auth/google/callback
    // Handle Google OAuth callback
    // ═══════════════════════════════════════════════════════════
    .get(
      '/google/callback',
      async ({ query, set }) => {
        try {
          const { code } = query;

          if (!code) {
            set.status = 400;
            return {
              error: true,
              message: 'No authorization code provided' satisfies AuthModel.noCodeError,
            };
          }

          // Exchange code for tokens
          const googleTokens = await GoogleOAuth.exchangeCode(
            container.config,
            code
          );

          // Get user info from Google
          const userInfo = await GoogleOAuth.getUserInfo(
            googleTokens.access_token
          );

          // Create or update user in database
          const user = await Auth.createOrUpdateUser(
            container.prisma,
            userInfo
          );

          // Generate JWT tokens
          const tokens = await Auth.generateTokens(
            container.jwtService,
            container.redis,
            container.prisma,
            user
          );

          return {
            user,
            tokens,
          };
        } catch (error: any) {
          console.error('OAuth callback error:', error);
          set.status = 500;
          return {
            error: true,
            message: 'Failed to authenticate with Google' satisfies AuthModel.authFailedError,
          };
        }
      },
      {
        query: AuthModel.googleCallbackQuery,
        response: {
          200: AuthModel.authResponse,
          400: AuthModel.errorResponse,
          500: AuthModel.errorResponse,
        },
      }
    )

    // ═══════════════════════════════════════════════════════════
    // POST /api/auth/refresh
    // Refresh access token using refresh token
    // ═══════════════════════════════════════════════════════════
    .post(
      '/refresh',
      async ({ body, set }) => {
        try {
          const tokens = await Auth.refreshTokens(
            container.jwtService,
            container.redis,
            container.prisma,
            body.refresh_token
          );

          if (!tokens) {
            set.status = 401;
            return {
              error: true,
              message: 'Invalid refresh token' satisfies AuthModel.invalidTokenError,
            };
          }

          return { tokens };
        } catch (error: any) {
          console.error('Token refresh error:', error);
          set.status = 500;
          return {
            error: true,
            message: 'Failed to refresh tokens',
          };
        }
      },
      {
        body: AuthModel.refreshTokenBody,
        response: {
          200: AuthModel.refreshTokenResponse,
          401: AuthModel.errorResponse,
          500: AuthModel.errorResponse,
        },
      }
    )

    // ═══════════════════════════════════════════════════════════
    // GET /api/auth/verify
    // Verify access token and return user info
    // ═══════════════════════════════════════════════════════════
    .get('/verify', async ({ headers, set }) => {
      try {
        const authHeader = headers.authorization;

        if (!authHeader || !authHeader.startsWith('Bearer ')) {
          set.status = 401;
          return {
            error: true,
            message: 'No token provided' satisfies AuthModel.noTokenError,
          };
        }

        const token = authHeader.substring(7);
        const user = await Auth.verifyAccessToken(
          container.jwtService,
          container.prisma,
          token
        );

        if (!user) {
          set.status = 401;
          return {
            error: true,
            message: 'Invalid token' satisfies AuthModel.invalidTokenError,
          };
        }

        return { user };
      } catch (error: any) {
        console.error('Token verification error:', error);
        set.status = 500;
        return {
          error: true,
          message: 'Failed to verify token',
        };
      }
    })

    // ═══════════════════════════════════════════════════════════
    // POST /api/auth/logout
    // Logout user and invalidate refresh token
    // ═══════════════════════════════════════════════════════════
    .post('/logout', async ({ headers, set }) => {
      try {
        const authHeader = headers.authorization;

        if (!authHeader || !authHeader.startsWith('Bearer ')) {
          set.status = 401;
          return {
            error: true,
            message: 'No token provided' satisfies AuthModel.noTokenError,
          };
        }

        const token = authHeader.substring(7);
        const payload = await container.jwtService.verifyAccessToken(token);

        if (payload) {
          await Auth.logout(container.redis, container.prisma, payload.user_id);
        }

        return { success: true };
      } catch (error: any) {
        console.error('Logout error:', error);
        set.status = 500;
        return {
          error: true,
          message: 'Failed to logout',
        };
      }
    });
