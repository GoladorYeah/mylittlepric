export interface ChatMessage {
  id: string;
  role: "user" | "assistant";
  content: string;
  timestamp: number;
  quick_replies?: string[];
  products?: ProductCard[];
  search_type?: string;
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

export interface ChatRequest {
  session_id: string;
  message: string;
  country: string;
  language: string;
  currency: string;
  new_search: boolean;
}

export interface ChatResponse {
  type: "text" | "product_card";
  output?: string;
  quick_replies?: string[];
  products?: ProductCard[];
  search_type?: string;
  session_id: string;
  message_count: number;
  search_state?: {
    status: string;
    category?: string;
    can_continue: boolean;
    search_count: number;
    max_searches: number;
    message?: string;
  };
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