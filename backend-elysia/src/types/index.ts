/**
 * Type definitions for MyLittlePrice Backend
 */

// ═══════════════════════════════════════════════════════════
// CHAT REQUEST/RESPONSE MODELS
// ═══════════════════════════════════════════════════════════

export interface ChatRequest {
  session_id: string;
  message: string;
  country: string;
  language: string;
  currency: string;
  new_search: boolean;
  current_category: string;
}

export interface ChatResponse {
  type: string;
  output?: string;
  quick_replies?: string[];
  products?: ProductCard[];
  search_type?: string;
  session_id: string;
  message_count: number;
  search_state?: SearchStateResponse;
}

export interface SearchStateResponse {
  status: string;
  category?: string;
  can_continue: boolean;
  search_count: number;
  max_searches: number;
  message?: string;
}

// ═══════════════════════════════════════════════════════════
// AI RESPONSE MODELS
// ═══════════════════════════════════════════════════════════

export interface GeminiResponse {
  response_type: 'dialogue' | 'search' | 'api_request';
  output: string;
  quick_replies: string[];
  search_phrase?: string;
  search_type?: 'exact' | 'parameters' | 'category';
  category?: string;
  price_filter?: 'cheaper' | 'expensive';
  min_price?: number;
  max_price?: number;
  product_type?: string;
  brand?: string;
  confidence?: number;
  requires_input?: boolean;
  api?: string;
  params?: Record<string, any>;
}

export interface SerpConfig {
  type: string;
  query: string;
  country: string;
  language: string;
  max_results: number;
}

// ═══════════════════════════════════════════════════════════
// SESSION MODELS
// ═══════════════════════════════════════════════════════════

export interface ChatSession {
  id: string;
  session_id: string;
  country_code: string;
  language_code: string;
  currency: string;
  message_count: number;
  search_state: SearchState;
  cycle_state: CycleState;
  conversation_context?: ConversationContext;
  created_at: Date;
  updated_at: Date;
  expires_at: Date;
}

export interface SearchState {
  status: SearchStatus;
  category: string;
  last_search_time?: Date;
  search_count: number;
  last_product?: ProductInfo;
}

export type SearchStatus = 'idle' | 'in_progress' | 'completed';

export interface CycleState {
  cycle_id: number;
  iteration: number;
  cycle_history: CycleMessage[];
  last_cycle_context?: LastCycleContext;
  last_defined: string[];
  prompt_id: string;
  prompt_hash: string;
}

export interface CycleMessage {
  role: 'user' | 'assistant';
  content: string;
  timestamp: Date;
}

export interface LastCycleContext {
  groups: string[];
  subgroups: string[];
  products: ProductInfo[];
  last_request: string;
}

export interface ConversationContext {
  summary: string;
  preferences: UserPreferences;
  last_search?: SearchContext;
  exclusions?: string[];
  updated_at: Date;
}

export interface UserPreferences {
  price_range?: PriceRange;
  brands?: string[];
  features?: string[];
  requirements?: string[];
}

export interface PriceRange {
  min?: number;
  max?: number;
  currency: string;
}

export interface SearchContext {
  query: string;
  category: string;
  products_shown?: ProductInfo[];
  user_feedback?: string;
  timestamp: Date;
}

// ═══════════════════════════════════════════════════════════
// PRODUCT MODELS
// ═══════════════════════════════════════════════════════════

export interface ProductInfo {
  name: string;
  price: number;
}

export interface ProductCard {
  name: string;
  price: string;
  old_price?: string;
  link: string;
  image: string;
  description?: string;
  badge?: string;
  page_token: string;
}

export interface ProductDetailsRequest {
  page_token: string;
  country: string;
}

export interface ProductDetailsResponse {
  type: string;
  title: string;
  price: string;
  rating?: number;
  reviews?: number;
  description?: string;
  images?: string[];
  specifications?: Specification[];
  variants?: Variant[];
  offers: Offer[];
  videos?: any[];
  more_options?: any[];
  rating_breakdown?: RatingBreakdownItem[];
}

export interface Specification {
  title: string;
  value: string;
}

export interface Variant {
  title: string;
  items: any[];
}

export interface Offer {
  merchant: string;
  logo?: string;
  price: string;
  extracted_price?: number;
  currency?: string;
  link: string;
  title?: string;
  availability?: string;
  shipping?: string;
  shipping_extracted?: number;
  total?: string;
  extracted_total?: number;
  rating?: number;
  reviews?: number;
  payment_methods?: string;
  tag?: string;
  details_and_offers?: string[];
  monthly_payment_duration?: number;
  down_payment?: string;
}

export interface RatingBreakdownItem {
  stars: number;
  amount: number;
}

// ═══════════════════════════════════════════════════════════
// AUTH MODELS
// ═══════════════════════════════════════════════════════════

export interface User {
  id: string;
  email: string;
  name: string;
  picture?: string;
  provider: string;
  provider_id: string;
  created_at: Date;
  updated_at: Date;
}

export interface AuthTokens {
  access_token: string;
  refresh_token: string;
  expires_in: number;
}

export interface GoogleUserInfo {
  id: string;
  email: string;
  verified_email: boolean;
  name: string;
  given_name: string;
  family_name: string;
  picture: string;
}

// ═══════════════════════════════════════════════════════════
// SEARCH HISTORY MODELS
// ═══════════════════════════════════════════════════════════

export interface SearchHistoryItem {
  id: string;
  user_id: string;
  session_id: string;
  query: string;
  category: string;
  search_type: string;
  country: string;
  language: string;
  result_count: number;
  created_at: Date;
}

// ═══════════════════════════════════════════════════════════
// MESSAGE MODELS
// ═══════════════════════════════════════════════════════════

export interface Message {
  role: 'user' | 'assistant';
  content: string;
  timestamp?: Date;
}

// ═══════════════════════════════════════════════════════════
// ERROR MODELS
// ═══════════════════════════════════════════════════════════

export interface ErrorResponse {
  error: string;
  message: string;
}

// ═══════════════════════════════════════════════════════════
// UTILITY TYPES
// ═══════════════════════════════════════════════════════════

export interface KeyRotatorStats {
  key_index: number;
  total_usage: string;
  success_count: string;
  failure_count: string;
  avg_response_time_ms?: number;
}
