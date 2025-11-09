export interface Product {
  position: number;
  title: string;
  link: string;
  product_link: string;
  product_id: string;
  serpapi_product_api: string;
  serpapi_product_api_comparative?: string;
  source: string;
  price: string;
  extracted_price: number;
  rating?: number;
  reviews?: number;
  thumbnail: string;
  delivery?: string;
  tag?: string;
  extensions?: string[];
  currency?: string;
  page_token?: string;
  relevance_score?: number;
}

export interface ChatMessage {
  id: string;
  role: "user" | "assistant";
  content: string;
  timestamp: number;
  quick_replies?: string[];
  products?: Product[];
  search_type?: string;
}

export interface SessionMessage {
  role: string;
  content: string;
  timestamp?: string;
  quick_replies?: string[];
  products?: Product[];
  search_type?: string;
}

export interface SessionResponse {
  session_id: string;
  messages: SessionMessage[];
  search_state?: {
    category?: string;
    status?: string;
    last_product?: {
      name: string;
      price: string;
    };
  };
}

export interface ChatResponse {
  session_id: string;
  message: string;
  quick_replies?: string[];
  products?: Product[];
  response_type?: string;
  search_type?: string;
}

export interface SearchHistoryItem {
  id: string;
  query: string;
  timestamp: number;
  category?: string;
  productsCount?: number;
  sessionId: string;
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

export interface SearchHistoryRecord {
  id: string;
  user_id?: string;
  session_id?: string;
  search_query: string;
  optimized_query?: string;
  search_type: string;
  category?: string;
  country_code: string;
  language_code: string;
  currency: string;
  result_count: number;
  products_found?: ProductCard[];
  clicked_product_id?: string;
  created_at: string;
  expires_at?: string;
}

export interface User {
  id: string;
  email: string;
  name: string;
  picture?: string;
}

export interface AuthState {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
}

export interface ProductDetailsResponse {
  type: string;
  title: string;
  price: string;
  rating?: number;
  reviews?: number;
  description?: string;
  images?: string[];
  specifications?: { title: string; value: string }[];
  variants?: { title: string; items: any[] }[];
  offers: {
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
  }[];
  videos?: any[];
  more_options?: any[];
  rating_breakdown?: { stars: number; amount: number }[];
}