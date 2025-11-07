/**
 * Chat Module - Controller (HTTP routing & validation)
 * Following Elysia Best Practices: Elysia instance = controller
 */

import { Elysia } from 'elysia';
import type { Container } from '../../container';
import { ChatModel } from './model';
import { Chat } from './service';

/**
 * Chat Module
 * Handles chat messages, product search, and product details
 */
export const chatModule = (container: Container) =>
  new Elysia({ prefix: '/api' })
    // ═══════════════════════════════════════════════════════════
    // POST /api/chat
    // Process chat message and return response
    // ═══════════════════════════════════════════════════════════
    .post(
      '/chat',
      async ({ body }) => {
        return Chat.processMessage(
          body.session_id,
          body.message,
          body.country,
          body.language,
          body.currency,
          body.current_category,
          container.sessionService,
          container.geminiService,
          container.serpService,
          container.config
        );
      },
      {
        body: ChatModel.chatRequest,
        response: {
          200: ChatModel.chatResponse,
        },
      }
    )

    // ═══════════════════════════════════════════════════════════
    // GET /api/chat/messages
    // Get conversation history for a session
    // ═══════════════════════════════════════════════════════════
    .get(
      '/chat/messages',
      async ({ query, set }) => {
        try {
          const { session_id } = query;

          if (!session_id) {
            set.status = 400;
            return {
              error: 'invalid_request' as const,
              message: 'session_id is required',
            };
          }

          // Get messages from session
          const messages = await container.sessionService.getMessages(
            session_id
          );

          return {
            session_id,
            messages,
            count: messages.length,
          };
        } catch (error: any) {
          console.error('Get messages error:', error);
          set.status = 500;
          return {
            error: 'server_error' as const,
            message: 'Failed to retrieve messages',
          };
        }
      },
      {
        query: ChatModel.getMessagesQuery,
        response: {
          200: ChatModel.getMessagesResponse,
          400: ChatModel.invalidRequestError,
          500: ChatModel.serverError,
        },
      }
    )

    // ═══════════════════════════════════════════════════════════
    // POST /api/product-details
    // Get detailed product information
    // ═══════════════════════════════════════════════════════════
    .post(
      '/product-details',
      async ({ body, set }) => {
        try {
          const details = await Chat.getProductDetails(
            body.page_token,
            body.country,
            container.serpService,
            container.config
          );

          return details;
        } catch (error: any) {
          console.error('Product details error:', error);
          set.status = 500;
          return {
            error: true,
            message: 'Failed to get product details',
          };
        }
      },
      {
        body: ChatModel.productDetailsRequest,
        response: {
          200: ChatModel.productDetailsResponse,
        },
      }
    );
