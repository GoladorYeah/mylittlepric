/**
 * WebSocket Module - Controller (WebSocket routing & validation)
 * Following Elysia Best Practices: Elysia instance = controller
 */

import { Elysia } from 'elysia';
import type { Container } from '../../container';
import { WebSocketModel } from './model';
import { Chat } from '../chat/service';

/**
 * WebSocket Module
 * Handles real-time chat communication via WebSocket
 * Uses Chat service to avoid code duplication with REST endpoints
 *
 * Following Elysia WebSocket Best Practices:
 * - Schema validation for body, query, and response
 * - Access validated data via ws.data
 * - Proper WebSocket configuration
 * - Type-safe message handling
 */
export const websocketModule = (container: Container) =>
  new Elysia()
    .ws('/ws', {
      // ═══════════════════════════════════════════════════════════
      // SCHEMA VALIDATION
      // ═══════════════════════════════════════════════════════════

      // Validate incoming message body
      body: WebSocketModel.body,

      // Validate query parameters (e.g., /ws?token=xxx)
      query: WebSocketModel.query,

      // Validate outgoing responses
      response: WebSocketModel.outgoingResponse,

      // ═══════════════════════════════════════════════════════════
      // WEBSOCKET CONFIGURATION
      // ═══════════════════════════════════════════════════════════

      // Enable compression for clients that support it
      perMessageDeflate: true,

      // Maximum message size (16MB)
      maxPayloadLength: 16 * 1024 * 1024,

      // Auto-close connection after 120 seconds of inactivity
      idleTimeout: 120,

      // Maximum backpressure buffer (16MB)
      backpressureLimit: 16 * 1024 * 1024,

      // Don't close connection if backpressure limit is reached (just slow down)
      closeOnBackpressureLimit: false,

      // ═══════════════════════════════════════════════════════════
      // LIFECYCLE HOOKS
      // ═══════════════════════════════════════════════════════════

      /**
       * Called when WebSocket connection is opened
       */
      open: (ws) => {
        const { token } = ws.data.query;
        console.log('🔌 WebSocket client connected', token ? `(authenticated)` : '(anonymous)');
      },

      /**
       * Called when WebSocket connection is closed
       */
      close: (ws, code, reason) => {
        console.log(`🔌 WebSocket client disconnected (code: ${code}, reason: ${reason})`);
      },

      /**
       * Called when WebSocket error occurs
       */
      error: (ws, error) => {
        console.error('❌ WebSocket error:', error);
      },

      /**
       * Transform and validate incoming messages before processing
       */
      transformMessage: (ws, message) => {
        // Parse JSON string messages
        if (typeof message === 'string') {
          try {
            return JSON.parse(message);
          } catch (error) {
            console.error('Failed to parse WebSocket message:', error);
            // Send error response and return null to skip message handling
            ws.send({
              type: 'error',
              error: 'INVALID_JSON',
              message: 'Invalid JSON format',
              session_id: '',
            });
            return null;
          }
        }
        return message;
      },

      // ═══════════════════════════════════════════════════════════
      // MESSAGE HANDLER
      // ═══════════════════════════════════════════════════════════

      /**
       * Handle incoming WebSocket messages
       * Message is automatically validated against body schema
       */
      message: async (ws, message) => {
        // Skip if message is null (from transformMessage error)
        if (message === null) {
          return;
        }

        try {
          // Type guard to ensure message has the expected structure
          if (typeof message !== 'object' || !message || !('type' in message)) {
            ws.send({
              type: 'error',
              error: 'INVALID_MESSAGE',
              message: 'Invalid message structure',
              session_id: '',
            });
            return;
          }

          const msg = message as WebSocketModel.body;
          console.log(`📨 WebSocket message type: ${msg.type}`);

          // ═══════════════════════════════════════════════════════════
          // CHAT MESSAGE
          // ═══════════════════════════════════════════════════════════
          if (msg.type === 'chat') {
            // Schema validation already ensures message exists and has minLength: 1
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

            // Save assistant message to session history
            if (response.type === 'dialogue' || response.type === 'search') {
              const assistantMessage = {
                role: 'assistant' as const,
                content: response.output || '',
                timestamp: new Date(),
              };
              await container.sessionService.addMessage(
                msg.session_id,
                assistantMessage
              );
            }

            ws.send(response);
            return;
          }

          // ═══════════════════════════════════════════════════════════
          // PRODUCT DETAILS REQUEST
          // ═══════════════════════════════════════════════════════════
          if (msg.type === 'product-details') {
            // Schema validation already ensures page_token exists and has minLength: 1
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
          // UNKNOWN MESSAGE TYPE (should never reach here due to schema)
          // ═══════════════════════════════════════════════════════════
          ws.send({
            type: 'error',
            error: 'INVALID_MESSAGE_TYPE',
            message: 'Invalid message type',
            session_id: msg.session_id,
          });
        } catch (error: any) {
          console.error('WebSocket message handler error:', error);

          // Extract session_id safely
          const msg = message as any;
          const sessionId = msg?.session_id || '';

          ws.send({
            type: 'error',
            error: 'SERVER_ERROR',
            message: error.message || 'An error occurred',
            session_id: sessionId,
          });
        }
      },
    });
