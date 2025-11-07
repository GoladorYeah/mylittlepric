/**
 * Chat module - handles chat and product search endpoints
 */

import { Elysia, t } from 'elysia';
import type { Container } from '../../container';
import type { ChatRequest, ChatResponse, Message } from '../../types';

export const chatModule = (container: Container) =>
  new Elysia({ prefix: '/api' })
    // Chat endpoint
    .post(
      '/chat',
      async ({ body }): Promise<ChatResponse> => {
        const {
          session_id,
          message,
          country,
          language,
          currency,
          new_search,
          current_category,
        } = body;

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
              const limitResponse: ChatResponse = {
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
              return limitResponse;
            }

            // Translate query to English
            const englishQuery =
              await container.geminiService.translateToEnglish(
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

            return {
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
            };
          } else {
            // Dialogue response
            await container.sessionService.updateSession(session);

            const assistantMessage: Message = {
              role: 'assistant',
              content: geminiResponse.output,
              timestamp: new Date(),
            };
            await container.sessionService.addMessage(session_id, assistantMessage);

            return {
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
            };
          }
        } catch (error: any) {
          console.error('Chat error:', error);
          return {
            type: 'error',
            output: 'An error occurred while processing your request.',
            session_id,
            message_count: 0,
          };
        }
      },
      {
        body: t.Object({
          session_id: t.String(),
          message: t.String(),
          country: t.String(),
          language: t.String(),
          currency: t.String(),
          new_search: t.Boolean(),
          current_category: t.String(),
        }),
      }
    )

    // Get session messages
    .get(
      '/chat/messages',
      async ({ query, set }) => {
        try {
          const sessionId = query.session_id;

          if (!sessionId) {
            set.status = 400;
            return {
              error: 'invalid_request',
              message: 'session_id is required',
            };
          }

          // Get messages from session
          const messages = await container.sessionService.getMessages(sessionId);

          return {
            session_id: sessionId,
            messages,
            count: messages.length,
          };
        } catch (error: any) {
          console.error('Get messages error:', error);
          set.status = 500;
          return {
            error: 'server_error',
            message: 'Failed to retrieve messages',
          };
        }
      },
      {
        query: t.Object({
          session_id: t.String(),
        }),
      }
    )

    // Product details endpoint
    .post(
      '/product-details',
      async ({ body }) => {
        try {
          const details = await container.serpService.getProductDetails(
            body.page_token,
            body.country || container.config.defaultCountry
          );

          return details;
        } catch (error: any) {
          console.error('Product details error:', error);
          return {
            error: true,
            message: 'Failed to get product details',
          };
        }
      },
      {
        body: t.Object({
          page_token: t.String(),
          country: t.String(),
        }),
      }
    );
