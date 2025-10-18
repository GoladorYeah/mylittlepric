import { useAuthStore, User, AuthTokens } from "./auth-store";
import { useChatStore } from "./store";

const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

export interface SignupRequest {
  email: string;
  password: string;
  full_name?: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface AuthResponse {
  access_token: string;
  refresh_token: string;
  user: User;
  expires_in: number;
}

export interface ErrorResponse {
  error: string;
  message: string;
}

class AuthAPI {
  private async fetch(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<Response> {
    const url = `${API_URL}/api/auth${endpoint}`;
    const headers = {
      "Content-Type": "application/json",
      ...options.headers,
    };

    return fetch(url, { ...options, headers });
  }

   async fetchWithAuth(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<Response> {
    const { accessToken, isTokenExpired, refreshToken } = useAuthStore.getState();

 // Only try to refresh if we have BOTH access token AND refresh token
  // and the access token is expired
  if (accessToken && refreshToken && isTokenExpired()) {
      try {
        await this.refreshAccessToken();
      } catch (error) {
        console.error("Failed to refresh token:", error);
        useAuthStore.getState().clearAuth();
        throw new Error("Authentication expired. Please login again.");
      }
    }

    const token = useAuthStore.getState().accessToken;
    const headers = {
      "Content-Type": "application/json",
      ...(token && { Authorization: `Bearer ${token}` }),
      ...options.headers,
    };

    return fetch(`${API_URL}/api/auth${endpoint}`, { ...options, headers });
  }

  async signup(data: SignupRequest): Promise<AuthResponse> {
    const response = await this.fetch("/signup", {
      method: "POST",
      body: JSON.stringify(data),
    });

    if (!response.ok) {
      const error: ErrorResponse = await response.json();
      throw new Error(error.message || "Signup failed");
    }

    const authData: AuthResponse = await response.json();

    // Save auth data to store
    useAuthStore.getState().setAuth(authData.user, {
      access_token: authData.access_token,
      refresh_token: authData.refresh_token,
      expires_in: authData.expires_in,
    });

    return authData;
  }

  async login(data: LoginRequest): Promise<AuthResponse> {
    const response = await this.fetch("/login", {
      method: "POST",
      body: JSON.stringify(data),
    });

    if (!response.ok) {
      const error: ErrorResponse = await response.json();
      throw new Error(error.message || "Login failed");
    }

    const authData: AuthResponse = await response.json();

    // Save auth data to store
    useAuthStore.getState().setAuth(authData.user, {
      access_token: authData.access_token,
      refresh_token: authData.refresh_token,
      expires_in: authData.expires_in,
    });

    return authData;
  }

  async logout(): Promise<void> {
    const { refreshToken } = useAuthStore.getState();

    if (refreshToken) {
      try {
        await this.fetch("/logout", {
          method: "POST",
          body: JSON.stringify({ refresh_token: refreshToken }),
        });
      } catch (error) {
        console.error("Logout error:", error);
      }
    }

    // Clear auth state
    useAuthStore.getState().clearAuth();
  }

  async refreshAccessToken(): Promise<void> {
    const { refreshToken } = useAuthStore.getState();

    if (!refreshToken) {
      throw new Error("No refresh token available");
    }

    const response = await this.fetch("/refresh", {
      method: "POST",
      body: JSON.stringify({ refresh_token: refreshToken }),
    });

    if (!response.ok) {
      throw new Error("Failed to refresh token");
    }

    const authData: AuthResponse = await response.json();

    // Update tokens
    useAuthStore.getState().updateTokens({
      access_token: authData.access_token,
      refresh_token: authData.refresh_token,
      expires_in: authData.expires_in,
    });
  }

  async getMe(): Promise<User> {
    const response = await this.fetchWithAuth("/me");

    if (!response.ok) {
      throw new Error("Failed to get user info");
    }

    const user: User = await response.json();
    useAuthStore.getState().setUser(user);
    return user;
  }

  async claimSessions(): Promise<void> {
    const { sessionId, searchHistory } = useChatStore.getState();

    // Collect all session IDs from current session and history
    const sessionIds = new Set<string>();
    if (sessionId) sessionIds.add(sessionId);
    searchHistory.forEach(item => {
      if (item.sessionId) sessionIds.add(item.sessionId);
    });

    if (sessionIds.size === 0) return;

    const response = await this.fetchWithAuth("/claim-sessions", {
      method: "POST",
      body: JSON.stringify({ session_ids: Array.from(sessionIds) }),
    });

    if (!response.ok) {
      throw new Error("Failed to claim sessions");
    }
  }
}

export const authAPI = new AuthAPI();

/**
 * Utility function for making authenticated API requests outside of auth endpoints
 * Automatically includes Authorization header if user is logged in
 * Handles token refresh if needed
 */
export async function fetchWithAuth(
  url: string,
  options: RequestInit = {}
): Promise<Response> {
  const { accessToken, isTokenExpired, refreshToken } = useAuthStore.getState();

  // Only try to refresh if we have BOTH access token AND refresh token
  // and the access token is expired
  if (accessToken && refreshToken && isTokenExpired()) {
    try {
      await authAPI["refreshAccessToken"]();
    } catch (error) {
      console.error("Failed to refresh token:", error);
      useAuthStore.getState().clearAuth();
      throw new Error("Authentication expired. Please login again.");
    }
  }

  const token = useAuthStore.getState().accessToken;
  const headers = {
    "Content-Type": "application/json",
    ...(token && { Authorization: `Bearer ${token}` }),
    ...options.headers,
  };

  return fetch(url, { ...options, headers });
}

