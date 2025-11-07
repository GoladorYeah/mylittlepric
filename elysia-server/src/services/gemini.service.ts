/**
 * GeminiService handles AI conversation processing using Google Gemini
 * Includes grounding strategy, prompt management, and response parsing
 */

import { GoogleGenerativeAI } from '@google/generative-ai';
import type { Config } from '../config';
import type { KeyRotator } from '../utils/key-rotator';
import type { EmbeddingService } from './embedding.service';
import type { GroundingStrategy } from './grounding-strategy.service';
import type { GeminiResponse, Message, ProductInfo } from '../types';

export class GeminiService {
  private genai: GoogleGenerativeAI;
  private keyRotator: KeyRotator;
  private config: Config;
  private embedding: EmbeddingService;
  private groundingStrategy: GroundingStrategy;
  private currentKeyIndex: number = 0;

  constructor(
    keyRotator: KeyRotator,
    config: Config,
    embedding: EmbeddingService,
    groundingStrategy: GroundingStrategy
  ) {
    this.keyRotator = keyRotator;
    this.config = config;
    this.embedding = embedding;
    this.groundingStrategy = groundingStrategy;

    // Initialize with first key
    const { key, index } = this.keyRotator.getKeyByIndex(0);
    this.genai = new GoogleGenerativeAI(key);
    this.currentKeyIndex = index;
  }

  /**
   * Rotate to next API key
   */
  private async rotateKey(): Promise<void> {
    const { key, index } = await this.keyRotator.getNextKey();
    this.genai = new GoogleGenerativeAI(key);
    this.currentKeyIndex = index;
  }

  /**
   * Process user message with conversation context
   */
  async processMessageWithContext(
    userMessage: string,
    conversationHistory: Message[],
    country: string,
    language: string,
    currentCategory: string,
    lastProduct?: ProductInfo
  ): Promise<{ response: GeminiResponse; tokensUsed: number }> {
    // Detect category if not set
    if (!currentCategory) {
      currentCategory = await this.embedding.detectCategory(userMessage);
    }

    // Build system prompt
    const systemPrompt = this.buildSystemPrompt(
      country,
      language,
      currentCategory,
      lastProduct
    );

    // Build conversation context
    const conversationContext = this.buildConversationContext(
      conversationHistory
    );

    // Check if grounding should be used
    const useGrounding = this.config.geminiUseGrounding
      ? await this.groundingStrategy.shouldUseGrounding(
          userMessage,
          conversationHistory,
          currentCategory
        )
      : { useGrounding: false, confidence: 0, reason: 'disabled' };

    // Prepare full prompt
    const fullPrompt = `${systemPrompt}\n\n# CONVERSATION HISTORY:\n${conversationContext}\n\nCurrent user message: ${userMessage}\n\nCRITICAL INSTRUCTIONS:\n- You MUST respond with valid JSON only\n- If using grounding/search results, incorporate the information naturally\n- ALWAYS end your response with valid JSON in this exact format:\n{"response_type":"dialogue","output":"...","quick_replies":[...],"category":"..."}\nOR\n{"response_type":"search","search_phrase":"...","search_type":"...","category":"..."}`;

    // Generate content
    const model = this.genai.getGenerativeModel({
      model: this.config.geminiModel,
    });

    const generationConfig: any = {
      temperature: this.config.geminiTemperature,
      maxOutputTokens: this.config.geminiMaxOutputTokens,
    };

    // Add grounding if needed
    if (useGrounding.useGrounding) {
      generationConfig.tools = [{ googleSearch: {} }];
      console.log(
        `🔍 Using grounding (confidence: ${useGrounding.confidence.toFixed(2)}, reason: ${useGrounding.reason})`
      );
    } else {
      generationConfig.responseMimeType = 'application/json';
    }

    try {
      const result = await model.generateContent({
        contents: [{ role: 'user', parts: [{ text: fullPrompt }] }],
        generationConfig,
      });

      const response = result.response;
      const text = response.text();

      // Parse JSON response
      const geminiResponse = this.parseGeminiResponse(text);

      // Calculate tokens used (approximation)
      const tokensUsed = this.estimateTokens(fullPrompt + text);

      return { response: geminiResponse, tokensUsed };
    } catch (error: any) {
      // Handle quota errors by rotating key
      if (
        error.message?.includes('quota') ||
        error.message?.includes('429') ||
        error.message?.includes('RESOURCE_EXHAUSTED')
      ) {
        console.log('⚠️ Quota exceeded, rotating key...');
        await this.rotateKey();

        // Retry with new key
        return this.processMessageWithContext(
          userMessage,
          conversationHistory,
          country,
          language,
          currentCategory,
          lastProduct
        );
      }

      throw error;
    }
  }

