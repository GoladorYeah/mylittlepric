/**
 * Chat Service - Business logic for chat processing
 * Following Elysia Best Practices: abstract static classes
 */

import type { SessionService } from '../../services/session.service';
import type { GeminiService } from '../../services/gemini.service';
import type { SerpService } from '../../services/serp.service';
import type { Message, ChatSession } from '../../types';
import type { ChatModel } from './model';
import type { Config } from '../../config';

/**
 * Chat Service
 * Handles chat message processing and product search logic
 */
export abstract class Chat {
  /**
   * Process chat message with full context
   * Returns chat response ready to send to client
   */
  static async processMessage(
    sessionId: string,
    message: string,
    country: string,
    language: string,
    currency: string,
    currentCategory: string,
    sessionService: SessionService,
    geminiService: GeminiService,
    serpService: SerpService,
    config: Config
  ): Promise<ChatModel.chatResponse> {
    try {
      // Get or create session
      const session = await sessionService.getOrCreateSession(
        sessionId,
        country || config.defaultCountry,
        language || config.defaultLanguage,
        currency || config.defaultCurrency
      );

      // Get conversation history
      const history = await sessionService.getMessages(sessionId);

      // Add user message
      const userMessage: Message = {
        role: 'user',
        content: message,
        timestamp: new Date(),
      };
      await sessionService.addMessage(sessionId, userMessage);

      // Process with Gemini
      const { response: geminiResponse } =
        await geminiService.processMessageWithContext(
          message,
          history,
          session.country_code,
          session.language_code,
          currentCategory || session.search_state.category,
          session.search_state.last_product
        );

      // Update message count
      session.message_count++;

      // Handle search response
      if (geminiResponse.response_type === 'search') {
        return await Chat.handleSearchResponse(
          session,
          message,
          geminiResponse,
          sessionService,
          geminiService,
          serpService,
          config
        );
      }

      // Handle dialogue response
      return await Chat.handleDialogueResponse(
        session,
        geminiResponse,
        sessionService,
        config
      );
    } catch (error: any) {
      console.error('Chat processing error:', error);
      return {
        type: 'error',
        output: 'An error occurred while processing your request.',
        session_id: sessionId,
        message_count: 0,
      };
    }
  }

  /**
   * Handle search type response
   */
  private static async handleSearchResponse(
    session: ChatSession,
    message: string,
    geminiResponse: any,
    sessionService: SessionService,
    geminiService: GeminiService,
    serpService: SerpService,
    config: Config
  ): Promise<ChatModel.searchResponse | ChatModel.errorResponse> {
    // Check search limit
    if (
      session.search_state.search_count >= config.maxSearchesPerSession
    ) {
      await sessionService.updateSession(session);

      return {
        type: 'error',
        output: `You have reached the maximum number of searches (${config.maxSearchesPerSession}) for this session.`,
        session_id: session.session_id,
        message_count: session.message_count,
        search_state: {
          status: 'idle',
          can_continue: false,
          search_count: session.search_state.search_count,
          max_searches: config.maxSearchesPerSession,
        },
      };
    }

    // Translate query to English for better search results
    const englishQuery = await geminiService.translateToEnglish(
      geminiResponse.search_phrase || message
    );

    // Search products
    const products = await serpService.searchProducts(
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

    await sessionService.updateSession(session);

    return {
      type: 'search',
      output: geminiResponse.output,
      products,
      search_type: geminiResponse.search_type,
      session_id: session.session_id,
      message_count: session.message_count,
      search_state: {
        status: session.search_state.status,
        category: session.search_state.category,
        can_continue:
          session.search_state.search_count < config.maxSearchesPerSession,
        search_count: session.search_state.search_count,
        max_searches: config.maxSearchesPerSession,
      },
    };
  }

  /**
   * Handle dialogue type response
   */
  private static async handleDialogueResponse(
    session: ChatSession,
    geminiResponse: any,
    sessionService: SessionService,
    config: Config
  ): Promise<ChatModel.dialogueResponse> {
    await sessionService.updateSession(session);

    const assistantMessage: Message = {
      role: 'assistant',
      content: geminiResponse.output,
      timestamp: new Date(),
    };
    await sessionService.addMessage(session.session_id, assistantMessage);

    return {
      type: 'dialogue',
      output: geminiResponse.output,
      quick_replies: geminiResponse.quick_replies || [],
      session_id: session.session_id,
      message_count: session.message_count,
      search_state: {
        status: session.search_state.status,
        category: session.search_state.category,
        can_continue:
          session.search_state.search_count < config.maxSearchesPerSession,
        search_count: session.search_state.search_count,
        max_searches: config.maxSearchesPerSession,
      },
    };
  }

  /**
   * Get product details by page token
   */
  static async getProductDetails(
    pageToken: string,
    country: string,
    serpService: SerpService,
    config: Config
  ): Promise<any> {
    return serpService.getProductDetails(
      pageToken,
      country || config.defaultCountry
    );
  }
}
