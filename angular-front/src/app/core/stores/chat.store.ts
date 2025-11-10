import {
  Injectable,
  signal,
  computed,
  effect,
  inject,
  PLATFORM_ID,
} from '@angular/core';
import { isPlatformBrowser } from '@angular/common';
import {
  ChatMessage,
  Product,
  SavedSearch,
  WebSocketMessage,
} from '../../shared/types';
import { generateId } from '../../shared/utils';
import { WebSocketService } from '../services/websocket.service';
import { ApiService } from '../services/api.service';

const STORAGE_KEYS = {
  MESSAGES: 'chat_messages',
  SESSION_ID: 'session_id',
  CATEGORY: 'search_category',
  SAVED_SEARCH: 'saved_search',
} as const;

@Injectable({
  providedIn: 'root',
})
export class ChatStore {
  private readonly platformId = inject(PLATFORM_ID);
  private readonly wsService = inject(WebSocketService);
  private readonly apiService = inject(ApiService);

  // Core state with signals
  readonly messages = signal<ChatMessage[]>([]);
  readonly sessionId = signal<string>(generateId());
  readonly isLoading = signal(false);
  readonly category = signal<string | null>(null);
  readonly searchStatus = signal<'idle' | 'searching' | 'completed'>('idle');

  // Computed values
  readonly hasMessages = computed(() => this.messages().length > 0);
  readonly lastMessage = computed(() => {
    const msgs = this.messages();
    return msgs.length > 0 ? msgs[msgs.length - 1] : null;
  });
  readonly userMessages = computed(() =>
    this.messages().filter((m) => m.role === 'user')
  );
  readonly assistantMessages = computed(() =>
    this.messages().filter((m) => m.role === 'assistant')
  );

  constructor() {
    // Load from localStorage on initialization
    if (isPlatformBrowser(this.platformId)) {
      this.loadFromStorage();
    }

    // Set up WebSocket message handling
    this.wsService.messages$.subscribe((msg) => this.handleWebSocketMessage(msg));

    // Persist to localStorage when state changes
    effect(() => {
      if (isPlatformBrowser(this.platformId)) {
        this.saveToStorage();
      }
    });
  }

  // Add a new message
  addMessage(message: Omit<ChatMessage, 'id' | 'timestamp'>): ChatMessage {
    const newMessage: ChatMessage = {
      ...message,
      id: generateId(),
      timestamp: Date.now(),
      isLocal: true,
    };

    this.messages.update((msgs) => [...msgs, newMessage]);
    return newMessage;
  }

  // Update an existing message
  updateMessage(id: string, updates: Partial<ChatMessage>): void {
    this.messages.update((msgs) =>
      msgs.map((msg) => (msg.id === id ? { ...msg, ...updates } : msg))
    );
  }

  // Remove a message
  removeMessage(id: string): void {
    this.messages.update((msgs) => msgs.filter((msg) => msg.id !== id));
  }

  // Clear all messages
  clearMessages(): void {
    this.messages.set([]);
    this.category.set(null);
    this.searchStatus.set('idle');
  }

  // Send a message via WebSocket
  sendMessage(
    content: string,
    accessToken?: string,
    country?: string,
    language?: string,
    currency?: string,
    newSearch?: boolean
  ): void {
    // Add user message immediately
    this.addMessage({
      role: 'user',
      content,
    });

    this.isLoading.set(true);

    // Send via WebSocket
    this.wsService.sendMessage(
      content,
      this.sessionId(),
      accessToken,
      country,
      language,
      currency,
      newSearch
    );
  }

  // Send a message via HTTP API (fallback)
  async sendMessageHttp(
    content: string,
    accessToken?: string,
    country?: string,
    language?: string,
    currency?: string,
    newSearch?: boolean
  ): Promise<void> {
    // Add user message immediately
    this.addMessage({
      role: 'user',
      content,
    });

    this.isLoading.set(true);

    try {
      const response = await this.apiService
        .sendMessage(
          content,
          this.sessionId(),
          accessToken,
          country,
          language,
          currency,
          newSearch
        )
        .toPromise();

      if (response) {
        this.addMessage({
          role: 'assistant',
          content: response.message,
          quick_replies: response.quick_replies,
          products: response.products,
          search_type: response.search_type,
        });

        // Update session ID if returned
        if (response.session_id) {
          this.sessionId.set(response.session_id);
        }
      }
    } catch (error) {
      console.error('Failed to send message:', error);
      this.addMessage({
        role: 'assistant',
        content: 'Sorry, I encountered an error. Please try again.',
      });
    } finally {
      this.isLoading.set(false);
    }
  }

