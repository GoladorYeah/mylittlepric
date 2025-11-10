import {
  Injectable,
  signal,
  computed,
  effect,
  inject,
  PLATFORM_ID,
} from '@angular/core';
import { isPlatformBrowser } from '@angular/common';
import { User, UserPreferences, SavedSearchBackend } from '../../shared/types';
import { ApiService } from '../services/api.service';

const STORAGE_KEYS = {
  USER: 'auth_user',
  ACCESS_TOKEN: 'access_token',
  PREFERENCES: 'user_preferences',
} as const;

@Injectable({
  providedIn: 'root',
})
export class AuthStore {
  private readonly platformId = inject(PLATFORM_ID);
  private readonly apiService = inject(ApiService);

  // Core auth state
  readonly user = signal<User | null>(null);
  readonly accessToken = signal<string | null>(null);
  readonly isLoading = signal(false);
  readonly error = signal<string | null>(null);

  // User preferences
  readonly preferences = signal<UserPreferences>({
    country: undefined,
    currency: undefined,
    language: undefined,
    saved_search: null,
  });

  // Computed values
  readonly isAuthenticated = computed(() => this.user() !== null);
  readonly userName = computed(() => this.user()?.name || null);
  readonly userEmail = computed(() => this.user()?.email || null);
  readonly userPicture = computed(() => this.user()?.picture || null);
  readonly hasSavedSearch = computed(
    () => this.preferences().saved_search !== null
  );

  constructor() {
    // Load from storage on init
    if (isPlatformBrowser(this.platformId)) {
      this.loadFromStorage();
    }

    // Persist to storage on changes
    effect(() => {
      if (isPlatformBrowser(this.platformId)) {
        this.saveToStorage();
      }
    });
  }

  // Check authentication status
  async checkAuth(): Promise<void> {
    this.isLoading.set(true);
    this.error.set(null);

    try {
      const response = await this.apiService.checkAuth().toPromise();

      if (response?.authenticated && response.user) {
        this.user.set(response.user);
        // Note: Access token is typically in httpOnly cookie, not returned here
      } else {
        this.user.set(null);
        this.accessToken.set(null);
      }
    } catch (error: any) {
      console.error('Auth check failed:', error);
      this.error.set(error.message || 'Authentication check failed');
      this.user.set(null);
      this.accessToken.set(null);
    } finally {
      this.isLoading.set(false);
    }
  }

  // Login (redirect to OAuth)
  login(provider: 'google' = 'google'): void {
    if (!isPlatformBrowser(this.platformId)) return;

    const redirectUrl = `${window.location.origin}/api/auth/${provider}`;
    window.location.href = redirectUrl;
  }

  // Logout
  async logout(): Promise<void> {
    this.isLoading.set(true);
    this.error.set(null);

    try {
      await this.apiService.logout().toPromise();
      this.clearAuth();
    } catch (error: any) {
      console.error('Logout failed:', error);
      this.error.set(error.message || 'Logout failed');
      // Clear anyway on error
      this.clearAuth();
    } finally {
      this.isLoading.set(false);
    }
  }

  // Clear authentication state
  private clearAuth(): void {
    this.user.set(null);
    this.accessToken.set(null);
    this.preferences.set({
      country: undefined,
      currency: undefined,
      language: undefined,
      saved_search: null,
    });
    this.clearStorage();
  }

  // Set user manually (after OAuth callback)
  setUser(user: User, accessToken?: string): void {
    this.user.set(user);
    if (accessToken) {
      this.accessToken.set(accessToken);
    }
  }

  // Load user preferences
  async loadPreferences(): Promise<void> {
    const token = this.accessToken();
    if (!token) {
      console.warn('Cannot load preferences: not authenticated');
      return;
    }

    this.isLoading.set(true);
    this.error.set(null);

    try {
      const response = await this.apiService.getPreferences(token).toPromise();

      if (response?.has_preferences && response.preferences) {
        this.preferences.set(response.preferences);
      }
    } catch (error: any) {
      console.error('Failed to load preferences:', error);
      this.error.set(error.message || 'Failed to load preferences');
    } finally {
      this.isLoading.set(false);
    }
  }

