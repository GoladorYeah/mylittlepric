// frontend/src/lib/search-history-api.ts
import { fetchWithAuth } from './auth-api';
import { useChatStore } from './store';

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

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

export interface SearchHistoryListResponse {
  items: SearchHistoryRecord[];
  total: number;
  limit: number;
  offset: number;
  has_more: boolean;
}

export class SearchHistoryAPI {
  /**
   * Get search history (supports both authenticated and anonymous users)
   */
  static async getSearchHistory(limit = 20, offset = 0): Promise<SearchHistoryListResponse> {
    // Get session_id for anonymous users
    const { sessionId } = useChatStore.getState();
    const sessionParam = sessionId ? `&session_id=${sessionId}` : '';

    const url = `${API_URL}/api/search-history?limit=${limit}&offset=${offset}${sessionParam}`;
    const response = await fetchWithAuth(url, {
      method: 'GET',
    });

    if (!response.ok) {
      let errorMessage = 'Failed to fetch search history';
      try {
        const error = await response.json();
        errorMessage = error.message || errorMessage;
      } catch (e) {
        // If response is not JSON, use default error message
      }
      throw new Error(errorMessage);
    }

    return response.json();
  }

  /**
   * Delete a specific search history entry
   */
  static async deleteSearchHistory(id: string): Promise<void> {
    const url = `${API_URL}/api/search-history/${id}`;
    const response = await fetchWithAuth(url, {
      method: 'DELETE',
    });

    if (!response.ok) {
      let errorMessage = 'Failed to delete search history';
      try {
        const error = await response.json();
        errorMessage = error.message || errorMessage;
      } catch (e) {
        // If response is not JSON, use default error message
      }
      throw new Error(errorMessage);
    }
  }

  /**
   * Delete all search history (authenticated users only)
   */
  static async deleteAllSearchHistory(): Promise<void> {
    const url = `${API_URL}/api/search-history`;
    const response = await fetchWithAuth(url, {
      method: 'DELETE',
    });

    if (!response.ok) {
      let errorMessage = 'Failed to delete all search history';
      try {
        const error = await response.json();
        errorMessage = error.message || errorMessage;
      } catch (e) {
        // If response is not JSON, use default error message
      }
      throw new Error(errorMessage);
    }
  }

  /**
   * Track which product was clicked in a search
   */
  static async trackProductClick(historyId: string, productId: string): Promise<void> {
    const url = `${API_URL}/api/search-history/${historyId}/click`;
    const response = await fetchWithAuth(url, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ product_id: productId }),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.message || 'Failed to track product click');
    }
  }
}
