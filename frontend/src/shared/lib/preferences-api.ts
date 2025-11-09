// frontend/src/shared/lib/preferences-api.ts
import { fetchWithAuth } from './auth-api';

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

export interface SavedSearchMessage {
  id: string;
  role: string;
  content: string;
  timestamp: number;
  quick_replies?: string[];
  products?: any[];
  search_type?: string;
}

export interface SavedSearchData {
  session_id: string;
  category: string;
  timestamp: number;
  messages: SavedSearchMessage[];
}

export interface UserPreferences {
  user_id: string;
  country?: string;
  currency?: string;
  language?: string;
  theme?: string;
  sidebar_open?: boolean;
  saved_search?: SavedSearchData;
  created_at: string;
  updated_at: string;
}

export interface UserPreferencesUpdate {
  country?: string;
  currency?: string;
  language?: string;
  theme?: string;
  sidebar_open?: boolean;
  saved_search?: SavedSearchData | null;
}

export interface GetPreferencesResponse {
  preferences: UserPreferences | null;
  has_preferences: boolean;
}

export interface UpdatePreferencesResponse {
  success: boolean;
  preferences: UserPreferences;
}

export class PreferencesAPI {
  /**
   * Get user preferences from server
   * Returns null if user hasn't set preferences yet
   */
  static async getUserPreferences(): Promise<GetPreferencesResponse> {
    const url = `${API_URL}/api/user/preferences`;
    const response = await fetchWithAuth(url, {
      method: 'GET',
    });

    if (!response.ok) {
      let errorMessage = 'Failed to fetch user preferences';
      try {
        const error = await response.json();
        errorMessage = error.message || error.error || errorMessage;
      } catch (e) {
        // If response is not JSON, use default error message
      }
      throw new Error(errorMessage);
    }

    return response.json();
  }

  /**
   * Update user preferences on server
   * Creates preferences if they don't exist yet
   */
  static async updateUserPreferences(
    update: UserPreferencesUpdate
  ): Promise<UpdatePreferencesResponse> {
    const url = `${API_URL}/api/user/preferences`;
    const response = await fetchWithAuth(url, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(update),
    });

    if (!response.ok) {
      let errorMessage = 'Failed to update user preferences';
      try {
        const error = await response.json();
        errorMessage = error.message || error.error || errorMessage;
      } catch (e) {
        // If response is not JSON, use default error message
      }
      throw new Error(errorMessage);
    }

    return response.json();
  }
}