  /**
   * Build system prompt based on category and context
   */
  private buildSystemPrompt(
    country: string,
    language: string,
    category: string,
    lastProduct?: ProductInfo
  ): string {
    const lastProductStr = lastProduct
      ? `${lastProduct.name} (${lastProduct.price})`
      : '';

    return `You are a helpful shopping assistant for ${country} market, speaking ${language}.
Category: ${category || 'general'}
Last discussed product: ${lastProductStr}

Your role is to:
1. Help users find the right products through natural conversation
2. Ask clarifying questions when needed
3. Trigger product searches when you have enough information
4. Provide product recommendations

Response format:
- For dialogue: {"response_type":"dialogue","output":"...","quick_replies":[...],"category":"..."}
- For search: {"response_type":"search","search_phrase":"...","search_type":"exact|parameters|category","category":"..."}

Available search types:
- exact: Specific product with brand and model
- parameters: Product with specific requirements
- category: General product category

Always be helpful, concise, and guide the conversation naturally.`;
  }

  /**
   * Build conversation context from history
   */
  private buildConversationContext(history: Message[]): string {
    if (history.length === 0) {
      return 'No previous conversation.';
    }

    return history
      .slice(-6) // Last 6 messages
      .map((msg) => `${msg.role === 'user' ? 'User' : 'Assistant'}: ${msg.content}`)
      .join('\n');
  }

  /**
   * Parse Gemini JSON response
   */
  private parseGeminiResponse(text: string): GeminiResponse {
    try {
      // Try to extract JSON from the text
      const jsonMatch = text.match(/\{[\s\S]*\}/);
      if (!jsonMatch) {
        throw new Error('No JSON found in response');
      }

      const parsed = JSON.parse(jsonMatch[0]);

      // Validate required fields
      if (!parsed.response_type) {
        throw new Error('Missing response_type');
      }

      return parsed as GeminiResponse;
    } catch (error) {
      console.error('Failed to parse Gemini response:', error);
      // Return fallback dialogue response
      return {
        response_type: 'dialogue',
        output: text || 'I apologize, but I encountered an error. Please try again.',
        quick_replies: [],
        category: '',
      };
    }
  }

  /**
   * Estimate tokens used (rough approximation)
   */
  private estimateTokens(text: string): number {
    // Rough approximation: 1 token ≈ 4 characters
    return Math.ceil(text.length / 4);
  }

  /**
   * Translate text to English for better search results
   */
  async translateToEnglish(text: string): Promise<string> {
    try {
      const model = this.genai.getGenerativeModel({
        model: this.config.geminiFallbackModel,
      });

      const result = await model.generateContent({
        contents: [
          {
            role: 'user',
            parts: [
              {
                text: `Translate the following product search query to English. Return ONLY the translated text, no explanations:\n\n${text}`,
              },
            ],
          },
        ],
        generationConfig: {
          temperature: this.config.geminiTranslationTemperature,
          maxOutputTokens: this.config.geminiTranslationMaxTokens,
        },
      });

      return result.response.text().trim();
    } catch (error) {
      console.error('Translation error:', error);
      return text; // Return original if translation fails
    }
  }
}
