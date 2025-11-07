/**
 * Auth module - handles authentication endpoints
 */

import { Elysia, t } from 'elysia';
import type { Container } from '../../container';

export const authModule = (container: Container) =>
  new Elysia({ prefix: '/api/auth' })
    // Get Google OAuth URL
    .get('/google/url', ({ query }) => {
      const state = query.state || Math.random().toString(36).substring(7);
      const url = container.googleOAuthService.getAuthUrl(state);

      return { url, state };
    })

    // Google OAuth callback
    .get('/google/callback', async ({ query }) => {
      try {
        const { code } = query;

        if (!code) {
          return { error: true, message: 'No authorization code provided' };
        }

        // Exchange code for tokens
        const googleTokens = await container.googleOAuthService.exchangeCode(
          code
        );

        // Get user info
        const userInfo = await container.googleOAuthService.getUserInfo(
          googleTokens.access_token
        );

        // Create or update user
        const user = await container.authService.createOrUpdateUser(userInfo);

        // Generate JWT tokens
        const tokens = await container.authService.generateTokens(user);

        return {
          user,
          tokens,
        };
      } catch (error: any) {
        console.error('OAuth callback error:', error);
        return {
          error: true,
          message: 'Failed to authenticate with Google',
        };
      }
    })

    // Refresh tokens
    .post(
      '/refresh',
      async ({ body }) => {
        try {
          const tokens = await container.authService.refreshTokens(
            body.refresh_token
          );

          if (!tokens) {
            return { error: true, message: 'Invalid refresh token' };
          }

          return { tokens };
        } catch (error: any) {
          console.error('Token refresh error:', error);
          return {
            error: true,
            message: 'Failed to refresh tokens',
          };
        }
      },
      {
        body: t.Object({
          refresh_token: t.String(),
        }),
      }
    )

    // Verify token
    .get('/verify', async ({ headers }) => {
      try {
        const authHeader = headers.authorization;
        if (!authHeader || !authHeader.startsWith('Bearer ')) {
          return { error: true, message: 'No token provided' };
        }

        const token = authHeader.substring(7);
        const user = await container.authService.verifyAccessToken(token);

        if (!user) {
          return { error: true, message: 'Invalid token' };
        }

        return { user };
      } catch (error: any) {
        console.error('Token verification error:', error);
        return {
          error: true,
          message: 'Failed to verify token',
        };
      }
    })

    // Logout
    .post('/logout', async ({ headers }) => {
      try {
        const authHeader = headers.authorization;
        if (!authHeader || !authHeader.startsWith('Bearer ')) {
          return { error: true, message: 'No token provided' };
        }

        const token = authHeader.substring(7);
        const payload = await container.jwtService.verifyAccessToken(token);

        if (payload) {
          await container.authService.logout(payload.user_id);
        }

        return { success: true };
      } catch (error: any) {
        console.error('Logout error:', error);
        return {
          error: true,
          message: 'Failed to logout',
        };
      }
    });
