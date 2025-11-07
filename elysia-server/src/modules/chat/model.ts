/**
 * Chat Model - TypeBox schemas and type definitions
 * Following Elysia Best Practices: namespace for DTOs
 */

import { t } from 'elysia';

export namespace ChatModel {
  // ═══════════════════════════════════════════════════════════
  // CHAT REQUEST/RESPONSE
  // ═══════════════════════════════════════════════════════════

  // POST /api/chat - Request body
  export const chatRequest = t.Object({
    session_id: t.String(),
    message: t.String(),
    country: t.String(),
    language: t.String(),
    currency: t.String(),
    new_search: t.Boolean(),
    current_category: t.String(),
  });

  export type chatRequest = typeof chatRequest.static;

  // Product card schema
  export const productCard = t.Object({
    name: t.String(),
    price: t.String(),
    old_price: t.Optional(t.String()),
    link: t.String(),
    image: t.String(),
    description: t.Optional(t.String()),
    badge: t.Optional(t.String()),
    page_token: t.String(),
  });

  export type productCard = typeof productCard.static;

  // Search state schema
  export const searchState = t.Object({
    status: t.String(),
    category: t.Optional(t.String()),
    can_continue: t.Boolean(),
    search_count: t.Number(),
    max_searches: t.Number(),
    message: t.Optional(t.String()),
  });

  export type searchState = typeof searchState.static;

  // Chat response - dialogue type
  export const dialogueResponse = t.Object({
    type: t.Literal('dialogue'),
    output: t.String(),
    quick_replies: t.Array(t.String()),
    session_id: t.String(),
    message_count: t.Number(),
    search_state: searchState,
  });

  export type dialogueResponse = typeof dialogueResponse.static;

  // Chat response - search type
  export const searchResponse = t.Object({
    type: t.Literal('search'),
    output: t.String(),
    products: t.Array(productCard),
    search_type: t.Optional(t.String()),
    session_id: t.String(),
    message_count: t.Number(),
    search_state: searchState,
  });

  export type searchResponse = typeof searchResponse.static;

  // Chat response - error type
  export const errorResponse = t.Object({
    type: t.Literal('error'),
    output: t.String(),
    session_id: t.String(),
    message_count: t.Number(),
    search_state: t.Optional(searchState),
  });

  export type errorResponse = typeof errorResponse.static;

  // Union of all response types
  export const chatResponse = t.Union([
    dialogueResponse,
    searchResponse,
    errorResponse,
  ]);

  export type chatResponse = typeof chatResponse.static;

  // ═══════════════════════════════════════════════════════════
  // GET MESSAGES
  // ═══════════════════════════════════════════════════════════

  // GET /api/chat/messages - Query params
  export const getMessagesQuery = t.Object({
    session_id: t.String(),
  });

  export type getMessagesQuery = typeof getMessagesQuery.static;

  // Message schema
  export const message = t.Object({
    role: t.Union([t.Literal('user'), t.Literal('assistant')]),
    content: t.String(),
    timestamp: t.Optional(t.Date()),
  });

  export type message = typeof message.static;

  // GET /api/chat/messages - Response
  export const getMessagesResponse = t.Object({
    session_id: t.String(),
    messages: t.Array(message),
    count: t.Number(),
  });

  export type getMessagesResponse = typeof getMessagesResponse.static;

  // ═══════════════════════════════════════════════════════════
  // PRODUCT DETAILS
  // ═══════════════════════════════════════════════════════════

  // POST /api/product-details - Request body
  export const productDetailsRequest = t.Object({
    page_token: t.String(),
    country: t.String(),
  });

  export type productDetailsRequest = typeof productDetailsRequest.static;

  // Specification schema
  export const specification = t.Object({
    title: t.String(),
    value: t.String(),
  });

  export type specification = typeof specification.static;

  // Offer schema
  export const offer = t.Object({
    merchant: t.String(),
    logo: t.Optional(t.String()),
    price: t.String(),
    extracted_price: t.Optional(t.Number()),
    currency: t.Optional(t.String()),
    link: t.String(),
    title: t.Optional(t.String()),
    availability: t.Optional(t.String()),
    shipping: t.Optional(t.String()),
    shipping_extracted: t.Optional(t.Number()),
    total: t.Optional(t.String()),
    extracted_total: t.Optional(t.Number()),
    rating: t.Optional(t.Number()),
    reviews: t.Optional(t.Number()),
    payment_methods: t.Optional(t.String()),
    tag: t.Optional(t.String()),
    details_and_offers: t.Optional(t.Array(t.String())),
  });

  export type offer = typeof offer.static;

  // Rating breakdown item
  export const ratingBreakdownItem = t.Object({
    stars: t.Number(),
    amount: t.Number(),
  });

  export type ratingBreakdownItem = typeof ratingBreakdownItem.static;

  // Product details response
  export const productDetailsResponse = t.Object({
    type: t.String(),
    title: t.String(),
    price: t.String(),
    rating: t.Optional(t.Number()),
    reviews: t.Optional(t.Number()),
    description: t.Optional(t.String()),
    images: t.Optional(t.Array(t.String())),
    specifications: t.Optional(t.Array(specification)),
    variants: t.Optional(t.Array(t.Any())),
    offers: t.Array(offer),
    videos: t.Optional(t.Array(t.Any())),
    more_options: t.Optional(t.Array(t.Any())),
    rating_breakdown: t.Optional(t.Array(ratingBreakdownItem)),
  });

  export type productDetailsResponse = typeof productDetailsResponse.static;

  // ═══════════════════════════════════════════════════════════
  // ERROR RESPONSES
  // ═══════════════════════════════════════════════════════════

  export const invalidRequestError = t.Object({
    error: t.Literal('invalid_request'),
    message: t.String(),
  });

  export type invalidRequestError = typeof invalidRequestError.static;

  export const serverError = t.Object({
    error: t.Literal('server_error'),
    message: t.String(),
  });

  export type serverError = typeof serverError.static;
}
