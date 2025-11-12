// frontend/src/shared/lib/session-api.ts
import { fetchWithAuth } from './auth-api';

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

export interface ActiveSessionResponse {
  session: {
    session_id: string;
    message_count: number;
    search_state: {
      status: string;
      category: string;
    };
    created_at: string;
    updated_at: string;
    expires_at: string;
  } | null;
  has_active_session: boolean;
}

export interface SignSessionResponse {
  signed_session_id: string;
  session_id: string;
  expires_at: string;
}

export class SessionAPI {
  /**
   * Get active session for authenticated user
   * Returns the most recent session that hasn't expired
   */
  static async getActiveSession(): Promise<ActiveSessionResponse> {
    const url = `${API_URL}/api/sessions/active`;
    const response = await fetchWithAuth(url, {
      method: 'GET',
    });

    if (!response.ok) {
      let errorMessage = 'Failed to fetch active session';
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
   * Link anonymous session to authenticated user (when user logs in)
   */
  static async linkSessionToUser(sessionId: string): Promise<void> {
    const url = `${API_URL}/api/sessions/link`;
    const response = await fetchWithAuth(url, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ session_id: sessionId }),
    });

    if (!response.ok) {
      let errorMessage = 'Failed to link session to user';
      try {
        const error = await response.json();
        errorMessage = error.message || error.error || errorMessage;
      } catch (e) {
        // If response is not JSON, use default error message
      }
      throw new Error(errorMessage);
    }
  }

  /**
   * Get signed session ID for additional security
   * Signed sessions are protected with HMAC and include ownership validation
   */
  static async signSession(sessionId: string): Promise<SignSessionResponse> {
    const url = `${API_URL}/api/sessions/sign`;
    const response = await fetchWithAuth(url, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ session_id: sessionId }),
    });

    if (!response.ok) {
      let errorMessage = 'Failed to sign session';
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
