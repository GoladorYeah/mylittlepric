import {
  Injectable,
  signal,
  computed,
  effect,
  inject,
  PLATFORM_ID,
} from '@angular/core';
import { isPlatformBrowser } from '@angular/common';
import { ApiService } from '../services/api.service';
import { AuthStore } from './auth.store';

const STORAGE_KEYS = {
  COUNTRY: 'user_country',
  LANGUAGE: 'user_language',
  CURRENCY: 'user_currency',
} as const;

// Helper function to get currency for a country
export function getCurrencyForCountry(countryCode: string): string {
  const currencyMap: Record<string, string> = {
    US: 'USD',
    GB: 'GBP',
    CA: 'CAD',
    AU: 'AUD',
    DE: 'EUR',
    FR: 'EUR',
    ES: 'EUR',
    IT: 'EUR',
    NL: 'EUR',
    BE: 'EUR',
    CH: 'CHF',
    AT: 'EUR',
    SE: 'SEK',
    NO: 'NOK',
    DK: 'DKK',
    FI: 'EUR',
    PL: 'PLN',
    CZ: 'CZK',
    PT: 'EUR',
    GR: 'EUR',
    IE: 'EUR',
    JP: 'JPY',
    KR: 'KRW',
    CN: 'CNY',
    IN: 'INR',
    SG: 'SGD',
    HK: 'HKD',
    TW: 'TWD',
    NZ: 'NZD',
    MX: 'MXN',
    BR: 'BRL',
    AR: 'ARS',
    CL: 'CLP',
    ZA: 'ZAR',
    AE: 'AED',
    SA: 'SAR',
    TR: 'TRY',
    RU: 'RUB',
    UA: 'UAH',
    IL: 'ILS',
    EG: 'EGP',
    TH: 'THB',
    MY: 'MYR',
    ID: 'IDR',
    PH: 'PHP',
    VN: 'VND',
  };

  return currencyMap[countryCode.toUpperCase()] || 'USD';
}

@Injectable({
  providedIn: 'root',
})
export class PreferencesStore {
  private readonly platformId = inject(PLATFORM_ID);
  private readonly apiService = inject(ApiService);
  private readonly authStore = inject(AuthStore);

  // Core preferences state
  readonly country = signal<string>('US');
  readonly language = signal<string>('en');
  readonly currency = signal<string>('USD');

  // Computed values
  readonly locale = computed(() => `${this.language()}-${this.country()}`);

  constructor() {
    // Load from localStorage on init
    if (isPlatformBrowser(this.platformId)) {
      this.loadFromStorage();
    }

    // Persist to localStorage when state changes
    effect(() => {
      if (isPlatformBrowser(this.platformId)) {
        this.saveToStorage();
      }
    });

    // Sync preferences when user logs in
    effect(() => {
      const isAuthenticated = this.authStore.isAuthenticated();
      if (isAuthenticated && isPlatformBrowser(this.platformId)) {
        // Small delay to ensure auth is fully established
        setTimeout(() => this.syncPreferencesToServer(), 100);
      }
    });
  }

  // Set country
  setCountry(countryCode: string): void {
    this.country.set(countryCode.toUpperCase());
    // Auto-update currency based on country
    const newCurrency = getCurrencyForCountry(countryCode);
    this.currency.set(newCurrency);
  }

  // Set language
  setLanguage(languageCode: string): void {
    this.language.set(languageCode.toLowerCase());
  }

  // Set currency
  setCurrency(currencyCode: string): void {
    this.currency.set(currencyCode.toUpperCase());
  }

  // Set all preferences at once
  setPreferences(
    country?: string,
    language?: string,
    currency?: string
  ): void {
    if (country) this.country.set(country.toUpperCase());
    if (language) this.language.set(language.toLowerCase());
    if (currency) this.currency.set(currency.toUpperCase());
  }

  // Sync preferences to server (for authenticated users)
  async syncPreferencesToServer(): Promise<void> {
    const accessToken = this.authStore.accessToken();
    if (!accessToken) {
      console.log('Not authenticated, skipping preference sync');
      return;
    }

    try {
      await this.apiService
        .updatePreferences(
          {
            country: this.country(),
            language: this.language(),
            currency: this.currency(),
          },
          accessToken
        )
        .toPromise();
      console.log('✅ Preferences synced to server');
    } catch (error) {
      console.error('Failed to sync preferences to server:', error);
    }
  }

  // Load preferences from server (for authenticated users)
  async loadPreferencesFromServer(): Promise<void> {
    const accessToken = this.authStore.accessToken();
    if (!accessToken) {
      return;
    }

    try {
      const response = await this.apiService
        .getPreferences(accessToken)
        .toPromise();

      if (response?.has_preferences && response.preferences) {
        const prefs = response.preferences;
        if (prefs.country) this.country.set(prefs.country);
        if (prefs.language) this.language.set(prefs.language);
        if (prefs.currency) this.currency.set(prefs.currency);
        console.log('✅ Loaded preferences from server');
      }
    } catch (error) {
      console.error('Failed to load preferences from server:', error);
    }
  }

  // Storage operations
  private saveToStorage(): void {
    if (!isPlatformBrowser(this.platformId)) return;

    try {
      localStorage.setItem(STORAGE_KEYS.COUNTRY, this.country());
      localStorage.setItem(STORAGE_KEYS.LANGUAGE, this.language());
      localStorage.setItem(STORAGE_KEYS.CURRENCY, this.currency());
    } catch (error) {
      console.error('Failed to save preferences to localStorage:', error);
    }
  }

  private loadFromStorage(): void {
    if (!isPlatformBrowser(this.platformId)) return;

    try {
      const country = localStorage.getItem(STORAGE_KEYS.COUNTRY);
      if (country) this.country.set(country);

      const language = localStorage.getItem(STORAGE_KEYS.LANGUAGE);
      if (language) this.language.set(language);

      const currency = localStorage.getItem(STORAGE_KEYS.CURRENCY);
      if (currency) this.currency.set(currency);
    } catch (error) {
      console.error('Failed to load preferences from localStorage:', error);
    }
  }

  clearStorage(): void {
    if (!isPlatformBrowser(this.platformId)) return;

    try {
      localStorage.removeItem(STORAGE_KEYS.COUNTRY);
      localStorage.removeItem(STORAGE_KEYS.LANGUAGE);
      localStorage.removeItem(STORAGE_KEYS.CURRENCY);
    } catch (error) {
      console.error('Failed to clear preferences from localStorage:', error);
    }
  }
}
