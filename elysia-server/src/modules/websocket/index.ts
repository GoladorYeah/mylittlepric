/**
 * WebSocket Module - Controller (WebSocket routing & validation)
 * Following Elysia Best Practices: Elysia instance = controller
 */

import { Elysia } from 'elysia';
import { ws } from '@elysiajs/ws';
import type { Container } from '../../container';
import { WebSocketModel } from './model';
import { Chat } from '../chat/service';

/**
 * WebSocket Module
 * Handles real-time chat communication via WebSocket
 * Uses Chat service to avoid code duplication with REST endpoints
 */
export const websocketModule = (container: Container) =>
  new Elysia()
    .use(
      ws({
        path: '/ws',
      })
    )
    .ws('/ws', {
      message: async (ws, msg: any) => {
        try {
          console.log(`📨 WebSocket message type: ${msg.type}`);

          // ═══════════════════════════════════════════════════════════
          // CHAT MESSAGE
          // ═══════════════════════════════════════════════════════════
          if (msg.type === 'chat') {
            if (!msg.message) {
              ws.send({
                type: 'error',
                error: 'INVALID_REQUEST',
                message: 'Message is required',
                session_id: msg.session_id,
              });
              return;
            }

            // Use Chat service to process message (DRY - Don't Repeat Yourself)
            const response = await Chat.processMessage(
              msg.session_id,
              msg.message,
              msg.country || container.config.defaultCountry,
              msg.language || container.config.defaultLanguage,
              msg.currency || container.config.defaultCurrency,
              msg.current_category || '',
              container.sessionService,
              container.geminiService,
              container.serpService,
              container.config
            );

            ws.send(response);
            return;
          }

          // ═══════════════════════════════════════════════════════════
          // PRODUCT DETAILS REQUEST
          // ═══════════════════════════════════════════════════════════
          if (msg.type === 'product-details') {
            if (!msg.page_token) {
              ws.send({
                type: 'error',
                error: 'INVALID_REQUEST',
                message: 'page_token is required',
                session_id: msg.session_id,
              });
              return;
            }

            const details = await Chat.getProductDetails(
              msg.page_token,
              msg.country || container.config.defaultCountry,
              container.serpService,
              container.config
            );

            ws.send({
              type: 'product_details',
              product_details: details,
              session_id: msg.session_id,
            });
            return;
          }

          // ═══════════════════════════════════════════════════════════
          // UNKNOWN MESSAGE TYPE
          // ═══════════════════════════════════════════════════════════
          ws.send({
            type: 'error',
            error: 'INVALID_MESSAGE_TYPE',
            message: 'Invalid message type',
            session_id: msg.session_id || '',
          });
        } catch (error: any) {
          console.error('WebSocket error:', error);
          ws.send({
            type: 'error',
            error: 'SERVER_ERROR',
            message: error.message || 'An error occurred',
            session_id: msg.session_id || '',
          });
        }
      },
      open: (ws) => {
        console.log('🔌 WebSocket client connected');
      },
      close: (ws) => {
        console.log('🔌 WebSocket client disconnected');
      },
      error: (ws, error) => {
        console.error('❌ WebSocket error:', error);
      },
    });