  // Update user preferences
  async updatePreferences(
    updates: Partial<UserPreferences>
  ): Promise<void> {
    const token = this.accessToken();
    if (!token) {
      console.warn('Cannot update preferences: not authenticated');
      return;
    }

    this.isLoading.set(true);
    this.error.set(null);

    try {
      await this.apiService.updatePreferences(updates, token).toPromise();
      this.preferences.update((prefs) => ({ ...prefs, ...updates }));
    } catch (error: any) {
      console.error('Failed to update preferences:', error);
      this.error.set(error.message || 'Failed to update preferences');
      throw error;
    } finally {
      this.isLoading.set(false);
    }
  }

  // Save search to preferences
  async saveSearch(search: SavedSearchBackend): Promise<void> {
    const token = this.accessToken();
    if (!token) {
      console.warn('Cannot save search: not authenticated');
      return;
    }

    this.isLoading.set(true);
    this.error.set(null);

    try {
      await this.apiService.saveSearch(search, token).toPromise();
      this.preferences.update((prefs) => ({
        ...prefs,
        saved_search: search,
      }));
    } catch (error: any) {
      console.error('Failed to save search:', error);
      this.error.set(error.message || 'Failed to save search');
      throw error;
    } finally {
      this.isLoading.set(false);
    }
  }

  // Clear saved search
  async clearSavedSearch(): Promise<void> {
    const token = this.accessToken();
    if (!token) {
      console.warn('Cannot clear saved search: not authenticated');
      return;
    }

    this.isLoading.set(true);
    this.error.set(null);

    try {
      await this.apiService.clearSavedSearch(token).toPromise();
      this.preferences.update((prefs) => ({
        ...prefs,
        saved_search: null,
      }));
    } catch (error: any) {
      console.error('Failed to clear saved search:', error);
      this.error.set(error.message || 'Failed to clear saved search');
      throw error;
    } finally {
      this.isLoading.set(false);
    }
  }

  // Get saved search
  getSavedSearch(): SavedSearchBackend | null {
    return this.preferences().saved_search || null;
  }

  // Update locale settings
  setLocale(country?: string, language?: string, currency?: string): void {
    this.preferences.update((prefs) => ({
      ...prefs,
      country: country || prefs.country,
      language: language || prefs.language,
      currency: currency || prefs.currency,
    }));
  }

  // Storage operations
  private saveToStorage(): void {
    if (!isPlatformBrowser(this.platformId)) return;

    try {
      const user = this.user();
      if (user) {
        localStorage.setItem(STORAGE_KEYS.USER, JSON.stringify(user));
      } else {
        localStorage.removeItem(STORAGE_KEYS.USER);
      }

      const token = this.accessToken();
      if (token) {
        localStorage.setItem(STORAGE_KEYS.ACCESS_TOKEN, token);
      } else {
        localStorage.removeItem(STORAGE_KEYS.ACCESS_TOKEN);
      }

      const prefs = this.preferences();
      localStorage.setItem(STORAGE_KEYS.PREFERENCES, JSON.stringify(prefs));
    } catch (error) {
      console.error('Failed to save auth to localStorage:', error);
    }
  }

  private loadFromStorage(): void {
    if (!isPlatformBrowser(this.platformId)) return;

    try {
      const userJson = localStorage.getItem(STORAGE_KEYS.USER);
      if (userJson) {
        const user = JSON.parse(userJson);
        this.user.set(user);
      }

      const token = localStorage.getItem(STORAGE_KEYS.ACCESS_TOKEN);
      if (token) {
        this.accessToken.set(token);
      }

      const prefsJson = localStorage.getItem(STORAGE_KEYS.PREFERENCES);
      if (prefsJson) {
        const prefs = JSON.parse(prefsJson);
        this.preferences.set(prefs);
      }
    } catch (error) {
      console.error('Failed to load auth from localStorage:', error);
    }
  }

  private clearStorage(): void {
    if (!isPlatformBrowser(this.platformId)) return;

    try {
      localStorage.removeItem(STORAGE_KEYS.USER);
      localStorage.removeItem(STORAGE_KEYS.ACCESS_TOKEN);
      localStorage.removeItem(STORAGE_KEYS.PREFERENCES);
    } catch (error) {
      console.error('Failed to clear auth from localStorage:', error);
    }
  }
}
