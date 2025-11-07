/**
 * WebSocket Model - TypeBox schemas and type definitions
 * Following Elysia Best Practices: namespace for DTOs
 */

import { t } from 'elysia';

export namespace WebSocketModel {
  // ═══════════════════════════════════════════════════════════
  // WEBSOCKET MESSAGE TYPES
  // ═══════════════════════════════════════════════════════════

  // Incoming chat message
  export const chatMessage = t.Object({
    type: t.Literal('chat'),
    session_id: t.String(),
    message: t.String(),
    country: t.Optional(t.String()),
    language: t.Optional(t.String()),
    currency: t.Optional(t.String()),
    new_search: t.Optional(t.Boolean()),
    current_category: t.Optional(t.String()),
    access_token: t.Optional(t.String()),
  });

  export type chatMessage = typeof chatMessage.static;

  // Incoming product details request
  export const productDetailsMessage = t.Object({
    type: t.Literal('product-details'),
    session_id: t.String(),
    page_token: t.String(),
    country: t.Optional(t.String()),
    access_token: t.Optional(t.String()),
  });

  export type productDetailsMessage = typeof productDetailsMessage.static;

  // Union of all incoming message types
  export const incomingMessage = t.Union([chatMessage, productDetailsMessage]);

  export type incomingMessage = typeof incomingMessage.static;

  // ═══════════════════════════════════════════════════════════
  // WEBSOCKET RESPONSE TYPES
  // ═══════════════════════════════════════════════════════════

  // Product card
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

  // Search state
  export const searchState = t.Object({
    status: t.String(),
    category: t.Optional(t.String()),
    can_continue: t.Boolean(),
    search_count: t.Number(),
    max_searches: t.Number(),
  });

  export type searchState = typeof searchState.static;

  // Dialogue response
  export const dialogueResponse = t.Object({
    type: t.Literal('dialogue'),
    output: t.String(),
    quick_replies: t.Optional(t.Array(t.String())),
    session_id: t.String(),
    message_count: t.Optional(t.Number()),
    search_state: t.Optional(searchState),
  });

  export type dialogueResponse = typeof dialogueResponse.static;

  // Search response
  export const searchResponse = t.Object({
    type: t.Literal('search'),
    output: t.String(),
    products: t.Array(productCard),
    search_type: t.Optional(t.String()),
    session_id: t.String(),
    message_count: t.Optional(t.Number()),
    search_state: t.Optional(searchState),
  });

  export type searchResponse = typeof searchResponse.static;

  // Product details response
  export const productDetailsResponse = t.Object({
    type: t.Literal('product_details'),
    product_details: t.Any(),
    session_id: t.String(),
  });

  export type productDetailsResponse = typeof productDetailsResponse.static;

  // Error response
  export const errorResponse = t.Object({
    type: t.Literal('error'),
    error: t.String(),
    message: t.String(),
    session_id: t.String(),
  });

  export type errorResponse = typeof errorResponse.static;

  // Union of all outgoing response types
  export const outgoingResponse = t.Union([
    dialogueResponse,
    searchResponse,
    productDetailsResponse,
    errorResponse,
  ]);

  export type outgoingResponse = typeof outgoingResponse.static;
}
