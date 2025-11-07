/**
 * Auth Middleware
 * Validates JWT tokens and adds user context
 */

import type { Context } from 'elysia';
import type { JWTService } from '../utils/jwt';

export interface AuthContext {
  userId?: string;
  userEmail?: string;
}

/**
 * Required auth middleware - returns error if no valid token
 */
export const authMiddleware = (jwtService: JWTService) => {
  return async ({ headers, set }: Context) => {
    const authHeader = headers.authorization;

    if (!authHeader) {
      set.status = 401;
      return {
        error: true,
        message: 'Missing authorization header',
      };
    }

    const parts = authHeader.split(' ');
    if (parts.length !== 2 || parts[0] !== 'Bearer') {
      set.status = 401;
      return {
        error: true,
        message: 'Invalid authorization header format',
      };
    }

    const token = parts[1];

    try {
      const claims = await jwtService.verifyAccessToken(token);

      if (!claims) {
        set.status = 401;
        return {
          error: true,
          message: 'Invalid token',
        };
      }

      // Return user context to be used in the handler
      return {
        userId: claims.user_id,
        userEmail: claims.email,
      };
    } catch (error: any) {
      if (error.message === 'Token expired') {
        set.status = 401;
        return {
          error: true,
          message: 'Token expired',
        };
      }

      set.status = 401;
      return {
        error: true,
        message: 'Invalid token',
      };
    }
  };
};

/**
 * Optional auth middleware - validates JWT if present, but doesn't require it
 */
export const optionalAuthMiddleware = (jwtService: JWTService) => {
  return async ({ headers }: Context): Promise<AuthContext> => {
    const authHeader = headers.authorization;

    if (!authHeader) {
      return {};
    }

    const parts = authHeader.split(' ');
    if (parts.length !== 2 || parts[0] !== 'Bearer') {
      return {};
    }

    const token = parts[1];

    try {
      const claims = await jwtService.verifyAccessToken(token);

      if (claims) {
        return {
          userId: claims.user_id,
          userEmail: claims.email,
        };
      }
    } catch (error) {
      // Silently fail for optional auth
    }

    return {};
  };
};

/**
 * Decorator to extract auth context from headers
 */
export const extractAuth = (jwtService: JWTService) => {
  return async (headers: Record<string, string | undefined>): Promise<AuthContext> => {
    const authHeader = headers.authorization;

    if (!authHeader) {
      return {};
    }

    const parts = authHeader.split(' ');
    if (parts.length !== 2 || parts[0] !== 'Bearer') {
      return {};
    }

    const token = parts[1];

    try {
      const claims = await jwtService.verifyAccessToken(token);

      if (claims) {
        return {
          userId: claims.user_id,
          userEmail: claims.email,
        };
      }
    } catch (error) {
      // Silently fail
    }

    return {};
  };
};
