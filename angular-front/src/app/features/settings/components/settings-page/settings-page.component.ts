import { Component, OnInit, signal, inject, computed } from '@angular/core';
import { CommonModule, Location } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatCardModule } from '@angular/material/card';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';
import { MatRadioModule } from '@angular/material/radio';

import { PreferencesStore, getCurrencyForCountry } from '../../../../core/stores/preferences.store';
import { AuthStore } from '../../../../core/stores/auth.store';
import { ThemeService, Theme } from '../../../../core/services/theme.service';
import { ApiService } from '../../../../core/services/api.service';

interface Country {
  code: string;
  name: string;
  flag: string;
}

interface Language {
  code: string;
  name: string;
  nativeName: string;
}

const COUNTRIES: Country[] = [
  { code: 'us', name: 'United States', flag: '🇺🇸' },
  { code: 'gb', name: 'United Kingdom', flag: '🇬🇧' },
  { code: 'ca', name: 'Canada', flag: '🇨🇦' },
  { code: 'au', name: 'Australia', flag: '🇦🇺' },
  { code: 'de', name: 'Germany', flag: '🇩🇪' },
  { code: 'fr', name: 'France', flag: '🇫🇷' },
  { code: 'es', name: 'Spain', flag: '🇪🇸' },
  { code: 'it', name: 'Italy', flag: '🇮🇹' },
  { code: 'nl', name: 'Netherlands', flag: '🇳🇱' },
  { code: 'be', name: 'Belgium', flag: '🇧🇪' },
  { code: 'ch', name: 'Switzerland', flag: '🇨🇭' },
  { code: 'at', name: 'Austria', flag: '🇦🇹' },
  { code: 'se', name: 'Sweden', flag: '🇸🇪' },
  { code: 'no', name: 'Norway', flag: '🇳🇴' },
  { code: 'dk', name: 'Denmark', flag: '🇩🇰' },
  { code: 'fi', name: 'Finland', flag: '🇫🇮' },
  { code: 'pl', name: 'Poland', flag: '🇵🇱' },
  { code: 'cz', name: 'Czech Republic', flag: '🇨🇿' },
  { code: 'pt', name: 'Portugal', flag: '🇵🇹' },
  { code: 'gr', name: 'Greece', flag: '🇬🇷' },
  { code: 'ie', name: 'Ireland', flag: '🇮🇪' },
  { code: 'jp', name: 'Japan', flag: '🇯🇵' },
  { code: 'kr', name: 'South Korea', flag: '🇰🇷' },
  { code: 'cn', name: 'China', flag: '🇨🇳' },
  { code: 'in', name: 'India', flag: '🇮🇳' },
  { code: 'sg', name: 'Singapore', flag: '🇸🇬' },
  { code: 'hk', name: 'Hong Kong', flag: '🇭🇰' },
  { code: 'tw', name: 'Taiwan', flag: '🇹🇼' },
  { code: 'nz', name: 'New Zealand', flag: '🇳🇿' },
  { code: 'mx', name: 'Mexico', flag: '🇲🇽' },
  { code: 'br', name: 'Brazil', flag: '🇧🇷' },
  { code: 'ar', name: 'Argentina', flag: '🇦🇷' },
  { code: 'cl', name: 'Chile', flag: '🇨🇱' },
  { code: 'za', name: 'South Africa', flag: '🇿🇦' },
  { code: 'ae', name: 'UAE', flag: '🇦🇪' },
  { code: 'sa', name: 'Saudi Arabia', flag: '🇸🇦' },
  { code: 'tr', name: 'Turkey', flag: '🇹🇷' },
  { code: 'ru', name: 'Russia', flag: '🇷🇺' },
  { code: 'ua', name: 'Ukraine', flag: '🇺🇦' },
  { code: 'il', name: 'Israel', flag: '🇮🇱' },
  { code: 'eg', name: 'Egypt', flag: '🇪🇬' },
  { code: 'th', name: 'Thailand', flag: '🇹🇭' },
  { code: 'my', name: 'Malaysia', flag: '🇲🇾' },
  { code: 'id', name: 'Indonesia', flag: '🇮🇩' },
  { code: 'ph', name: 'Philippines', flag: '🇵🇭' },
  { code: 'vn', name: 'Vietnam', flag: '🇻🇳' },
];