  // Handle incoming WebSocket messages
  private handleWebSocketMessage(msg: WebSocketMessage): void {
    switch (msg.type) {
      case 'chat':
        if (msg.output) {
          this.addMessage({
            role: 'assistant',
            content: msg.output,
            quick_replies: msg.quick_replies,
            products: msg.products,
            search_type: msg.search_type,
          });
        }

        if (msg.session_id) {
          this.sessionId.set(msg.session_id);
        }

        if (msg.search_state?.category) {
          this.category.set(msg.search_state.category);
        }

        if (msg.search_state?.status) {
          this.searchStatus.set(
            msg.search_state.status as 'idle' | 'searching' | 'completed'
          );
        }

        this.isLoading.set(false);
        break;

      case 'error':
        this.addMessage({
          role: 'assistant',
          content: msg.error || 'An error occurred',
        });
        this.isLoading.set(false);
        break;

      case 'assistant_message_sync':
      case 'message_sync':
        if (msg.output && msg.message_id) {
          // Check if message already exists
          const exists = this.messages().some((m) => m.id === msg.message_id);
          if (!exists) {
            this.addMessage({
              role: 'assistant',
              content: msg.output,
              quick_replies: msg.quick_replies,
              products: msg.products,
              search_type: msg.search_type,
              isLocal: false,
            });
          }
        }
        break;

      case 'user_message_sync':
        if (msg.message && msg.message_id) {
          const exists = this.messages().some((m) => m.id === msg.message_id);
          if (!exists) {
            this.addMessage({
              role: 'user',
              content: msg.message,
              isLocal: false,
            });
          }
        }
        break;

      case 'session_changed':
        if (msg.session_id) {
          this.loadSession(msg.session_id);
        }
        break;

      case 'pong':
        // Heartbeat response
        break;

      default:
        console.warn('[ChatStore] Unknown WebSocket message type:', msg.type);
    }
  }

  // Load session from server
  async loadSession(sessionId: string, accessToken?: string): Promise<void> {
    try {
      const response = await this.apiService
        .getSession(sessionId, accessToken)
        .toPromise();

      if (response) {
        this.sessionId.set(response.session_id);
        this.messages.set(
          response.messages.map((msg) => ({
            id: generateId(),
            role: msg.role as 'user' | 'assistant',
            content: msg.content,
            timestamp: msg.timestamp ? new Date(msg.timestamp).getTime() : Date.now(),
            quick_replies: msg.quick_replies,
            products: msg.products,
            search_type: msg.search_type,
            isLocal: false,
          }))
        );

        if (response.search_state?.category) {
          this.category.set(response.search_state.category);
        }

        if (response.search_state?.status) {
          this.searchStatus.set(
            response.search_state.status as 'idle' | 'searching' | 'completed'
          );
        }
      }
    } catch (error) {
      console.error('Failed to load session:', error);
    }
  }

  // Start a new session
  startNewSession(): void {
    this.sessionId.set(generateId());
    this.clearMessages();
  }

  // Save current search
  saveSearch(): SavedSearch {
    return {
      messages: this.messages(),
      sessionId: this.sessionId(),
      category: this.category() || 'general',
      timestamp: Date.now(),
    };
  }

  // Load a saved search
  loadSavedSearch(search: SavedSearch): void {
    this.messages.set(search.messages);
    this.sessionId.set(search.sessionId);
    this.category.set(search.category);
  }

  // Storage operations
  private saveToStorage(): void {
    if (!isPlatformBrowser(this.platformId)) return;

    try {
      localStorage.setItem(STORAGE_KEYS.MESSAGES, JSON.stringify(this.messages()));
      localStorage.setItem(STORAGE_KEYS.SESSION_ID, this.sessionId());
      if (this.category()) {
        localStorage.setItem(STORAGE_KEYS.CATEGORY, this.category()!);
      }
    } catch (error) {
      console.error('Failed to save to localStorage:', error);
    }
  }

  private loadFromStorage(): void {
    if (!isPlatformBrowser(this.platformId)) return;

    try {
      const messagesJson = localStorage.getItem(STORAGE_KEYS.MESSAGES);
      if (messagesJson) {
        const messages = JSON.parse(messagesJson);
        this.messages.set(messages);
      }

      const sessionId = localStorage.getItem(STORAGE_KEYS.SESSION_ID);
      if (sessionId) {
        this.sessionId.set(sessionId);
      }

      const category = localStorage.getItem(STORAGE_KEYS.CATEGORY);
      if (category) {
        this.category.set(category);
      }
    } catch (error) {
      console.error('Failed to load from localStorage:', error);
    }
  }

  clearStorage(): void {
    if (!isPlatformBrowser(this.platformId)) return;

    try {
      localStorage.removeItem(STORAGE_KEYS.MESSAGES);
      localStorage.removeItem(STORAGE_KEYS.SESSION_ID);
      localStorage.removeItem(STORAGE_KEYS.CATEGORY);
      localStorage.removeItem(STORAGE_KEYS.SAVED_SEARCH);
    } catch (error) {
      console.error('Failed to clear localStorage:', error);
    }
  }
}
