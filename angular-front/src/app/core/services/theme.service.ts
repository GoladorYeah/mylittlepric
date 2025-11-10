import {
  Injectable,
  signal,
  effect,
  inject,
  PLATFORM_ID,
} from '@angular/core';
import { DOCUMENT, isPlatformBrowser } from '@angular/common';

export type Theme = 'light' | 'dark' | 'system';

const STORAGE_KEY = 'theme_preference';

@Injectable({
  providedIn: 'root',
})
export class ThemeService {
  private readonly platformId = inject(PLATFORM_ID);
  private readonly document = inject(DOCUMENT);

  readonly theme = signal<Theme>('system');
  readonly actualTheme = signal<'light' | 'dark'>('light');

  constructor() {
    if (isPlatformBrowser(this.platformId)) {
      this.loadTheme();
      this.setupMediaQuery();
    }

    // Apply theme when it changes
    effect(() => {
      if (isPlatformBrowser(this.platformId)) {
        this.applyTheme();
      }
    });
  }

  setTheme(theme: Theme): void {
    this.theme.set(theme);
    if (isPlatformBrowser(this.platformId)) {
      localStorage.setItem(STORAGE_KEY, theme);
    }
  }

  toggleTheme(): void {
    const current = this.actualTheme();
    this.setTheme(current === 'light' ? 'dark' : 'light');
  }

  private loadTheme(): void {
    const stored = localStorage.getItem(STORAGE_KEY) as Theme;
    if (stored && ['light', 'dark', 'system'].includes(stored)) {
      this.theme.set(stored);
    }
  }

  private setupMediaQuery(): void {
    const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)');

    const updateSystemTheme = () => {
      if (this.theme() === 'system') {
        this.actualTheme.set(mediaQuery.matches ? 'dark' : 'light');
      }
    };

    // Initial check
    updateSystemTheme();

    // Listen for changes
    mediaQuery.addEventListener('change', updateSystemTheme);
  }

  private applyTheme(): void {
    const theme = this.theme();
    const root = this.document.documentElement;

    // Remove existing theme classes
    root.classList.remove('light', 'dark');

    if (theme === 'system') {
      const systemTheme = window.matchMedia('(prefers-color-scheme: dark)').matches
        ? 'dark'
        : 'light';
      this.actualTheme.set(systemTheme);
      root.classList.add(systemTheme);
    } else {
      this.actualTheme.set(theme);
      root.classList.add(theme);
    }
  }
}
