/**
 * Auth Model - TypeBox schemas and type definitions
 * Following Elysia Best Practices: namespace for DTOs
 */

import { t } from 'elysia';

export namespace AuthModel {
  // ═══════════════════════════════════════════════════════════
  // GOOGLE OAUTH
  // ═══════════════════════════════════════════════════════════

  // GET /api/auth/google/url - Query params
  export const googleUrlQuery = t.Object({
    state: t.Optional(t.String()),
  });

  export type googleUrlQuery = typeof googleUrlQuery.static;

  // GET /api/auth/google/url - Response
  export const googleUrlResponse = t.Object({
    url: t.String(),
    state: t.String(),
  });

  export type googleUrlResponse = typeof googleUrlResponse.static;

  // GET /api/auth/google/callback - Query params
  export const googleCallbackQuery = t.Object({
    code: t.String(),
    state: t.Optional(t.String()),
  });

  export type googleCallbackQuery = typeof googleCallbackQuery.static;

  // Google OAuth callback response
  export const authResponse = t.Object({
    user: t.Object({
      id: t.String(),
      email: t.String(),
      name: t.String(),
      picture: t.Optional(t.String()),
      provider: t.String(),
      provider_id: t.String(),
    }),
    tokens: t.Object({
      access_token: t.String(),
      refresh_token: t.String(),
      expires_in: t.Number(),
    }),
  });

  export type authResponse = typeof authResponse.static;

  // ═══════════════════════════════════════════════════════════
  // TOKEN REFRESH
  // ═══════════════════════════════════════════════════════════

  // POST /api/auth/refresh - Request body
  export const refreshTokenBody = t.Object({
    refresh_token: t.String(),
  });

  export type refreshTokenBody = typeof refreshTokenBody.static;

  // POST /api/auth/refresh - Response
  export const refreshTokenResponse = t.Object({
    tokens: t.Object({
      access_token: t.String(),
      refresh_token: t.String(),
      expires_in: t.Number(),
    }),
  });

  export type refreshTokenResponse = typeof refreshTokenResponse.static;

  // ═══════════════════════════════════════════════════════════
  // TOKEN VERIFICATION
  // ═══════════════════════════════════════════════════════════

  // GET /api/auth/verify - Response
  export const verifyTokenResponse = t.Object({
    user: t.Object({
      id: t.String(),
      email: t.String(),
      name: t.String(),
      picture: t.Optional(t.String()),
    }),
  });

  export type verifyTokenResponse = typeof verifyTokenResponse.static;

  // ═══════════════════════════════════════════════════════════
  // ERROR RESPONSES
  // ═══════════════════════════════════════════════════════════

  export const errorResponse = t.Object({
    error: t.Boolean(),
    message: t.String(),
  });

  export type errorResponse = typeof errorResponse.static;

  // Specific error messages
  export const noCodeError = t.Literal('No authorization code provided');
  export type noCodeError = typeof noCodeError.static;

  export const noTokenError = t.Literal('No token provided');
  export type noTokenError = typeof noTokenError.static;

  export const invalidTokenError = t.Literal('Invalid token');
  export type invalidTokenError = typeof invalidTokenError.static;

  export const authFailedError = t.Literal('Failed to authenticate with Google');
  export type authFailedError = typeof authFailedError.static;
}
