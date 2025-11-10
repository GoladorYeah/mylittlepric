import { Component, OnInit, signal, computed, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router } from '@angular/router';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatCardModule } from '@angular/material/card';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatMenuModule } from '@angular/material/menu';
import { MatTooltipModule } from '@angular/material/tooltip';
import { formatDistanceToNow } from 'date-fns';
import { enUS, ru, uk, type Locale } from 'date-fns/locale';

import { AuthStore, ChatStore } from '../../../../core/stores';
import { PreferencesStore } from '../../../../core/stores/preferences.store';
import { ApiService } from '../../../../core/services/api.service';
import { SearchHistoryRecord } from '../../../../shared/types';
import { ProductCardComponent } from '../../../products/components/product-card/product-card.component';

const localeMap: Record<string, Locale> = {
  ru,
  uk,
  en: enUS,
};

@Component({
  selector: 'app-history-page',
  standalone: true,
  imports: [
    CommonModule,
    MatButtonModule,
    MatIconModule,
    MatCardModule,
    MatProgressSpinnerModule,
    MatMenuModule,
    MatTooltipModule,
    ProductCardComponent,
  ],
  templateUrl: './history-page.component.html',
  styleUrl: './history-page.component.scss',
})
export class HistoryPageComponent implements OnInit {
  private readonly router = inject(Router);
  private readonly authStore = inject(AuthStore);
  private readonly chatStore = inject(ChatStore);
  private readonly preferencesStore = inject(PreferencesStore);
  private readonly apiService = inject(ApiService);

  // Reactive state using signals
  readonly history = signal<SearchHistoryRecord[]>([]);
  readonly loading = signal(true);
  readonly error = signal<string | null>(null);
  readonly hasMore = signal(false);
  readonly offset = signal(0);
  readonly expandedItems = signal<Set<string>>(new Set());
  readonly openMenuId = signal<string | null>(null);

  readonly limit = 50;

  // Computed values
  readonly isAuthenticated = this.authStore.isAuthenticated;
  readonly authLoading = this.authStore.isLoading;
  readonly language = this.preferencesStore.language;

  ngOnInit(): void {
    // Load history for both authenticated and anonymous users
    if (!this.authLoading()) {
      this.loadHistory(true);
    }
  }

  async loadHistory(resetOffset = false): Promise<void> {
    // Allow both authenticated and anonymous users to view history
    if (this.authLoading()) {
      this.loading.set(false);
      return;
    }

    try {
      this.loading.set(true);
      this.error.set(null);
      const currentOffset = resetOffset ? 0 : this.offset();

      const accessToken = this.authStore.accessToken();
      const sessionId = this.chatStore.sessionId();

      const response = await this.apiService
        .getSearchHistory(this.limit, currentOffset, sessionId, accessToken || undefined)
        .toPromise();

      if (!response) {
        throw new Error('No response from server');
      }

      if (resetOffset) {
        this.history.set(response.items);
        this.offset.set(0);
      } else {
        this.history.update((prev) => [...prev, ...response.items]);
      }

      this.hasMore.set(response.has_more);
      if (!resetOffset) {
        this.offset.set(currentOffset + response.items.length);
      }
    } catch (err) {
      this.error.set(err instanceof Error ? err.message : 'Failed to load history');
    } finally {
      this.loading.set(false);
    }
  }

  toggleExpanded(id: string): void {
    this.expandedItems.update((prev) => {
      const newSet = new Set(prev);
      if (newSet.has(id)) {
        newSet.delete(id);
      } else {
        newSet.add(id);
      }
      return newSet;
    });
  }

  isExpanded(id: string): boolean {
    return this.expandedItems().has(id);
  }

  async handleDelete(id: string, event: Event): Promise<void> {
    event.stopPropagation();

    try {
      const deletedItem = this.history().find((item) => item.id === id);
      const accessToken = this.authStore.accessToken();

      await this.apiService
        .deleteSearchHistory(id, accessToken || undefined)
        .toPromise();

      this.history.update((prev) => prev.filter((item) => item.id !== id));

      // Clear chat if this was the current session
      if (deletedItem && deletedItem.session_id === this.chatStore.sessionId()) {
        this.chatStore.clearMessages();
        localStorage.removeItem('chat_session_id');
      }
    } catch (err) {
      this.error.set(err instanceof Error ? err.message : 'Failed to delete');
    }
  }

  async handleClearAll(): Promise<void> {
    if (!confirm('Are you sure you want to delete all search history?')) {
      return;
    }

    try {
      const accessToken = this.authStore.accessToken();
      await this.apiService
        .deleteAllSearchHistory(accessToken || undefined)
        .toPromise();

      this.history.set([]);
      this.chatStore.clearMessages();
      localStorage.removeItem('chat_session_id');
    } catch (err) {
      this.error.set(
        err instanceof Error ? err.message : 'Failed to delete all history'
      );
    }
  }

  handleViewChat(sessionId: string): void {
    if (sessionId) {
      // Navigate to chat with this session
      this.router.navigate(['/chat'], { queryParams: { session_id: sessionId } });
    }
  }

  getTimeAgo(dateString: string): string {
    try {
      const locale = localeMap[this.language()] || enUS;
      return formatDistanceToNow(new Date(dateString), {
        addSuffix: true,
        locale,
      });
    } catch {
      return new Date(dateString).toLocaleDateString();
    }
  }

  setOpenMenuId(id: string | null): void {
    this.openMenuId.set(id);
  }

  hasProducts(item: SearchHistoryRecord): boolean {
    return !!(item.products_found && item.products_found.length > 0);
  }

  goToChat(): void {
    this.router.navigate(['/chat']);
  }
}