const LANGUAGES: Language[] = [
  { code: 'en', name: 'English', nativeName: 'English' },
  { code: 'es', name: 'Spanish', nativeName: 'Español' },
  { code: 'fr', name: 'French', nativeName: 'Français' },
  { code: 'de', name: 'German', nativeName: 'Deutsch' },
  { code: 'it', name: 'Italian', nativeName: 'Italiano' },
  { code: 'pt', name: 'Portuguese', nativeName: 'Português' },
  { code: 'ru', name: 'Russian', nativeName: 'Русский' },
  { code: 'zh', name: 'Chinese', nativeName: '中文' },
  { code: 'ja', name: 'Japanese', nativeName: '日本語' },
  { code: 'ko', name: 'Korean', nativeName: '한국어' },
  { code: 'ar', name: 'Arabic', nativeName: 'العربية' },
  { code: 'hi', name: 'Hindi', nativeName: 'हिन्दी' },
  { code: 'nl', name: 'Dutch', nativeName: 'Nederlands' },
  { code: 'pl', name: 'Polish', nativeName: 'Polski' },
  { code: 'tr', name: 'Turkish', nativeName: 'Türkçe' },
  { code: 'sv', name: 'Swedish', nativeName: 'Svenska' },
  { code: 'no', name: 'Norwegian', nativeName: 'Norsk' },
  { code: 'da', name: 'Danish', nativeName: 'Dansk' },
  { code: 'fi', name: 'Finnish', nativeName: 'Suomi' },
  { code: 'cs', name: 'Czech', nativeName: 'Čeština' },
  { code: 'el', name: 'Greek', nativeName: 'Ελληνικά' },
  { code: 'he', name: 'Hebrew', nativeName: 'עברית' },
  { code: 'th', name: 'Thai', nativeName: 'ไทย' },
  { code: 'vi', name: 'Vietnamese', nativeName: 'Tiếng Việt' },
  { code: 'id', name: 'Indonesian', nativeName: 'Bahasa Indonesia' },
  { code: 'ms', name: 'Malay', nativeName: 'Bahasa Melayu' },
  { code: 'uk', name: 'Ukrainian', nativeName: 'Українська' },
  { code: 'ro', name: 'Romanian', nativeName: 'Română' },
  { code: 'hu', name: 'Hungarian', nativeName: 'Magyar' },
  { code: 'hr', name: 'Croatian', nativeName: 'Hrvatski' },
];

@Component({
  selector: 'app-settings-page',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    MatButtonModule,
    MatIconModule,
    MatCardModule,
    MatFormFieldModule,
    MatInputModule,
    MatSelectModule,
    MatRadioModule,
  ],
  templateUrl: './settings-page.component.html',
  styleUrl: './settings-page.component.scss',
})
export class SettingsPageComponent implements OnInit {
  private readonly location = inject(Location);
  private readonly preferencesStore = inject(PreferencesStore);
  private readonly authStore = inject(AuthStore);
  private readonly themeService = inject(ThemeService);
  private readonly apiService = inject(ApiService);

  // Data
  readonly countries = COUNTRIES;
  readonly languages = LANGUAGES;

  // Search queries for filtering
  readonly countrySearchQuery = signal('');
  readonly languageSearchQuery = signal('');

  // Store references
  readonly country = this.preferencesStore.country;
  readonly language = this.preferencesStore.language;
  readonly currency = this.preferencesStore.currency;
  readonly theme = this.themeService.theme;
  readonly accessToken = this.authStore.accessToken;

  // Computed values
  readonly selectedCountry = computed(() => {
    const code = this.country().toLowerCase();
    return COUNTRIES.find((c) => c.code === code) || COUNTRIES[0];
  });

  readonly selectedLanguage = computed(() => {
    const code = this.language().toLowerCase();
    return LANGUAGES.find((l) => l.code === code) || LANGUAGES[0];
  });

  readonly filteredCountries = computed(() => {
    const query = this.countrySearchQuery().toLowerCase();
    if (!query) return COUNTRIES;
    return COUNTRIES.filter(
      (c) =>
        c.name.toLowerCase().includes(query) ||
        c.code.toLowerCase().includes(query)
    );
  });

  readonly filteredLanguages = computed(() => {
    const query = this.languageSearchQuery().toLowerCase();
    if (!query) return LANGUAGES;
    return LANGUAGES.filter(
      (l) =>
        l.name.toLowerCase().includes(query) ||
        l.nativeName.toLowerCase().includes(query) ||
        l.code.toLowerCase().includes(query)
    );
  });

  ngOnInit(): void {
    // Component initialization
  }

  async handleCountrySelect(countryCode: string): Promise<void> {
    this.preferencesStore.setCountry(countryCode);
    this.countrySearchQuery.set('');

    // Sync to server if user is authenticated
    if (this.accessToken()) {
      try {
        await this.preferencesStore.syncPreferencesToServer();
      } catch (error) {
        console.error('Failed to sync country preference:', error);
      }
    }
  }

  async handleLanguageSelect(languageCode: string): Promise<void> {
    this.preferencesStore.setLanguage(languageCode);
    this.languageSearchQuery.set('');

    // Sync to server if user is authenticated
    if (this.accessToken()) {
      try {
        await this.preferencesStore.syncPreferencesToServer();
      } catch (error) {
        console.error('Failed to sync language preference:', error);
      }
    }
  }

  async handleThemeChange(newTheme: Theme): Promise<void> {
    this.themeService.setTheme(newTheme);

    // Sync theme to server if user is authenticated
    if (this.accessToken()) {
      try {
        await this.apiService
          .updatePreferences({ theme: newTheme } as any, this.accessToken()!)
          .toPromise();
        console.log('✅ Synced theme to server:', newTheme);
      } catch (error) {
        console.error('Failed to sync theme preference:', error);
      }
    }
  }

  handleBack(): void {
    this.location.back();
  }

  isCountrySelected(countryCode: string): boolean {
    return this.country().toLowerCase() === countryCode.toLowerCase();
  }

  isLanguageSelected(languageCode: string): boolean {
    return this.language().toLowerCase() === languageCode.toLowerCase();
  }
}
