/**
 * WebSocket module - handles real-time chat communication
 */

import { Elysia, t } from 'elysia';
import { ws } from '@elysiajs/ws';
import type { Container } from '../../container';
import type { ChatRequest, ChatResponse, Message } from '../../types';

interface WSMessage {
  type: 'chat' | 'product-details';
  session_id: string;
  message?: string;
  country?: string;
  language?: string;
  currency?: string;
  new_search?: boolean;
  current_category?: string;
  page_token?: string;
  access_token?: string;
}

interface WSResponse {
  type: string;
  output?: string;
  quick_replies?: string[];
  products?: any[];
  search_type?: string;
  session_id: string;
  message_count?: number;
  search_state?: any;
  product_details?: any;
  error?: string;
  message?: string;
}

export const websocketModule = (container: Container) =>
  new Elysia()
    .use(
      ws({
        path: '/ws',
      })
    )
    .ws('/ws', {
      message: async (ws, msg: WSMessage) => {
        try {
          console.log(`📨 WebSocket message type: ${msg.type}`);

          if (msg.type === 'chat') {
            await handleChatMessage(ws, msg, container);
          } else if (msg.type === 'product-details') {
            await handleProductDetails(ws, msg, container);
          } else {
            ws.send({
              type: 'error',
              error: 'INVALID_MESSAGE_TYPE',
              message: 'Invalid message type',
              session_id: msg.session_id || '',
            });
          }
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

/**
 * Handle chat message
 */
async function handleChatMessage(
  ws: any,
  msg: WSMessage,
  container: Container
): Promise<void> {
  const {
    session_id,
    message,
    country,
    language,
    currency,
    new_search,
    current_category,
  } = msg;

  if (!message) {
    ws.send({
      type: 'error',
      error: 'INVALID_REQUEST',
      message: 'Message is required',
      session_id,
    });
    return;
  }

  try {
    // Get or create session
    const session = await container.sessionService.getOrCreateSession(
      session_id,
      country || container.config.defaultCountry,
      language || container.config.defaultLanguage,
      currency || container.config.defaultCurrency
    );

    // Get conversation history
    const history = await container.sessionService.getMessages(session_id);

    // Add user message
    const userMessage: Message = {
      role: 'user',
      content: message,
      timestamp: new Date(),
    };
    await container.sessionService.addMessage(session_id, userMessage);

    // Process with Gemini
    const { response: geminiResponse } =
      await container.geminiService.processMessageWithContext(
        message,
        history,
        session.country_code,
        session.language_code,
        current_category || session.search_state.category,
        session.search_state.last_product
      );

    // Update message count
    session.message_count++;

    // Handle response based on type
    if (geminiResponse.response_type === 'search') {
      // Check search limit
      if (
        session.search_state.search_count >=
        container.config.maxSearchesPerSession
      ) {
        const limitResponse: WSResponse = {
          type: 'error',
          output: `You have reached the maximum number of searches (${container.config.maxSearchesPerSession}) for this session.`,
          session_id,
          message_count: session.message_count,
          search_state: {
            status: 'idle',
            can_continue: false,
            search_count: session.search_state.search_count,
            max_searches: container.config.maxSearchesPerSession,
          },
        };

        await container.sessionService.updateSession(session);
        ws.send(limitResponse);
        return;
      }

      // Translate query to English
      const englishQuery = await container.geminiService.translateToEnglish(
        geminiResponse.search_phrase || message
      );

      // Search products
      const products = await container.serpService.searchProducts(
        englishQuery,
        geminiResponse.search_type || 'parameters',
        session.country_code,
        geminiResponse.min_price,
        geminiResponse.max_price
      );

      // Update search state
      session.search_state.search_count++;
      session.search_state.category = geminiResponse.category || '';
      session.search_state.status = 'completed';

      if (products.length > 0) {
        session.search_state.last_product = {
          name: products[0].name,
          price: parseFloat(products[0].price.replace(/[^0-9.]/g, '')) || 0,
        };
      }

      await container.sessionService.updateSession(session);

      ws.send({
        type: 'search',
        output: geminiResponse.output,
        products,
        search_type: geminiResponse.search_type,
        session_id,
        message_count: session.message_count,
        search_state: {
          status: session.search_state.status,
          category: session.search_state.category,
          can_continue:
            session.search_state.search_count <
            container.config.maxSearchesPerSession,
          search_count: session.search_state.search_count,
          max_searches: container.config.maxSearchesPerSession,
        },
      });
    } else {
      // Dialogue response
      await container.sessionService.updateSession(session);

      const assistantMessage: Message = {
        role: 'assistant',
        content: geminiResponse.output,
        timestamp: new Date(),
      };
      await container.sessionService.addMessage(session_id, assistantMessage);

      ws.send({
        type: 'dialogue',
        output: geminiResponse.output,
        quick_replies: geminiResponse.quick_replies || [],
        session_id,
        message_count: session.message_count,
        search_state: {
          status: session.search_state.status,
          category: session.search_state.category,
          can_continue:
            session.search_state.search_count <
            container.config.maxSearchesPerSession,
          search_count: session.search_state.search_count,
          max_searches: container.config.maxSearchesPerSession,
        },
      });
    }
  } catch (error: any) {
    console.error('Chat processing error:', error);
    ws.send({
      type: 'error',
      error: 'PROCESSING_ERROR',
      message: 'An error occurred while processing your message.',
      session_id,
    });
  }
}

/**
 * Handle product details request
 */
async function handleProductDetails(
  ws: any,
  msg: WSMessage,
  container: Container
): Promise<void> {
  const { page_token, country, session_id } = msg;

  if (!page_token) {
    ws.send({
      type: 'error',
      error: 'INVALID_REQUEST',
      message: 'page_token is required',
      session_id,
    });
    return;
  }

  try {
    const details = await container.serpService.getProductDetails(
      page_token,
      country || container.config.defaultCountry
    );

    ws.send({
      type: 'product_details',
      product_details: details,
      session_id,
    });
  } catch (error: any) {
    console.error('Product details error:', error);
    ws.send({
      type: 'error',
      error: 'PRODUCT_DETAILS_ERROR',
      message: 'Failed to get product details',
      session_id,
    });
  }
}
